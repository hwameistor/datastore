// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBaseModels implements BaseModelInterface
type FakeBaseModels struct {
	Fake *FakeDatastoreV1alpha1
	ns   string
}

var basemodelsResource = schema.GroupVersionResource{Group: "datastore.io", Version: "v1alpha1", Resource: "basemodels"}

var basemodelsKind = schema.GroupVersionKind{Group: "datastore.io", Version: "v1alpha1", Kind: "BaseModel"}

// Get takes name of the baseModel, and returns the corresponding baseModel object, and an error if there is any.
func (c *FakeBaseModels) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.BaseModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(basemodelsResource, c.ns, name), &v1alpha1.BaseModel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BaseModel), err
}

// List takes label and field selectors, and returns the list of BaseModels that match those selectors.
func (c *FakeBaseModels) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BaseModelList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(basemodelsResource, basemodelsKind, c.ns, opts), &v1alpha1.BaseModelList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BaseModelList{ListMeta: obj.(*v1alpha1.BaseModelList).ListMeta}
	for _, item := range obj.(*v1alpha1.BaseModelList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested baseModels.
func (c *FakeBaseModels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(basemodelsResource, c.ns, opts))

}

// Create takes the representation of a baseModel and creates it.  Returns the server's representation of the baseModel, and an error, if there is any.
func (c *FakeBaseModels) Create(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.CreateOptions) (result *v1alpha1.BaseModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(basemodelsResource, c.ns, baseModel), &v1alpha1.BaseModel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BaseModel), err
}

// Update takes the representation of a baseModel and updates it. Returns the server's representation of the baseModel, and an error, if there is any.
func (c *FakeBaseModels) Update(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (result *v1alpha1.BaseModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(basemodelsResource, c.ns, baseModel), &v1alpha1.BaseModel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BaseModel), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBaseModels) UpdateStatus(ctx context.Context, baseModel *v1alpha1.BaseModel, opts v1.UpdateOptions) (*v1alpha1.BaseModel, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(basemodelsResource, "status", c.ns, baseModel), &v1alpha1.BaseModel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BaseModel), err
}

// Delete takes name of the baseModel and deletes it. Returns an error if one occurs.
func (c *FakeBaseModels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(basemodelsResource, c.ns, name), &v1alpha1.BaseModel{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBaseModels) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(basemodelsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.BaseModelList{})
	return err
}

// Patch applies the patch and returns the patched baseModel.
func (c *FakeBaseModels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BaseModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(basemodelsResource, c.ns, name, pt, data, subresources...), &v1alpha1.BaseModel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BaseModel), err
}