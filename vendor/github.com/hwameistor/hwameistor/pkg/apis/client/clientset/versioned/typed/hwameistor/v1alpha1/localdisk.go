// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/hwameistor/hwameistor/pkg/apis/client/clientset/versioned/scheme"
	v1alpha1 "github.com/hwameistor/hwameistor/pkg/apis/hwameistor/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// LocalDisksGetter has a method to return a LocalDiskInterface.
// A group's client should implement this interface.
type LocalDisksGetter interface {
	LocalDisks() LocalDiskInterface
}

// LocalDiskInterface has methods to work with LocalDisk resources.
type LocalDiskInterface interface {
	Create(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.CreateOptions) (*v1alpha1.LocalDisk, error)
	Update(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.UpdateOptions) (*v1alpha1.LocalDisk, error)
	UpdateStatus(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.UpdateOptions) (*v1alpha1.LocalDisk, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.LocalDisk, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.LocalDiskList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LocalDisk, err error)
	LocalDiskExpansion
}

// localDisks implements LocalDiskInterface
type localDisks struct {
	client rest.Interface
}

// newLocalDisks returns a LocalDisks
func newLocalDisks(c *HwameistorV1alpha1Client) *localDisks {
	return &localDisks{
		client: c.RESTClient(),
	}
}

// Get takes name of the localDisk, and returns the corresponding localDisk object, and an error if there is any.
func (c *localDisks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LocalDisk, err error) {
	result = &v1alpha1.LocalDisk{}
	err = c.client.Get().
		Resource("localdisks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of LocalDisks that match those selectors.
func (c *localDisks) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LocalDiskList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.LocalDiskList{}
	err = c.client.Get().
		Resource("localdisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested localDisks.
func (c *localDisks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("localdisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a localDisk and creates it.  Returns the server's representation of the localDisk, and an error, if there is any.
func (c *localDisks) Create(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.CreateOptions) (result *v1alpha1.LocalDisk, err error) {
	result = &v1alpha1.LocalDisk{}
	err = c.client.Post().
		Resource("localdisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(localDisk).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a localDisk and updates it. Returns the server's representation of the localDisk, and an error, if there is any.
func (c *localDisks) Update(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.UpdateOptions) (result *v1alpha1.LocalDisk, err error) {
	result = &v1alpha1.LocalDisk{}
	err = c.client.Put().
		Resource("localdisks").
		Name(localDisk.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(localDisk).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *localDisks) UpdateStatus(ctx context.Context, localDisk *v1alpha1.LocalDisk, opts v1.UpdateOptions) (result *v1alpha1.LocalDisk, err error) {
	result = &v1alpha1.LocalDisk{}
	err = c.client.Put().
		Resource("localdisks").
		Name(localDisk.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(localDisk).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the localDisk and deletes it. Returns an error if one occurs.
func (c *localDisks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("localdisks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *localDisks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("localdisks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched localDisk.
func (c *localDisks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LocalDisk, err error) {
	result = &v1alpha1.LocalDisk{}
	err = c.client.Patch(pt).
		Resource("localdisks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
