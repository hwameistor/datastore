// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/scheme"
	v1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BaseModelsGetter has a method to return a BaseModelInterface.
// A group's client should implement this interface.
type BaseModelsGetter interface {
	BaseModels(namespace string) BaseModelInterface
}

// BaseModelInterface has methods to work with BaseModel resources.
type BaseModelInterface interface {
	Create(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.CreateOptions) (*v1alpha1.BaseModel, error)
	Update(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (*v1alpha1.BaseModel, error)
	UpdateStatus(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (*v1alpha1.BaseModel, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.BaseModel, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.BaseModelList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BaseModel, err error)
	BaseModelExpansion
}

// baseModels implements BaseModelInterface
type baseModels struct {
	client rest.Interface
	ns     string
}

// newBaseModels returns a BaseModels
func newBaseModels(c *DatastoreV1alpha1Client, namespace string) *baseModels {
	return &baseModels{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the baseModel, and returns the corresponding baseModel object, and an error if there is any.
func (c *baseModels) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.BaseModel, err error) {
	result = &v1alpha1.BaseModel{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("basemodels").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BaseModels that match those selectors.
func (c *baseModels) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BaseModelList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BaseModelList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("basemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested baseModels.
func (c *baseModels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("basemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a baseModel and creates it.  Returns the server's representation of the baseModel, and an error, if there is any.
func (c *baseModels) Create(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.CreateOptions) (result *v1alpha1.BaseModel, err error) {
	result = &v1alpha1.BaseModel{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("basemodels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(baseModel).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a baseModel and updates it. Returns the server's representation of the baseModel, and an error, if there is any.
func (c *baseModels) Update(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (result *v1alpha1.BaseModel, err error) {
	result = &v1alpha1.BaseModel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("basemodels").
		Name(baseModel.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(baseModel).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *baseModels) UpdateStatus(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (result *v1alpha1.BaseModel, err error) {
	result = &v1alpha1.BaseModel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("basemodels").
		Name(baseModel.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(baseModel).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the baseModel and deletes it. Returns an error if one occurs.
func (c *baseModels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("basemodels").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *baseModels) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("basemodels").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched baseModel.
func (c *baseModels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BaseModel, err error) {
	result = &v1alpha1.BaseModel{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("basemodels").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
