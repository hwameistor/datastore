package nfs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	datastoreapis "github.com/hwameistor/datastore/pkg/apis"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/willscott/go-nfs-client/nfs"
	"github.com/willscott/go-nfs-client/nfs/rpc"
)

func IsConnected(spec *datastorev1alpha1.NFSSpec) (bool, error) {

	mount, err := nfs.DialMount(spec.Endpoint, time.Second)
	if err != nil {
		return false, err
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mount.Mount(spec.Export, auth.Auth())
	if err != nil {
		return false, err
	}
	defer target.Close()

	mount.Unmount()

	return true, nil

}

func LoadObjectMetadata(spec *datastorev1alpha1.NFSSpec) ([]*datastoreapis.DataObject, error) {

	files := []*datastoreapis.DataObject{}

	mount, err := nfs.DialMount(spec.Endpoint, time.Second)
	if err != nil {
		return files, err
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mount.Mount(spec.Export, auth.Auth())
	if err != nil {
		return files, err
	}
	defer target.Close()

	dirs := []string{spec.RootDir}
	for len(dirs) > 0 {
		subdirs := []string{}
		for _, dirpath := range dirs {
			objs, err := target.ReadDirPlus(dirpath)
			if err != nil {
				return files, err
			}
			for _, obj := range objs {
				path := strings.TrimPrefix(fmt.Sprintf("%s/%s", dirpath, obj.FileName), "./")
				if obj.IsDir() {
					fmt.Printf("Directory:  name: %s, size: %d\n", path, obj.Size())
					subdirs = append(subdirs, path)
				} else {
					file := datastoreapis.DataObject{Name: obj.FileName, Path: path, Size: obj.Size(), MTime: obj.ModTime()}
					fmt.Printf("     File:  %v+\n", file)
					files = append(files, &file)
				}
			}
		}
		dirs = subdirs
	}

	mount.Unmount()

	return files, nil
}

func LoadObjectToLocal(spec *datastorev1alpha1.NFSSpec, localRootDir string) error {

	mount, err := nfs.DialMount(spec.Endpoint, time.Second)
	if err != nil {
		return err
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("hasselhoff", 1001, 1001)
	target, err := mount.Mount(spec.Export, auth.Auth())
	if err != nil {
		return err
	}
	defer target.Close()

	dirs := []string{spec.RootDir}
	for len(dirs) > 0 {
		subdirs := []string{}
		for _, dirpath := range dirs {
			objs, err := target.ReadDirPlus(dirpath)
			if err != nil {
				return err
			}
			for _, obj := range objs {
				path := strings.TrimPrefix(fmt.Sprintf("%s/%s", dirpath, obj.FileName), "./")
				if obj.IsDir() {
					//fmt.Printf("Directory:  name: %s, size: %d\n", path, obj.Size())
					subdirs = append(subdirs, path)
				} else {
					localDirPath := localRootDir + "/" + dirpath
					if _, err = target.Mkdir(localDirPath, 0775); err != nil {
						return fmt.Errorf("failed to create a local directory %s with error : %s", localDirPath, err.Error())
					}
					localFilePath := localDirPath + "/" + obj.FileName
					lf, err := os.Open(localFilePath)
					if err != nil {
						return fmt.Errorf("failed to create a local file %s with error : %s", localFilePath, err.Error())
					}
					rf, err := target.OpenFile(path, 0777)
					if err != nil {
						return fmt.Errorf("failed to open a nfs file %s with error : %s", path, err.Error())
					}
					n, err := io.CopyN(rf, lf, obj.Size())
					if err != nil {
						return fmt.Errorf("failed to load data into local file %s (n=%d) with error : %s", localFilePath, n, err.Error())
					}
				}
			}
		}
		dirs = subdirs
	}

	mount.Unmount()

	return nil
}
