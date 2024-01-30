package dataloader

import (
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/nfs"

	log "github.com/sirupsen/logrus"
)

type nfsDataLoader struct {
}

func newDataLoaderForNFS() DataLoader {
	return &nfsDataLoader{}
}

func (dlr *nfsDataLoader) Load(request *datastorev1alpha1.DataLoadRequest, localRootDir string) error {

	log.WithFields(log.Fields{"request": request.Name}).Debug("Handling a NFS data loading request ...")

	return nfs.LoadObjectToLocal(request.Spec.NFS, localRootDir)

}
