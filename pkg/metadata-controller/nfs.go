package metadatacontroller

import (
	"fmt"
	"strings"
	"sync"
	"time"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"

	log "github.com/sirupsen/logrus"
	"github.com/willscott/go-nfs-client/nfs"
	"github.com/willscott/go-nfs-client/nfs/rpc"
)

var nfsLock sync.Mutex

func (mgr *storageBackendManager) _checkConnectionForNFS(backend *datastorev1alpha1.StorageBackend) (bool, error) {
	logCtx := log.WithFields(log.Fields{"backend": backend.Name, "endpoint": backend.Spec.NFS.Endpoint})
	if backend.Spec.NFS == nil {
		return false, fmt.Errorf("invaild NFS spec info")
	}

	nfsLock.Lock()
	defer nfsLock.Unlock()

	mount, err := nfs.DialMount(backend.Spec.NFS.Endpoint, time.Second)
	if err != nil {
		return false, err
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mount.Mount(backend.Spec.NFS.Export, auth.Auth())
	if err != nil {
		return false, err
	}
	defer target.Close()

	if err = mount.Unmount(); err != nil {
		logCtx.WithError(err).Warning("Failed to unmount")
	}

	return true, nil
}

func (mgr *storageBackendManager) handleStorageBackendForNFS(backend *datastorev1alpha1.StorageBackend) error {
	logCtx := log.WithFields(log.Fields{"backend": backend.Name, "endpoint": backend.Spec.NFS.Endpoint, "export": backend.Spec.NFS.Export})
	logCtx.Debug("Handling a NFS storage backend ...")

	nfsLock.Lock()
	defer nfsLock.Unlock()

	mount, err := nfs.DialMount(backend.Spec.NFS.Endpoint, time.Second)
	if err != nil {
		return err
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mount.Mount(backend.Spec.NFS.Export, auth.Auth())
	if err != nil {
		return err
	}
	defer target.Close()

	if err = mgr.refreshDataFromStorageBackendForNFS(target, backend); err != nil {
		logCtx.WithError(err).Error("Error happened when refreshing the data")
	}
	if err = mount.Unmount(); err != nil {
		logCtx.WithError(err).Error("Failed to unmount")
	}

	return err
}

func (mgr *storageBackendManager) refreshDataFromStorageBackendForNFS(nfsTarget *nfs.Target, backend *datastorev1alpha1.StorageBackend) error {
	nfsSpec := backend.Spec.NFS
	logCtx := log.WithFields(log.Fields{"endpoint": nfsSpec.Endpoint, "export": nfsSpec.Export, "rootdir": nfsSpec.RootDir})
	logCtx.Debug("Refreshing data from NFS server ...")
	files := []*DataObject{}
	dirs := []string{nfsSpec.RootDir}
	for len(dirs) > 0 {
		subdirs := []string{}
		for _, dirpath := range dirs {
			objs, err := nfsTarget.ReadDirPlus(dirpath)
			if err != nil {
				logCtx.WithField("dir", dirpath).WithError(err).Error("Failed to list directory")
				return err
			}
			for _, obj := range objs {
				path := strings.TrimPrefix(fmt.Sprintf("%s/%s", dirpath, obj.FileName), "./")
				if obj.IsDir() {
					fmt.Printf("Directory:  name: %s, size: %d\n", path, obj.Size())
					subdirs = append(subdirs, path)
				} else {
					file := DataObject{Name: obj.FileName, Path: path, Size: obj.Size(), MTime: obj.ModTime()}
					fmt.Printf("     File:  %v+\n", file)
					files = append(files, &file)
				}
			}
		}
		dirs = subdirs
	}

	logCtx.Debug("Refresh completed")

	mgr.globalView.UpdateDataObjects(backend, files)

	return nil
}
