package metadatacontroller

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// Evictor interface
type MetadataController interface {
	Run(stopCh <-chan struct{}) error
}

type controller struct {
	clientset *kubernetes.Clientset
}

/* steps:
1. watch for Pod update event, insert the Pod with Evicted status into evictedPodQueue;
2. pick up a Pod from evictedPodQueue, check if it is using HwameiStor volume. If yes, insert it into the migrateVolumeQueue; if not, ignore
3. pick up a volume form migrateVolumeQueue, and migrate it. Make sure there is no replica located at the node where the pod is evicted;
*/

// New an assistant instance
func New(clientset *kubernetes.Clientset) MetadataController {
	return &controller{
		clientset: clientset,
	}
}

func (ctrl *controller) Run(stopCh <-chan struct{}) error {
	log.Debug("start informer factory")

	return nil
}
