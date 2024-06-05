package dataload

import (
	"context"
	dlrclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	dlrinformers "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions/datastore/v1alpha1"
	dlrlisters "github.com/hwameistor/datastore/pkg/apis/client/listers/datastore/v1alpha1"
	datastore "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/common"
	"github.com/hwameistor/datastore/pkg/storage/minio"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
	"path/filepath"
	"strings"
)

var mounter *mount.SafeFormatAndMount

func init() {
	mounter = NewMounter()
}

func NewMounter() *mount.SafeFormatAndMount {
	return &mount.SafeFormatAndMount{
		Interface: mount.New("/bin/mount"),
		Exec:      utilexec.New(),
	}
}

type DLRController interface {
	Run(stopCh <-chan struct{})
}

// dlrController is a controller to manage DataLoadRequest
type dlrController struct {
	nodeName        string
	dlrClientset    *dlrclientset.Clientset
	kubeClient      *kubernetes.Clientset
	dlrLister       dlrlisters.DataLoadRequestLister
	dlrListerSynced cache.InformerSynced
	dlrQueue        *common.TaskQueue
}

func New(nodeName string, kubeClientset *kubernetes.Clientset, dlrClientset *dlrclientset.Clientset, dlrInformer dlrinformers.DataLoadRequestInformer) DLRController {
	ctr := &dlrController{
		nodeName:     nodeName,
		dlrClientset: dlrClientset,
		kubeClient:   kubeClientset,
		dlrQueue:     common.NewTaskQueue("DataLoadTask", 0),
	}

	dlrInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctr.dlrAdded,
		UpdateFunc: ctr.dlrUpdated,
		DeleteFunc: ctr.dlrDeleted,
	})
	ctr.dlrLister = dlrInformer.Lister()
	ctr.dlrListerSynced = dlrInformer.Informer().HasSynced

	return ctr
}

func (ctr *dlrController) dlrAdded(obj interface{}) {
	dlr := obj.(*datastore.DataLoadRequest)
	if dlr.Spec.Node == ctr.nodeName {
		ctr.dlrQueue.Add(dlr.Namespace + "/" + dlr.Name)
	}
}

func (ctr *dlrController) dlrUpdated(oldObj, newObj interface{}) {
	ctr.dlrAdded(newObj)
}

func (ctr *dlrController) dlrDeleted(obj interface{}) {
	ctr.dlrAdded(obj)
}

func (ctr *dlrController) Run(stopCh <-chan struct{}) {
	defer ctr.dlrQueue.Shutdown()
	log.Infof("Starting DataLoadRequest controller")
	if !cache.WaitForCacheSync(stopCh, ctr.dlrListerSynced) {
		log.Fatalf("Cannot sync caches")
	}

	go wait.Until(ctr.syncDataLoad, 0, stopCh)
	<-stopCh
}

func (ctr *dlrController) syncDataLoad() {
	key, quiet := ctr.dlrQueue.Get()
	if quiet {
		return
	}
	defer ctr.dlrQueue.Done(key)
	klog.V(0).Infof("Started DataLoadRequest porcessing %q", key)
	dlrNamespace := strings.Split(key, "/")[0]
	dlrName := strings.Split(key, "/")[1]

	// get DataLoadRequest to process
	dlr, err := ctr.dlrLister.DataLoadRequests(dlrNamespace).Get(dlrName)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Errorf("DataLoadRequest %q has been deleted, ignoring", key)
			return
		}
		log.Errorf("Error getting DataLoadRequest %q: %v", key, err)
		ctr.dlrQueue.AddRateLimited(key)
		return
	}
	ctr.SyncNewOrUpdatedDataLoad(dlr)
}

func (ctr *dlrController) SyncNewOrUpdatedDataLoad(dlr *datastore.DataLoadRequest) {
	log.Infof("Processing DataLoadRequest %s/%s", dlr.Namespace, dlr.Name)
	var err error
	log.Debugf("dlr.Status.State: %s", dlr.Status.State)
	switch dlr.Status.State {

	case datastore.OperationStateStart, "":
		err := ctr.DataLoadStart(dlr)
		if err != nil {
			log.Errorf("Error DataLoadSubmit %s/%s: %v", dlr.Namespace, dlr.Name, err)
			ctr.dlrQueue.AddRateLimited(dlr.Namespace + "/" + dlr.Name)
			return
		}

	case datastore.OperationStateCompleted:
		err := ctr.DataLoadComplete(dlr)
		if err != nil {
			log.Errorf("Error DataLoadComplete %s/%s: %v", dlr.Namespace, dlr.Name, err)
			ctr.dlrQueue.AddRateLimited(dlr.Namespace + "/" + dlr.Name)
			return
		}
	}

	if err != nil {
		log.Errorf("Error processing DataLoadRequest %s/%s: %v", dlr.Namespace, dlr.Name, err)
		ctr.dlrQueue.AddRateLimited(dlr.Namespace + "/" + dlr.Name)
		return
	}

	ctr.dlrQueue.Forget(dlr.Namespace + "/" + dlr.Name)
	log.Infof("Finished processing DataLoadRequest %s/%s", dlr.Namespace, dlr.Name)
}

func (ctr *dlrController) DataLoadStart(dlr *datastore.DataLoadRequest) error {

	ds, err := ctr.dlrClientset.DatastoreV1alpha1().DataSets(dlr.Namespace).Get(context.Background(), dlr.Spec.DataSet, metav1.GetOptions{})
	if err != nil {
		log.WithFields(log.Fields{"namespace": dlr.Namespace, "dataset": dlr.Spec.DataSet}).WithError(err).Error("Failed to get dataset")
		return err
	}

	if ds.Spec.Type == "minio" && ds.Spec.MinIO != nil {
		if dlr.Spec.SubDir != "" {
			ds.Spec.MinIO.Prefix = filepath.Join(ds.Spec.MinIO.Prefix, dlr.Spec.SubDir)
		}
		err := minio.LoadObjectsFromDragonflyV2(ctr.kubeClient, ds.Spec.MinIO, dlr.Spec.DstDir, ds.Name)
		if err != nil {
			return err
		}
	}

	dlr.Status.State = datastore.OperationStateCompleted
	_, err = ctr.dlrClientset.DatastoreV1alpha1().DataLoadRequests(dlr.Namespace).UpdateStatus(context.TODO(), dlr, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("dlr.Status.State update OperationStateCompleted")
	return nil
}

func (ctr *dlrController) DataLoadComplete(dlr *datastore.DataLoadRequest) error {
	err := ctr.dlrClientset.DatastoreV1alpha1().DataLoadRequests(dlr.Namespace).Delete(context.TODO(), dlr.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
