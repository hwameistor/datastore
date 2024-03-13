package metadatacontroller

import (
	"context"
	"fmt"
	"time"

	datastoreclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	datastoreinformers "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type dataSourceManager struct {
	dsClientSet *datastoreclientset.Clientset

	namespace string

	globalView *GlobalViewFileSystem
}

func newdataSourceManager(namespace string) *dataSourceManager {
	return &dataSourceManager{
		namespace:  namespace,
		globalView: &GlobalViewFileSystem{},
	}
}

func (mgr *dataSourceManager) Run(stopCh <-chan struct{}) {

	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to get kubernetes cluster config")
	}

	log.Debug("start datastore informer factory")
	mgr.dsClientSet = datastoreclientset.NewForConfigOrDie(cfg)
	dsFactory := datastoreinformers.NewSharedInformerFactory(mgr.dsClientSet, 0)
	dsFactory.Start(stopCh)

	sbInformer := dsFactory.Datastore().V1alpha1().DataSources()
	sbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mgr.onDataSourceAdd,
		UpdateFunc: mgr.onDataSourceUpdate,
		DeleteFunc: mgr.onDataSourceDelete,
	})
	go sbInformer.Informer().Run(stopCh)

	go mgr.checkConnectionForever(stopCh)
}

func (mgr *dataSourceManager) checkConnectionForever(stopCh <-chan struct{}) {
	log.Debug("Starting a worker to check storage dss for the connection regularly")
	mgr.checkConnection()
	for {
		select {
		case <-time.After(datastorev1alpha1.DataSourceConnectionCheckInterval):
			mgr.checkConnection()
		case <-stopCh:
			log.Debug("Exit the node status synchronizing")
			return
		}
	}
}

func (mgr *dataSourceManager) checkConnection() {
	log.Debug("Starting to check storage dss' connection")

	ctx := context.Background()
	dsList, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to list the storage dss")
		return
	}
	for i := range dsList.Items {
		ds := &dsList.Items[i]
		if connected, err := mgr._checkConnection(ds); err != nil {
			if ds.Status.Error != err.Error() {
				ds.Status.Error = err.Error()
				if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, ds, metav1.UpdateOptions{}); err != nil {
					log.WithField("Message", ds.Status.Error).WithError(err).Error("ailed to update storage ds")
				}
			}
		} else {
			if ds.Status.Connected != connected {
				ds.Status.Connected = connected
				if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, ds, metav1.UpdateOptions{}); err != nil {
					log.WithField("Connected", ds.Status.Connected).WithError(err).Error("Failed to update status of storage ds")
				}
			}
		}

	}

}

func (mgr *dataSourceManager) _checkConnection(ds *datastorev1alpha1.DataSource) (bool, error) {
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeMinIO {
		return mgr._checkConnectionForMinIO(ds)
	}
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeNFS {
		return mgr._checkConnectionForNFS(ds)
	}

	return false, fmt.Errorf("unknown storage ds")
}

func (mgr *dataSourceManager) onDataSourceAdd(obj interface{}) {
	ds, _ := obj.(*datastorev1alpha1.DataSource)
	mgr.handleDataSource(ds, true)
}

func (mgr *dataSourceManager) onDataSourceUpdate(oldObj, newObj interface{}) {
	ds, _ := newObj.(*datastorev1alpha1.DataSource)
	mgr.handleDataSource(ds, false)
}

func (mgr *dataSourceManager) onDataSourceDelete(obj interface{}) {
	ds, _ := obj.(*datastorev1alpha1.DataSource)
	mgr.removeDataSource(ds)
}

func (mgr *dataSourceManager) handleDataSource(ds *datastorev1alpha1.DataSource, forceRefresh bool) {
	mgr.globalView.UpdateDataServer(ds)

	if !forceRefresh && !ds.Spec.Refresh {
		return
	}

	var err error
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeMinIO {
		err = mgr.handleDataSourceForMinIO(ds)
	} else if ds.Spec.Type == datastorev1alpha1.DataSourceTypeNFS {
		err = mgr.handleDataSourceForNFS(ds)
	} else {
		err = mgr.handleDataSourceForUnknown(ds)
	}

	ctx := context.Background()
	ds.Spec.Refresh = false
	newBackend, e := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).Update(ctx, ds, metav1.UpdateOptions{})
	if e != nil {
		log.WithField("Message", ds.Status.Error).WithError(e).Error("Failed to update storage ds")
	}
	newBackend.Status.LastRefreshTimestamp = &metav1.Time{Time: time.Now()}
	newBackend.Status.RefreshCount++
	newBackend.Status.Connected = true
	if err == nil {
		newBackend.Status.Error = ""
	} else {
		newBackend.Status.Error = err.Error()
	}
	if _, err = mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, newBackend, metav1.UpdateOptions{}); err != nil {
		log.WithField("Message", newBackend.Status.Error).WithError(err).Error("Failed to update status of storage ds")
	}

	dumpGlobalView(mgr.globalView)
}

func (mgr *dataSourceManager) handleDataSourceForUnknown(ds *datastorev1alpha1.DataSource) error {
	log.WithFields(log.Fields{"ds": ds.Name}).Debug("Handling a unknown storage ds")

	return fmt.Errorf("unknown storage ds")

}

func (mgr *dataSourceManager) removeDataSource(ds *datastorev1alpha1.DataSource) {
	mgr.globalView.RemoveDataServer(ds)

	dumpGlobalView(mgr.globalView)
}
