package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	datastoreclientsetv1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/typed/datastore/v1alpha1"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/ftp"
	"github.com/hwameistor/datastore/pkg/storage/minio"
	"github.com/hwameistor/datastore/pkg/storage/nfs"
	"github.com/hwameistor/datastore/pkg/storage/ssh"
	"github.com/hwameistor/datastore/pkg/storage/web"
)

func BuildInClusterDataStoreClientset() *datastoreclientsetv1alpha1.DatastoreV1alpha1Client {
	config, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return datastoreclientsetv1alpha1.NewForConfigOrDie(config)
}

// BuildInClusterClientset builds a kubernetes in-cluster clientset
func BuildInClusterClientset() *kubernetes.Clientset {
	config, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return kubernetes.NewForConfigOrDie(config)
}

func DownloadFileByMinIO(spec *datastorev1alpha1.MinIOSpec, objName string, localFilePath string) error {
	minioClient, err := minio.NewClient(spec)
	if err != nil {
		log.WithField("endpoint", spec.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}
	return minio.DownloadObject(minioClient, spec.Bucket, fmt.Sprintf("%s/%s", spec.Prefix, objName), localFilePath)

}

func UploadFileByMinIO(spec *datastorev1alpha1.MinIOSpec, objName string, localFilePath string) error {
	minioClient, err := minio.NewClient(spec)
	if err != nil {
		log.WithField("endpoint", spec.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}
	return minio.UploadObject(minioClient, spec.Bucket, fmt.Sprintf("%s/%s", spec.Prefix, objName), localFilePath)
}

func DownloadFileByHttp(spec *datastorev1alpha1.HTTPSpec, fName string, localFilePath string) error {
	return web.DownloadObject(spec, fName, localFilePath)
}

func DownloadFileByNFS(spec *datastorev1alpha1.NFSSpec, fName string, localFilePath string) error {
	return nfs.DownloadObject(spec, fName, localFilePath)

}

func DownloadFileByFTP(spec *datastorev1alpha1.FTPSpec, fName string, localFilePath string) error {
	return ftp.DownloadObject(spec, fName, localFilePath)
}

func DownloadFileBySSH(node string, remoteFilePath string, localFilePath string) error {
	return ssh.DownloadObject(node, remoteFilePath, localFilePath)
}

func Checksum(fpath string) (string, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Mkdir(dirpath string) error {
	f, err := os.Stat(dirpath)
	if err != nil {
		return os.MkdirAll(dirpath, os.ModePerm)
	}
	if !f.IsDir() {
		return os.MkdirAll(dirpath, os.ModePerm)
	}
	return nil
}
