package dataloader

import (
	"context"
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	datastoreclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	datastoreinformers "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

const (
	localDataRootDirDefault   = "~/_build/llm/data"
	localCheckpointDirDefault = "~/_build/llm/checkpoints"
)

var (
	nodeName      = flag.String("nodename", "", "Node name")
	targetDir     = flag.String("targetDir", localDataRootDirDefault, fmt.Sprintf("Target directory to load data, default = %s", localDataRootDirDefault))
	checkpointDir = flag.String("checkpointDir", localCheckpointDirDefault, fmt.Sprintf("Directory to save checkpoints, default = %s", localCheckpointDirDefault))
)

type DataLoader interface {
	Load(request *datastorev1alpha1.DataLoadRequest, rootDir string) error
}

type Manager struct {
	dsClientSet *datastoreclientset.Clientset

	localDataRootDir string

	localCheckpointDir string

	nodeName string
}

func NewManager() *Manager {
	flag.Parse()

	return &Manager{
		localDataRootDir:   *targetDir,
		localCheckpointDir: *checkpointDir,
		nodeName:           *nodeName,
	}
}

func (mgr *Manager) Run(stopCh <-chan struct{}) {
	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to get kubernetes cluster config")
	}

	log.Debug("start datastore informer factory")
	mgr.dsClientSet = datastoreclientset.NewForConfigOrDie(cfg)
	dsFactory := datastoreinformers.NewSharedInformerFactory(mgr.dsClientSet, 0)
	dsFactory.Start(stopCh)

	dataloadreqs := dsFactory.Datastore().V1alpha1().DataLoadRequests()
	dataloadreqs.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mgr.onDataLoaderAdd,
		UpdateFunc: mgr.onDataLoaderUpdate,
		DeleteFunc: mgr.onDataLoaderDelete,
	})
	go dataloadreqs.Informer().Run(stopCh)

	<-stopCh
}

func (mgr *Manager) onDataLoaderAdd(obj interface{}) {
	request, _ := obj.(*datastorev1alpha1.DataLoadRequest)
	if request.Spec.Node == mgr.nodeName || len(request.Spec.Node) == 0 {
		mgr.handleDataLoadRequest(request)
	}
}

func (mgr *Manager) onDataLoaderUpdate(oldObj, newObj interface{}) {
	mgr.onDataLoaderAdd(newObj)
}

func (mgr *Manager) onDataLoaderDelete(obj interface{}) {

}

func (mgr *Manager) handleDataLoadRequest(request *datastorev1alpha1.DataLoadRequest) {
	if request.Status.Phase == datastorev1alpha1.DataLoadingPhaseCompleted {
		return
	}
	if request.Status.Error != "" && !request.Spec.Retry {
		return
	}

	var err error
	request.Status.LoadingStartTime = metav1.NewTime(time.Now())
	if request.Spec.Type == datastorev1alpha1.StorageBackendTypeMinIO {
		err = newDataLoaderForMinIO().Load(request, mgr.localDataRootDir)
	} else if request.Spec.Type == datastorev1alpha1.StorageBackendTypeNFS {
		err = newDataLoaderForNFS().Load(request, mgr.localDataRootDir)
	}

	if err != nil {
		request.Status.Error = err.Error()
		request.Status.Phase = datastorev1alpha1.DataLoadingPhaseFailed
	} else {
		request.Status.LoadingCompleteTime = metav1.NewTime(time.Now())
		request.Status.Phase = datastorev1alpha1.DataLoadingPhaseCompleted
		request.Status.Error = ""
	}

	if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataLoadRequests().UpdateStatus(context.Background(), request, metav1.UpdateOptions{}); err != nil {
		log.WithField("request", request.Name).WithError(err).Error("Failed to update dataloading request's status")
	}
}
