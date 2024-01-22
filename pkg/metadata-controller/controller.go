package metadatacontroller

import (
	log "github.com/sirupsen/logrus"
)

// MetadataController interface
type MetadataController interface {
	Run(stopCh <-chan struct{}) error
}

type controller struct {
	backendManager *storageBackendManager
}

// New an assistant instance
func NewMetadataController() MetadataController {
	return &controller{
		backendManager: newstorageBackendManager(),
	}
}

func (ctrl *controller) Run(stopCh <-chan struct{}) error {
	log.Debug("Start metadata controller ...")

	ctrl.backendManager.Run(stopCh)

	return nil
}
