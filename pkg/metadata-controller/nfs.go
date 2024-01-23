package metadatacontroller

import (
	"fmt"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
)

func (mgr *storageBackendManager) _checkConnectionForNFS(backend *datastorev1alpha1.StorageBackend) (bool, error) {
	logCtx := log.WithFields(log.Fields{"backend": backend.Name, "endpoint": backend.Spec.NFS.Endpoint})
	logCtx.Debug("Checking connection for a NFS storage backend")
	if backend.Spec.NFS == nil {
		return false, fmt.Errorf("invaild NFS spec info")
	}

	mountClient, err := nfs.DialMount(backend.Spec.NFS.Endpoint)
	if err != nil {
		logCtx.WithError(err).Error("Failed to setup the mount client")
		return false, err
	}
	defer mountClient.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mountClient.Mount(backend.Spec.NFS.RootDir, auth.Auth())
	if err != nil {
		logCtx.WithError(err).Error("Failed to mount")
		return false, err
	}
	defer target.Close()

	if err = mountClient.Unmount(); err != nil {
		logCtx.WithError(err).Error("Failed to unmount")
	}

	return true, nil
}

func (mgr *storageBackendManager) handleStorageBackendForNFS(backend *datastorev1alpha1.StorageBackend) error {
	logCtx := log.WithFields(log.Fields{"backend": backend.Name, "endpoint": backend.Spec.NFS.Endpoint, "target": backend.Spec.NFS.RootDir})
	logCtx.Debug("Handling a NFS storage backend")

	mountClient, err := nfs.DialMount(backend.Spec.NFS.Endpoint)
	if err != nil {
		logCtx.WithError(err).Error("Failed to setup the mount client")
		return err
	}
	defer mountClient.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mountClient.Mount(backend.Spec.NFS.RootDir, auth.Auth())
	if err != nil {
		logCtx.WithError(err).Error("Failed to mount")
		return err
	}
	defer target.Close()

	if err = mgr.refreshDataFromStorageBackendForNFS(target, backend); err != nil {
		logCtx.WithError(err).Error("Error happened when refreshing the data")
	}
	if err = mountClient.Unmount(); err != nil {
		logCtx.WithError(err).Error("Failed to unmount")
	}

	return err
}

func (mgr *storageBackendManager) refreshDataFromStorageBackendForNFS(nfsTarget *nfs.Target, backend *datastorev1alpha1.StorageBackend) error {
	return nil
}
