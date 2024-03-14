// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	versioned "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	internalinterfaces "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/listers/datastore/v1alpha1"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DataSourceInformer provides access to a shared informer and lister for
// DataSources.
type DataSourceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.DataSourceLister
}

type dataSourceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDataSourceInformer constructs a new informer for DataSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDataSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDataSourceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDataSourceInformer constructs a new informer for DataSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDataSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DatastoreV1alpha1().DataSources(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DatastoreV1alpha1().DataSources(namespace).Watch(context.TODO(), options)
			},
		},
		&datastorev1alpha1.DataSource{},
		resyncPeriod,
		indexers,
	)
}

func (f *dataSourceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDataSourceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *dataSourceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&datastorev1alpha1.DataSource{}, f.defaultInformer)
}

func (f *dataSourceInformer) Lister() v1alpha1.DataSourceLister {
	return v1alpha1.NewDataSourceLister(f.Informer().GetIndexer())
}