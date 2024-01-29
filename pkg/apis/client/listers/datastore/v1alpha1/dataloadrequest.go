// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DataLoadRequestLister helps list DataLoadRequests.
type DataLoadRequestLister interface {
	// List lists all DataLoadRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.DataLoadRequest, err error)
	// Get retrieves the DataLoadRequest from the index for a given name.
	Get(name string) (*v1alpha1.DataLoadRequest, error)
	DataLoadRequestListerExpansion
}

// dataLoadRequestLister implements the DataLoadRequestLister interface.
type dataLoadRequestLister struct {
	indexer cache.Indexer
}

// NewDataLoadRequestLister returns a new DataLoadRequestLister.
func NewDataLoadRequestLister(indexer cache.Indexer) DataLoadRequestLister {
	return &dataLoadRequestLister{indexer: indexer}
}

// List lists all DataLoadRequests in the indexer.
func (s *dataLoadRequestLister) List(selector labels.Selector) (ret []*v1alpha1.DataLoadRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DataLoadRequest))
	})
	return ret, err
}

// Get retrieves the DataLoadRequest from the index for a given name.
func (s *dataLoadRequestLister) Get(name string) (*v1alpha1.DataLoadRequest, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("dataloadrequest"), name)
	}
	return obj.(*v1alpha1.DataLoadRequest), nil
}
