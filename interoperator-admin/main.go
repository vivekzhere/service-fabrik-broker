package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudfoundry-incubator/service-fabrik-broker/interoperator-admin/internal/config"
	"github.com/cloudfoundry-incubator/service-fabrik-broker/interoperator-admin/internal/router"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	setupLog := ctrl.Log.WithName("setup")

	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		setupLog.Error(err, "Error while reading kubeconfig")
		os.Exit(1)
	}
	adminConfig := config.NewAdminConfig()
	adminConfig.InitConfig()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", adminConfig.ServerPort),
		Handler: router.GetAdminRouter(kubeConfig, adminConfig),
	}

	setupLog.Info("Server starting to listen on port: ", "Port", adminConfig.ServerPort)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			setupLog.Error(err, "Could not start server")
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	setupLog.Info("Got OS shutdown signal, shutting down server gracefully...")
	server.Shutdown(context.Background())

}