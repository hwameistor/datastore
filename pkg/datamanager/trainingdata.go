package datamanager

import (
	"context"
	"fmt"
	"github.com/hwameistor/datastore/pkg/storage/minio"
	"os"
	"path/filepath"
	"time"

	datastoreclientsetv1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/typed/datastore/v1alpha1"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/hwameistor/datastore/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type TrainingDataParameters struct {
	Namespace          string
	ConfigName         string
	NodeName           string
	LocalPathDir       string
	LocalPathDirOnHost string
}

type TrainingDataManager struct {
	params *TrainingDataParameters

	clientset *datastoreclientsetv1alpha1.DatastoreV1alpha1Client
}

func NewTrainingDataManager(params *TrainingDataParameters) Manager {
	if err := utils.Mkdir(params.LocalPathDir); err != nil {
		log.WithField("localDir", params.LocalPathDir).WithError(err).Fatal("Failed to make local directory for training data")
	}

	return &TrainingDataManager{
		params:    params,
		clientset: utils.BuildInClusterDataStoreClientset(),
	}
}

func (mgr *TrainingDataManager) Cook() error {
	files, err := os.ReadDir(mgr.params.LocalPathDir)
	if err != nil {
		log.WithField("dir", mgr.params.LocalPathDir).WithError(err).Error("Failed to access the local directory for the training data")
		return err
	}
	ctx := context.Background()
	dlr, err := mgr.clientset.DataLoadRequests(mgr.params.Namespace).Get(ctx, mgr.params.ConfigName, metav1.GetOptions{})
	if err != nil {
		log.WithFields(log.Fields{"namespace": mgr.params.Namespace, "dataloadrequest": mgr.params.ConfigName}).WithError(err).Error("Failed to get dataloadrequest")
		return err
	}

	if !dlr.Spec.IsGlobal && dlr.Spec.Node != mgr.params.NodeName {
		log.WithFields(log.Fields{"config": mgr.params.ConfigName, "spec": dlr.Spec}).Debug("The config is not for me")
		return fmt.Errorf("no valid config")
	}

	isLoadedBefore := utils.IsStringInSet(mgr.params.NodeName, dlr.Status.ReadyNodes)

	if len(files) > 0 && isLoadedBefore {
		return nil
	}

	if err := mgr.loadData(dlr); err != nil {
		return err
	}
	if !isLoadedBefore {
		for range []int{1, 2, 3} {
			dlr, err := mgr.clientset.DataLoadRequests(mgr.params.Namespace).Get(ctx, mgr.params.ConfigName, metav1.GetOptions{})
			if err != nil {
				log.WithFields(log.Fields{"namespace": mgr.params.Namespace, "dataloadrequest": mgr.params.ConfigName}).WithError(err).Error("Failed to get dataloadrequest")
				continue
			}
			oldDlr := dlr.DeepCopy()
			dlr.Status.ReadyNodes = append(dlr.Status.ReadyNodes, mgr.params.NodeName)
			patch := client.MergeFrom(oldDlr)
			patchData, _ := patch.Data(dlr)
			_, err = mgr.clientset.DataLoadRequests(mgr.params.Namespace).Patch(ctx, dlr.Name, patch.Type(), patchData, metav1.PatchOptions{}, "status")
			if err == nil {
				return nil
			}
			log.WithFields(log.Fields{"namespace": mgr.params.Namespace, "dataloadrequest": mgr.params.ConfigName}).WithError(err).Error("Failed to update status of dataloadrequest")
			time.Sleep(time.Second)
		}
		return fmt.Errorf("failed to update status of dataloadrequest")
	}

	return nil
}

func (mgr *TrainingDataManager) loadData(dlr *datastorev1alpha1.DataLoadRequest) error {
	if err := os.RemoveAll(mgr.params.LocalPathDir); err != nil {
		log.WithField("localdir", mgr.params.LocalPathDir).WithError(err).Error("Failed to clean up the local directory for training data")
		return err
	}
	ds, err := mgr.clientset.DataSets(dlr.Namespace).Get(context.Background(), dlr.Spec.DataSet, metav1.GetOptions{})
	if err != nil {
		log.WithFields(log.Fields{"namespace": dlr.Namespace, "datasource": dlr.Spec.DataSet}).WithError(err).Error("Failed to get datasource")
		return err
	}
	if ds.Spec.Type == "minio" && ds.Spec.MinIO != nil {
		ds.Spec.MinIO.Prefix = filepath.Join(ds.Spec.MinIO.Prefix, dlr.Spec.SubDir)
		log.WithField("minio", ds.Spec.MinIO).Debug("Start to load data ...")
		err := minio.LoadObjectsFromDragonfly(ds.Spec.MinIO, mgr.params.LocalPathDir, ds.Name)
		if err != nil {
			// clean up the directory when loading data fails
			os.RemoveAll(mgr.params.LocalPathDir)
		}
		return err
	}
	log.WithField("datasource", filepath.Join(dlr.Namespace, dlr.Name)).Debug("The data source is not supported yet")
	return fmt.Errorf("unsupported data source")
}

func (mgr *TrainingDataManager) Run(stopCh <-chan struct{}) {
	log.Debug("Watching for training data loading request ...")
}
