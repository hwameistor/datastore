package dataload_manager

import (
	dlrclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	dlrinformers "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/dataload-manager/dataload"
	"k8s.io/client-go/kubernetes"
)

// dataLoadRequestController interface
type DataLoadController interface {
	Run(stopCh <-chan struct{})
}

type dataLoadRequestController struct {
	dlrController dataload.DLRController
}

func New(nodeName string, kubeClientset *kubernetes.Clientset, dlrClientset *dlrclientset.Clientset, dlrInformer dlrinformers.DataLoadRequestInformer) DataLoadController {
	return &dataLoadRequestController{
		dlrController: dataload.New(nodeName, kubeClientset, dlrClientset, dlrInformer),
	}
}

func (dl *dataLoadRequestController) Run(stopCh <-chan struct{}) {
	go dl.dlrController.Run(stopCh)
	<-stopCh
}
