/*
Copyright 2019 The Service Fabrik Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/cloudfoundry-incubator/service-fabrik-broker/interoperator/api/resource/v1alpha1"
	scheme "github.com/cloudfoundry-incubator/service-fabrik-broker/interoperator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SFClustersGetter has a method to return a SFClusterInterface.
// A group's client should implement this interface.
type SFClustersGetter interface {
	SFClusters(namespace string) SFClusterInterface
}

// SFClusterInterface has methods to work with SFCluster resources.
type SFClusterInterface interface {
	Create(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.CreateOptions) (*v1alpha1.SFCluster, error)
	Update(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.UpdateOptions) (*v1alpha1.SFCluster, error)
	UpdateStatus(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.UpdateOptions) (*v1alpha1.SFCluster, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.SFCluster, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.SFClusterList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.SFCluster, err error)
	SFClusterExpansion
}

// sFClusters implements SFClusterInterface
type sFClusters struct {
	client rest.Interface
	ns     string
}

// newSFClusters returns a SFClusters
func newSFClusters(c *ResourceV1alpha1Client, namespace string) *sFClusters {
	return &sFClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sFCluster, and returns the corresponding sFCluster object, and an error if there is any.
func (c *sFClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.SFCluster, err error) {
	result = &v1alpha1.SFCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sfclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SFClusters that match those selectors.
func (c *sFClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SFClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.SFClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sfclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sFClusters.
func (c *sFClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sfclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a sFCluster and creates it.  Returns the server's representation of the sFCluster, and an error, if there is any.
func (c *sFClusters) Create(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.CreateOptions) (result *v1alpha1.SFCluster, err error) {
	result = &v1alpha1.SFCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sfclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sFCluster).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a sFCluster and updates it. Returns the server's representation of the sFCluster, and an error, if there is any.
func (c *sFClusters) Update(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.UpdateOptions) (result *v1alpha1.SFCluster, err error) {
	result = &v1alpha1.SFCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sfclusters").
		Name(sFCluster.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sFCluster).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *sFClusters) UpdateStatus(ctx context.Context, sFCluster *v1alpha1.SFCluster, opts v1.UpdateOptions) (result *v1alpha1.SFCluster, err error) {
	result = &v1alpha1.SFCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sfclusters").
		Name(sFCluster.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sFCluster).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the sFCluster and deletes it. Returns an error if one occurs.
func (c *sFClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sfclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sFClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sfclusters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched sFCluster.
func (c *sFClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.SFCluster, err error) {
	result = &v1alpha1.SFCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sfclusters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
