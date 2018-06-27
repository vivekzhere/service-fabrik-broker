'use strict';

const _ = require('lodash');
const Promise = require('bluebird');
const catalog = require('../../broker/lib/models/catalog');
const eventmesh = require('../../eventmesh');
const config = require('../../common/config');
const logger = require('../../common/logger');
const CONST = require('../../common/constants');
const bm = require('./');
const BaseManager = require('../BaseManager');
const DBManager = require('../../broker/lib/fabrik/DBManager');
const errors = require('../../common/errors');
const Conflict = errors.Conflict;

/* jshint nonew:false */
new DBManager(); //to start the BnRStatusPoller

class DefaultBackupManager extends BaseManager {

  /**
   * @description Registers backup watcher with worker callback
   */
  registerWatcher() {
    logger.info(`Registering Backup watcher`);
    return eventmesh.server.registerWatcher('backup', 'defaultbackup', this.worker);
  }
  /**
   * @description This method does following in order
   * 1. Tries to acquire processing lock => In case of HA, Only one backup manager process will process the request
   * 2. Processes the request based on the state of the changed resource
   * 3. Release processing lock
   * @param {object} change - Change object that comes as part of apiserver watch event
   */

  worker(change) {
    /**
     * @description Patches resource with annotation key lockedByManager and value broker ip
     */
    function acquireProcessingLock() {
      logger.info('Trying to acquire processing lock for the backup request for backup guid: ', changedOptions.guid);
      // Set lockedManager annotations to true
      const patchBody = _.cloneDeep(changeObjectBody);
      let patchAnnotations = patchBody.metadata.annotations;
      patchAnnotations = patchAnnotations ? patchAnnotations : {};
      patchAnnotations.lockedByManager = config.broker_ip;
      return eventmesh.server.updateResource(CONST.APISERVER.RESOURCE_TYPES.BACKUP, CONST.APISERVER.RESOURCE_NAMES.DEFAULT_BACKUP, changedOptions.guid, patchBody)
        .tap((resource) => logger.info(`Successfully acquired processing lock for the backup request for backup guid: ${changedOptions.guid}\n` +
          `Updated resource with annotations is: `, resource));
    }

    /**
     * @description Sets lockedByManager annotation to empty string
     */

    function releaseProcessingLock() {
      logger.info('Trying to release processing lock for the backup request for backup guid: ', changedOptions.guid);
      const patchBody = {
        metadata: {
          annotations: {
            lockedByManager: ''
          }
        }
      };
      return eventmesh.server.updateResource(CONST.APISERVER.RESOURCE_TYPES.BACKUP, CONST.APISERVER.RESOURCE_NAMES.DEFAULT_BACKUP, changedOptions.guid, patchBody)
        .tap((resource) => logger.info(`Successfully released processing lock for the backup request for backup guid: ${changedOptions.guid}\n` +
          `Updated resource with annotations is: `, resource));
    }

    function processBackup() {
      logger.info('Triggering backup with the following options:', changedOptions);
      const plan = catalog.getPlan(changedOptions.plan_id);
      return bm.createManager(plan)
        .then(manager => manager.startBackup(changedOptions));
    }

    function processAbort() {
      return eventmesh.server.getOperationOptions({
          resourceId: changedOptions.instance_guid,
          operationName: CONST.APISERVER.ANNOTATION_NAMES.BACKUP,
          operationType: CONST.APISERVER.ANNOTATION_TYPES.BACKUP,
          operationId: changedOptions.guid
        })
        .then(options => {
          const changedOptions = JSON.parse(options);
          return Promise.try(() => {
            const plan = catalog.getPlan(changedOptions.plan_id);
            return bm.createManager(plan);
          }).then(manager => manager.abortLastBackup(changedOptions));
        });
    }

    function processDelete() {
      return eventmesh.server.getOperationOptions({
          resourceId: changedOptions.instance_guid,
          operationName: CONST.APISERVER.ANNOTATION_NAMES.BACKUP,
          operationType: CONST.APISERVER.ANNOTATION_TYPES.BACKUP,
          operationId: changedOptions.guid
        })
        .then(options => {
          const changedOptions = JSON.parse(options);
          return this.backupStore
            .deleteBackupFile(options)
            .catch(err => {
              return Promise
                .try(() => logger.error(`Error during delete of backup`, err))
                .tap(() => eventmesh.server.updateOperationState({
                  resourceId: changedOptions.instance_guid,
                  operationName: CONST.APISERVER.ANNOTATION_NAMES.BACKUP,
                  operationType: CONST.APISERVER.ANNOTATION_TYPES.BACKUP,
                  operationId: changedOptions.guid,
                  stateValue: CONST.APISERVER.RESOURCE_STATE.ERROR
                }))
                .tap(() => eventmesh.server.updateOperationResult({
                  resourceId: changedOptions.instance_guid,
                  operationName: CONST.APISERVER.ANNOTATION_NAMES.BACKUP,
                  operationType: CONST.APISERVER.ANNOTATION_TYPES.BACKUP,
                  operationId: changedOptions.guid,
                  value: err
                }));
            });
        });
    }

    logger.info('Changed Resource:', change);
    logger.debug('Changed resource options:', change.object.spec.options);
    const changeObjectBody = change.object;
    const changedOptions = JSON.parse(changeObjectBody.spec.options);
    logger.debug('Changed resource options(parsed):', changedOptions);
    let processingLockConflict = false;
    if (changeObjectBody.status.state === CONST.APISERVER.RESOURCE_STATE.IN_QUEUE || changeObjectBody.status.state === CONST.OPERATION.ABORT) {
      // Acquire processing lock so that in HA scenerio, only one backup-manager process processes the request
      return Promise.try(() => {
          if (!changeObjectBody.metadata.annotations || changeObjectBody.metadata.annotations.lockedByManager === '') {
            return acquireProcessingLock()
              .catch(err => {
                if (err instanceof Conflict) {
                  processingLockConflict = true;
                  logger.info(`Not able to acquire processing lock, Backup request for guid ${changedOptions.guid} is probably picked by other worker`);
                } else {
                  logger.error(`Error while trying to get processing lock for backup with guid ${changedOptions.guid}`, err);
                }
                throw err;
              });
          } else if (changeObjectBody.metadata.annotations.lockedByManager === config.broker_ip) {
            processingLockConflict = false;
          } else {
            processingLockConflict = true;
            logger.info(`Backup request for guid ${changedOptions.guid} is picked by other worker`);
          }
        })
        .then(() => {
          if (!processingLockConflict) {
            if (changeObjectBody.status.state === CONST.APISERVER.RESOURCE_STATE.IN_QUEUE) {
              return processBackup();
            } else if (changeObjectBody.status.state === CONST.OPERATION.ABORT) {
              return processAbort();
            } else if (changeObjectBody.status.state === 'delete') {
              return processDelete();
            }
          }
        })
        .catch(e => {
          // If conflict error is in acquiring processingLock then don't throw error
          if (!processingLockConflict) {
            logger.error(`Caught error while starting backup ${changedOptions.guid}`, e);
          }
        })
        .then(() => releaseProcessingLock());
    }
  }
}

const defaultBackupManager = new DefaultBackupManager();
defaultBackupManager.registerWatcher();