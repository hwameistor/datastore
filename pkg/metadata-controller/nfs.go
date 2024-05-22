package metadatacontroller

import (
	"fmt"
	"sync"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/nfs"

	log "github.com/sirupsen/logrus"
)

var nfsLock sync.Mutex

func (mgr *dataSourceManager) _checkConnectionForNFS(backend *datastorev1alpha1.DataSet) (bool, error) {
	if backend.Spec.NFS == nil {
		return false, fmt.Errorf("invaild NFS spec info")
	}

	nfsLock.Lock()
	defer nfsLock.Unlock()

	return nfs.IsConnected(backend.Spec.NFS)
}

func (mgr *dataSourceManager) handleDataSourceForNFS(backend *datastorev1alpha1.DataSet) error {
	logCtx := log.WithFields(log.Fields{"backend": backend.Name, "endpoint": backend.Spec.NFS.Endpoint, "export": backend.Spec.NFS.Export})
	logCtx.Debug("Handling a NFS storage backend ...")

	nfsLock.Lock()
	defer nfsLock.Unlock()

	files, err := nfs.LoadObjectMetadata(backend.Spec.NFS)

	mgr.globalView.UpdateDataObjects(backend, files)

	return err
}
