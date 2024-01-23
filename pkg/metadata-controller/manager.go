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

type storageBackendManager struct {
	dsClientSet *datastoreclientset.Clientset

	globalView *GlobalViewFileSystem
}

func newstorageBackendManager() *storageBackendManager {
	return &storageBackendManager{
		globalView: &GlobalViewFileSystem{},
	}
}

func (mgr *storageBackendManager) Run(stopCh <-chan struct{}) {

	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to get kubernetes cluster config")
	}

	log.Debug("start datastore informer factory")
	mgr.dsClientSet = datastoreclientset.NewForConfigOrDie(cfg)
	dsFactory := datastoreinformers.NewSharedInformerFactory(mgr.dsClientSet, 0)
	dsFactory.Start(stopCh)

	sbInformer := dsFactory.Datastore().V1alpha1().StorageBackends()
	sbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mgr.onStorageBackendAdd,
		UpdateFunc: mgr.onStorageBackendUpdate,
		DeleteFunc: mgr.onStorageBackendDelete,
	})
	go sbInformer.Informer().Run(stopCh)

	go mgr.checkConnectionForever(stopCh)
}

func (mgr *storageBackendManager) checkConnectionForever(stopCh <-chan struct{}) {
	log.Debug("Starting a worker to check storage backends for the connection regularly")
	mgr.checkConnection()
	for {
		select {
		case <-time.After(datastorev1alpha1.StorageBackendConnectionCheckInterval):
			mgr.checkConnection()
		case <-stopCh:
			log.Debug("Exit the node status synchronizing")
			return
		}
	}
}

func (mgr *storageBackendManager) checkConnection() {
	log.Debug("Starting to check storage backends' connection")

	ctx := context.Background()
	backendList, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to list the storage backends")
		return
	}
	for i := range backendList.Items {
		backend := &backendList.Items[i]
		if connected, err := mgr._checkConnection(backend); err != nil {
			if backend.Status.Error != err.Error() {
				backend.Status.Error = err.Error()
				if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
					log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
				}
			}
		} else {
			if backend.Status.Connected != connected {
				backend.Status.Connected = connected
				if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
					log.WithField("Connected", backend.Status.Connected).WithError(err).Error("Failed to update storage backend")
				}
			}
		}
		log.WithFields(log.Fields{"backend": backend.Name}).Debug("Checking completed")

	}

}

func (mgr *storageBackendManager) _checkConnection(backend *datastorev1alpha1.StorageBackend) (bool, error) {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		return mgr._checkConnectionForMinIO(backend)
	}
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		return mgr._checkConnectionForNFS(backend)
	}

	return false, fmt.Errorf("unknown storage backend")
}

func (mgr *storageBackendManager) onStorageBackendAdd(obj interface{}) {
	backend, _ := obj.(*datastorev1alpha1.StorageBackend)
	mgr.handleStorageBackend(backend, true)
}

func (mgr *storageBackendManager) onStorageBackendUpdate(oldObj, newObj interface{}) {
	backend, _ := newObj.(*datastorev1alpha1.StorageBackend)
	mgr.handleStorageBackend(backend, false)
}

func (mgr *storageBackendManager) onStorageBackendDelete(obj interface{}) {
	backend, _ := obj.(*datastorev1alpha1.StorageBackend)
	mgr.removeStorageBackend(backend)
}

func (mgr *storageBackendManager) handleStorageBackend(backend *datastorev1alpha1.StorageBackend, forceRefresh bool) {
	mgr.globalView.UpdateDataServer(backend)

	if !forceRefresh && !backend.Spec.Refresh {
		return
	}

	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		mgr.handleStorageBackendForMinIO(backend)
	} else if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		mgr.handleStorageBackendForNFS(backend)
	} else {
		mgr.handleStorageBackendForUnknown(backend)
	}

	ctx := context.Background()
	backend.Spec.Refresh = false
	newBackend, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().Update(ctx, backend, metav1.UpdateOptions{})
	if err != nil {
		log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
	}
	newBackend.Status.LastRefreshTimestamp = &metav1.Time{Time: time.Now()}
	newBackend.Status.RefreshCount++
	newBackend.Status.Connected = true
	if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, newBackend, metav1.UpdateOptions{}); err != nil {
		log.WithField("Message", newBackend.Status.Error).WithError(err).Error("Failed to update storage backend")
	}

	dumpGlobalView(mgr.globalView)
}

func (mgr *storageBackendManager) handleStorageBackendForUnknown(backend *datastorev1alpha1.StorageBackend) {
	log.WithFields(log.Fields{"backend": backend.Name}).Debug("Handling a unknown storage backend")

}

func (mgr *storageBackendManager) removeStorageBackend(backend *datastorev1alpha1.StorageBackend) {
	mgr.globalView.RemoveDataServer(backend)

	dumpGlobalView(mgr.globalView)
}
