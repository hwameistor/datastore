package datamanager

import (
	"context"
	"fmt"
	"os"
	"strings"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/ftp"
	"github.com/hwameistor/datastore/pkg/storage/minio"
	"github.com/hwameistor/datastore/pkg/storage/nfs"
	"github.com/hwameistor/datastore/pkg/storage/web"
	"github.com/hwameistor/datastore/pkg/utils"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BaseModelParameters struct {
	Namespace    string
	ConfigName   string
	LocalPathDir string
}

type BaseModelManager struct {
	conf *datastorev1alpha1.BaseModelSpec

	LocalPathDir string
}

func NewBaseModelManager(params *BaseModelParameters) Manager {
	mgr := &BaseModelManager{}
	mgr.LocalPathDir = params.LocalPathDir
	mgr.conf = mgr.getConfigOrDie(params.Namespace, params.ConfigName)

	return mgr
}

func (mgr *BaseModelManager) getConfigOrDie(namespace string, configName string) *datastorev1alpha1.BaseModelSpec {
	clientset := utils.BuildInClusterDataStoreClientset()
	cfg, err := clientset.BaseModels(namespace).Get(context.Background(), configName, metav1.GetOptions{})
	if err != nil {
		log.WithField("config", configName).WithError(err).Fatal("Failed to get config for base model")
	}
	return &cfg.Spec
}

func (mgr *BaseModelManager) isReady(fpath string) bool {

	if _, err := os.Stat(fpath); err != nil {
		log.WithField("model", fpath).WithError(err).Warn("Not found the base model")
		return false
	}

	return true
}

func (mgr *BaseModelManager) Cook() error {
	if err := utils.Mkdir(mgr.LocalPathDir); err != nil {
		return err
	}

	fpath := fmt.Sprintf("%s/%s", strings.TrimRight(mgr.LocalPathDir, "/"), mgr.conf.ModelFileName)

	if mgr.isReady(fpath) {
		return nil
	}

	log.WithField("file", mgr.conf.ModelFileName).Debug("Loading base model ...")
	if mgr.conf.Proto == "http" && mgr.conf.HTTP != nil {
		return web.DownloadObject(mgr.conf.HTTP, mgr.conf.ModelFileName, fpath)
	}
	if mgr.conf.Proto == "minio" && mgr.conf.MinIO != nil {
		return minio.DownloadObject(mgr.conf.MinIO, mgr.conf.ModelFileName, fpath)
	}
	if mgr.conf.Proto == "nfs" && mgr.conf.NFS != nil {
		return nfs.DownloadObject(mgr.conf.NFS, mgr.conf.ModelFileName, fpath)
	}
	if mgr.conf.Proto == "ftp" && mgr.conf.FTP != nil {
		return ftp.DownloadObject(mgr.conf.FTP, mgr.conf.ModelFileName, fpath)
	}

	return fmt.Errorf("unsupported url")
}

func (mgr *BaseModelManager) Run(stopCh <-chan struct{}) {

	log.Debug("No action for base model management")

}
