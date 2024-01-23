package metadatacontroller

import (
	"fmt"
	"strings"
	"time"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

type GlobalViewFileSystem struct {
	Servers map[string]*DataServer
}

type DataServer struct {
	Type      string
	Connected bool
	Endpoint  string

	Dirs map[string]*DataDirectory

	AllDirs map[string]*DataDirectory
	AllObjs map[string]*DataObject
}

type DataDirectory struct {
	Name   string
	Path   string
	Prefix string // s3 bucket only

	SubDirs map[string]*DataDirectory
	Objects map[string]*DataObject
}

type DataObject struct {
	Name  string
	Path  string
	Size  int64
	MTime time.Time
}

func (gs *GlobalViewFileSystem) UpdateDataServer(backend *datastorev1alpha1.StorageBackend) {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		gs._updateDataServerForMinIO(backend)
	}

}

func (gs *GlobalViewFileSystem) _updateDataServerForMinIO(backend *datastorev1alpha1.StorageBackend) {
	if len(gs.Servers) == 0 {
		return
	}
	if gs.Servers[backend.Name] == nil {
		return
	}
	gs.Servers[backend.Name].Endpoint = backend.Spec.MinIO.Endpoint
	gs.Servers[backend.Name].Connected = backend.Status.Connected
}

func (gs *GlobalViewFileSystem) resetDataServer(backend *datastorev1alpha1.StorageBackend) *DataServer {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		return gs._resetDataServerForMinIO(backend)
	}

	return nil
}

func (gs *GlobalViewFileSystem) _resetDataServerForMinIO(backend *datastorev1alpha1.StorageBackend) *DataServer {
	if len(gs.Servers) == 0 {
		gs.Servers = map[string]*DataServer{}
	}
	if gs.Servers[backend.Name] == nil {
		gs.Servers[backend.Name] = &DataServer{
			Type:      backend.Spec.Type,
			Connected: backend.Status.Connected,
			Endpoint:  backend.Spec.MinIO.Endpoint,
			Dirs:      map[string]*DataDirectory{},
			AllDirs:   map[string]*DataDirectory{},
			AllObjs:   map[string]*DataObject{},
		}
	} else {
		gs.Servers[backend.Name].Endpoint = backend.Spec.MinIO.Endpoint
		gs.Servers[backend.Name].Connected = backend.Status.Connected
		gs.Servers[backend.Name].Dirs = map[string]*DataDirectory{}
		gs.Servers[backend.Name].AllDirs = map[string]*DataDirectory{}
		gs.Servers[backend.Name].AllObjs = map[string]*DataObject{}
	}

	minioInfo := backend.Spec.MinIO
	bucket := DataDirectory{
		Name:    minioInfo.Bucket,
		Prefix:  minioInfo.Prefix,
		Path:    minioInfo.Bucket + "/",
		SubDirs: map[string]*DataDirectory{},
		Objects: map[string]*DataObject{},
	}
	gs.Servers[backend.Name].Dirs[minioInfo.Bucket] = &bucket

	return gs.Servers[backend.Name]
}

func (gs *GlobalViewFileSystem) UpdateDataObjects(backend *datastorev1alpha1.StorageBackend, objs []*DataObject) {
	if backend.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		gs._updateDataObjectsForMinIO(backend, objs)
	}
}

func (gs *GlobalViewFileSystem) _updateDataObjectsForMinIO(backend *datastorev1alpha1.StorageBackend, objs []*DataObject) {
	server := gs.resetDataServer(backend)
	bucket := server.Dirs[backend.Spec.MinIO.Bucket]

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
					upDir.SubDirs[dirPath] = &DataDirectory{
						Name:    items[pos],
						Path:    dirPath,
						SubDirs: map[string]*DataDirectory{},
						Objects: map[string]*DataObject{},
					}
				}
				upDir = upDir.SubDirs[dirPath]
				pos++
			}
			upDir.Objects[obj.Path] = obj
		}
	}
}

func dumpGlobalView(gv *GlobalViewFileSystem) {
	for name, server := range gv.Servers {
		dumpServer(name, server)
	}
}

func dumpServer(name string, svr *DataServer) {
	fmt.Printf("Server:  Backend: %s, Type: %s, Endpoint: %s, Connected: %t\n", name, svr.Type, svr.Endpoint, svr.Connected)
	for _, dir := range svr.Dirs {
		fmt.Printf("      Dir: Path: %s, Type: Bucket\n", dir.Path)
	}

	for _, dir := range svr.Dirs {
		dumpDir(dir)
	}
}

func dumpDir(dir *DataDirectory) {
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
