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

func (gs *GlobalViewFileSystem) UpdateDataServer(ds *datastorev1alpha1.DataSet) {
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeMinIO {
		gs._updateDataServerForMinIO(ds)
	}
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeNFS {
		gs._updateDataServerForNFS(ds)
	}
}

func (gs *GlobalViewFileSystem) _updateDataServerForMinIO(ds *datastorev1alpha1.DataSet) {
	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[ds.Name] == nil {
		return
	}

	gs.lock.Lock()

	gs.Servers[ds.Name].Endpoint = ds.Spec.MinIO.Endpoint
	gs.Servers[ds.Name].Connected = ds.Status.Connected

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) _updateDataServerForNFS(ds *datastorev1alpha1.DataSet) {
	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[ds.Name] == nil {
		return
	}

	gs.lock.Lock()

	gs.Servers[ds.Name].Endpoint = ds.Spec.NFS.Endpoint
	gs.Servers[ds.Name].Connected = ds.Status.Connected

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) RemoveDataServer(ds *datastorev1alpha1.DataSet) {
	log.WithFields(log.Fields{"ds": ds.Name}).Debug("Removing a storage ds ...")

	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[ds.Name] == nil {
		return
	}

	gs.lock.Lock()

	delete(gs.Servers, ds.Name)

	gs.lock.Unlock()
}

func (gs *GlobalViewFileSystem) resetDataServer(ds *datastorev1alpha1.DataSet) *datastoreapis.DataServer {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeMinIO {
		return gs._resetDataServerForMinIO(ds)
	}
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeNFS {
		return gs._resetDataServerForNFS(ds)
	}

	return nil
}

func (gs *GlobalViewFileSystem) _resetDataServerForMinIO(ds *datastorev1alpha1.DataSet) *datastoreapis.DataServer {

	if len(gs.Servers) == 0 {
		gs.Servers = map[string]*datastoreapis.DataServer{}
	}
	if gs.Servers[ds.Name] == nil {
		gs.Servers[ds.Name] = &datastoreapis.DataServer{
			Type:      ds.Spec.Type,
			Connected: ds.Status.Connected,
			Endpoint:  ds.Spec.MinIO.Endpoint,
			SubDirs:   map[string]*datastoreapis.DataDirectory{},
			AllDirs:   map[string]*datastoreapis.DataDirectory{},
			AllObjs:   map[string]*datastoreapis.DataObject{},
		}
	} else {
		gs.Servers[ds.Name].Endpoint = ds.Spec.MinIO.Endpoint
		gs.Servers[ds.Name].Connected = ds.Status.Connected
		gs.Servers[ds.Name].SubDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[ds.Name].AllDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[ds.Name].AllObjs = map[string]*datastoreapis.DataObject{}
	}

	minioInfo := ds.Spec.MinIO
	bucket := datastoreapis.DataDirectory{
		Name:    minioInfo.Bucket,
		Prefix:  minioInfo.Prefix,
		Path:    minioInfo.Bucket + "/",
		SubDirs: map[string]*datastoreapis.DataDirectory{},
		Objects: map[string]*datastoreapis.DataObject{},
	}
	gs.Servers[ds.Name].SubDirs[minioInfo.Bucket] = &bucket

	return gs.Servers[ds.Name]
}

func (gs *GlobalViewFileSystem) _resetDataServerForNFS(ds *datastorev1alpha1.DataSet) *datastoreapis.DataServer {

	if len(gs.Servers) == 0 {
		gs.Servers = map[string]*datastoreapis.DataServer{}
	}
	if gs.Servers[ds.Name] == nil {
		gs.Servers[ds.Name] = &datastoreapis.DataServer{
			Type:      ds.Spec.Type,
			Connected: ds.Status.Connected,
			Endpoint:  ds.Spec.NFS.Endpoint,
			SubDirs:   map[string]*datastoreapis.DataDirectory{},
			AllDirs:   map[string]*datastoreapis.DataDirectory{},
			AllObjs:   map[string]*datastoreapis.DataObject{},
		}
	} else {
		gs.Servers[ds.Name].Endpoint = ds.Spec.NFS.Endpoint
		gs.Servers[ds.Name].Connected = ds.Status.Connected
		gs.Servers[ds.Name].SubDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[ds.Name].AllDirs = map[string]*datastoreapis.DataDirectory{}
		gs.Servers[ds.Name].AllObjs = map[string]*datastoreapis.DataObject{}
	}

	spec := ds.Spec.NFS
	rootdir := datastoreapis.DataDirectory{
		Name:    spec.RootDir,
		Path:    spec.RootDir + "/",
		SubDirs: map[string]*datastoreapis.DataDirectory{},
		Objects: map[string]*datastoreapis.DataObject{},
	}
	gs.Servers[ds.Name].SubDirs[spec.RootDir] = &rootdir

	return gs.Servers[ds.Name]
}

func (gs *GlobalViewFileSystem) UpdateDataObjects(ds *datastorev1alpha1.DataSet, objs []*datastoreapis.DataObject) {
	if ds.Spec.Type == datastorev1alpha1.DataSourceTypeMinIO {
		gs._updateDataObjectsForMinIO(ds, objs)
	} else if ds.Spec.Type == datastorev1alpha1.DataSourceTypeNFS {
		gs._updateDataObjectsForNFS(ds, objs)
	}
}

func (gs *GlobalViewFileSystem) _updateDataObjectsForMinIO(ds *datastorev1alpha1.DataSet, objs []*datastoreapis.DataObject) {

	server := gs.resetDataServer(ds)
	bucket := server.SubDirs[ds.Spec.MinIO.Bucket]

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

func (gs *GlobalViewFileSystem) _updateDataObjectsForNFS(ds *datastorev1alpha1.DataSet, objs []*datastoreapis.DataObject) {

	server := gs.resetDataServer(ds)
	rootdir := server.SubDirs[ds.Spec.NFS.RootDir]

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
