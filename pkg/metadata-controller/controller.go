package metadatacontroller

import (
	log "github.com/sirupsen/logrus"
)

// MetadataController interface
type MetadataController interface {
	Run(stopCh <-chan struct{}) error
}

type controller struct {
	dsMgr *dataSourceManager
}

// New an assistant instance
func NewMetadataController(namespace string) MetadataController {
	return &controller{
		dsMgr: newdataSourceManager(namespace),
	}
}

func (ctrl *controller) Run(stopCh <-chan struct{}) error {
	log.Debug("Start metadata controller ...")

	ctrl.dsMgr.Run(stopCh)

	return nil
}
