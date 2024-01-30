package metadatacontroller

import (
	"fmt"
	"strings"
	"sync"

	datastoreapis "github.com/hwameistor/datastore/pkg/apis"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type GlobalViewFileSystem struct {
	Servers map[string]*datastoreapis.DataServer

	lock sync.Mutex
}

func (gs *GlobalViewFileSystem) UpdateDataServer(backend *datastorev1alpha1.StorageBackend) {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		gs._updateDataServerForMinIO(backend)
	}
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		gs._updateDataServerForNFS(backend)
	}
}

func (gs *GlobalViewFileSystem) _updateDataServerForMinIO(backend *datastorev1alpha1.StorageBackend) {
	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[backend.Name] == nil {
		return
	}

	gs.lock.Lock()

	gs.Servers[backend.Name].Endpoint = backend.Spec.MinIO.Endpoint
	gs.Servers[backend.Name].Connected = backend.Status.Connected

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) _updateDataServerForNFS(backend *datastorev1alpha1.StorageBackend) {
	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[backend.Name] == nil {
		return
	}

	gs.lock.Lock()

	gs.Servers[backend.Name].Endpoint = backend.Spec.NFS.Endpoint
	gs.Servers[backend.Name].Connected = backend.Status.Connected

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) RemoveDataServer(backend *datastorev1alpha1.StorageBackend) {
	log.WithFields(log.Fields{"backend": backend.Name}).Debug("Removing a storage backend ...")

	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[backend.Name] == nil {
		return
	}

	gs.lock.Lock()

	delete(gs.Servers, backend.Name)

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) resetDataServer(backend *datastorev1alpha1.StorageBackend) *datastoreapis.DataServer {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		return gs._resetDataServerForMinIO(backend)
	}
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		return gs._resetDataServerForNFS(backend)
	}

	return nil
}

func (gs *GlobalViewFileSystem) _resetDataServerForMinIO(backend *datastorev1alpha1.StorageBackend) *datastoreapis.DataServer {

	if len(gs.Servers) == 0 {
		gs.Servers = map[string]*datastoreapis.DataServer{}
	}
	if gs.Servers[backend.Name] == nil {
		gs.Servers[backend.Name] = &datastoreapis.DataServer{
			Type:      backend.Spec.Type,
			Connected: backend.Status.Connected,
			Endpoint:  backend.Spec.MinIO.Endpoint,
			SubDirs:   map[string]*datastoreapis.DataDirectory{},
			AllDirs:   map[string]*datastoreapis.DataDirectory{},
			AllObjs:   map[string]*datastoreapis.DataObject{},
		}
	} else {
		gs.Servers[backend.Name].Endpoint = backend.Spec.MinIO.Endpoint
		gs.Servers[backend.Name].Connected = backend.Status.Connected
		gs.Servers[backend.Name].SubDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[backend.Name].AllDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[backend.Name].AllObjs = map[string]*datastoreapis.DataObject{}
	}

	minioInfo := backend.Spec.MinIO
	bucket := datastoreapis.DataDirectory{
		Name:    minioInfo.Bucket,
		Prefix:  minioInfo.Prefix,
		Path:    minioInfo.Bucket + "/",
		SubDirs: map[string]*datastoreapis.DataDirectory{},
		Objects: map[string]*datastoreapis.DataObject{},
	}
	gs.Servers[backend.Name].SubDirs[minioInfo.Bucket] = &bucket

	return gs.Servers[backend.Name]
}

func (gs *GlobalViewFileSystem) _resetDataServerForNFS(backend *datastorev1alpha1.StorageBackend) *datastoreapis.DataServer {

	if len(gs.Servers) == 0 {
		gs.Servers = map[string]*datastoreapis.DataServer{}
	}
	if gs.Servers[backend.Name] == nil {
		gs.Servers[backend.Name] = &datastoreapis.DataServer{
			Type:      backend.Spec.Type,
			Connected: backend.Status.Connected,
			Endpoint:  backend.Spec.NFS.Endpoint,
			SubDirs:   map[string]*datastoreapis.DataDirectory{},
			AllDirs:   map[string]*datastoreapis.DataDirectory{},
			AllObjs:   map[string]*datastoreapis.DataObject{},
		}
	} else {
		gs.Servers[backend.Name].Endpoint = backend.Spec.NFS.Endpoint
		gs.Servers[backend.Name].Connected = backend.Status.Connected
		gs.Servers[backend.Name].SubDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[backend.Name].AllDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[backend.Name].AllObjs = map[string]*datastoreapis.DataObject{}
	}

	spec := backend.Spec.NFS
	rootdir := datastoreapis.DataDirectory{
		Name:    spec.RootDir,
		Path:    spec.RootDir + "/",
		SubDirs: map[string]*datastoreapis.DataDirectory{},
		Objects: map[string]*datastoreapis.DataObject{},
	}
	gs.Servers[backend.Name].SubDirs[spec.RootDir] = &rootdir

	return gs.Servers[backend.Name]
}

func (gs *GlobalViewFileSystem) UpdateDataObjects(backend *datastorev1alpha1.StorageBackend, objs []*datastoreapis.DataObject) {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		gs._updateDataObjectsForMinIO(backend, objs)
	} else if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		gs._updateDataObjectsForNFS(backend, objs)
	}
}

func (gs *GlobalViewFileSystem) _updateDataObjectsForMinIO(backend *datastorev1alpha1.StorageBackend, objs []*datastoreapis.DataObject) {

	server := gs.resetDataServer(backend)
	bucket := server.SubDirs[backend.Spec.MinIO.Bucket]

	gs.lock.Lock()
	defer gs.lock.Unlock()

	for _, obj := range objs {
		server.AllObjs[obj.Path] = obj
		items := strings.Split(obj.Path, "/")
		if len(items) == 1 {
			// no folder, e.g picture.jpg
			bucket.Objects[obj.Path] = obj
		} else {
			// e.g. words/test/kube/shell/init-kube.sh
			pos := 0
			dirPath := ""
			upDir := bucket
			for pos < len(items)-1 {
				dirPath = dirPath + items[pos] + "/"
				if upDir.SubDirs[dirPath] == nil {
					fmt.Printf("Created a directory: %s\n", dirPath)
					upDir.SubDirs[dirPath] = &datastoreapis.DataDirectory{
						Name:    items[pos],
						Path:    dirPath,
						SubDirs: map[string]*datastoreapis.DataDirectory{},
						Objects: map[string]*datastoreapis.DataObject{},
					}
				}
				upDir = upDir.SubDirs[dirPath]
				pos++
			}
			upDir.Objects[obj.Path] = obj
		}
	}
}

func (gs *GlobalViewFileSystem) _updateDataObjectsForNFS(backend *datastorev1alpha1.StorageBackend, objs []*datastoreapis.DataObject) {

	server := gs.resetDataServer(backend)
	rootdir := server.SubDirs[backend.Spec.NFS.RootDir]

	gs.lock.Lock()
	defer gs.lock.Unlock()

	for _, obj := range objs {
		server.AllObjs[obj.Path] = obj
		log.WithField("file", obj.Path).Debug("Updating a data object ...")
		items := strings.Split(obj.Path, "/")
		if len(items) == 1 {
			// no folder, e.g picture.jpg
			rootdir.Objects[obj.Path] = obj
		} else {
			// e.g. ./words/test/kube/shell/init-kube.sh
			pos := 0
			dirPath := ""
			upDir := rootdir
			for pos < len(items)-1 {
				dirPath = dirPath + items[pos] + "/"
				if upDir.SubDirs[dirPath] == nil {
					fmt.Printf("Created a directory: %s\n", dirPath)
					upDir.SubDirs[dirPath] = &datastoreapis.DataDirectory{
						Name:    items[pos],
						Path:    dirPath,
						SubDirs: map[string]*datastoreapis.DataDirectory{},
						Objects: map[string]*datastoreapis.DataObject{},
					}
				}
				upDir = upDir.SubDirs[dirPath]
				pos++
			}
			upDir.Objects[obj.Path] = obj
		}
	}
	log.WithField("files", len(objs)).Debug("Update completed")
}

func dumpGlobalView(gv *GlobalViewFileSystem) {
	for name, server := range gv.Servers {
		dumpServer(name, server)
	}
}

func dumpServer(name string, svr *datastoreapis.DataServer) {
	fmt.Printf("Server:  Backend: %s, Type: %s, Endpoint: %s, Connected: %t\n", name, svr.Type, svr.Endpoint, svr.Connected)
	for _, dir := range svr.SubDirs {
		fmt.Printf("      Dir: Path: %s, Type: Root\n", dir.Path)
	}

	for _, dir := range svr.SubDirs {
		dumpDir(dir)
	}
}

func dumpDir(dir *datastoreapis.DataDirectory) {
	fmt.Printf("Directory: Path: %s\n", dir.Path)
	for _, dir := range dir.SubDirs {
		fmt.Printf("      Dir: Path: %s\n", dir.Path)
	}

	for _, obj := range dir.Objects {
		fmt.Printf("   Object: Path: %s, Size: %d\n", obj.Path, obj.Size)
	}

	for _, dir := range dir.SubDirs {
		dumpDir(dir)
	}
}
