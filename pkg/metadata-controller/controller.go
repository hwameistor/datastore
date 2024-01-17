package metadatacontroller

import (
	"github.com/hwameistor/datastore/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

// It has some data to be shared among all the controllers.
// So, it's a global variable
var instance MetadataController

func Instance() MetadataController {
	if instance == nil {
		instance = newMetadataController()
	}
	return instance
}

// MetadataController interface
type MetadataController interface {
	Run(stopCh <-chan struct{}) error

	ReconcileStorageBackend(backend *datastorev1alpha1.StorageBackend)
}

type controller struct {
	clientset *kubernetes.Clientset
}

// New an assistant instance
func newMetadataController() MetadataController {
	return &controller{
		clientset: utils.BuildInClusterClientset(),
	}
}

func (ctrl *controller) Run(stopCh <-chan struct{}) error {
	log.Debug("start informer factory")

	return nil
}

func (ctrl *controller) ReconcileStorageBackend(backend *datastorev1alpha1.StorageBackend) {
	log.Debug("Reconciling a StorageBackend")

}
