package persistentvolume

import (
	"context"
	dsclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	hmclientset "github.com/hwameistor/hwameistor/pkg/apis/client/clientset/versioned"
	hwameistor "github.com/hwameistor/hwameistor/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/hwameistor/pkg/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	v12 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

var localVolumeTemplate = &hwameistor.LocalVolume{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"hwameistor.io/acceleration-dataset": "true",
		},
	},
	Spec: hwameistor.LocalVolumeSpec{},
}

type PVController interface {
	Run(stopCh <-chan struct{})
}

// pvController is a controller to manage PersistentVolume
type pvController struct {
	dsClientset *dsclientset.Clientset
	hmClientset *hmclientset.Clientset
	kubeClient  *kubernetes.Clientset

	pvLister       corelisters.PersistentVolumeLister
	pvListerSynced cache.InformerSynced
	pvQueue        *common.TaskQueue
}

func New(kubeClientset *kubernetes.Clientset, dsClientset *dsclientset.Clientset, hmClientset *hmclientset.Clientset, pvInformer v12.PersistentVolumeInformer) PVController {
	ctr := &pvController{
		dsClientset: dsClientset,
		kubeClient:  kubeClientset,
		hmClientset: hmClientset,
		pvQueue:     common.NewTaskQueue("PersistentVolumeTask", 0),
	}

	pvInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctr.pvAdded,
		UpdateFunc: ctr.pvUpdated,
		DeleteFunc: ctr.pvDeleted,
	})
	ctr.pvLister = pvInformer.Lister()
	ctr.pvListerSynced = pvInformer.Informer().HasSynced

	return ctr
}

func (ctr *pvController) pvAdded(obj interface{}) {
	pv := obj.(*v1.PersistentVolume)
	if !isPVForDataset(pv) {
		return
	}
	ctr.pvQueue.Add(pv.Name)
}

func isPVForDataset(pv *v1.PersistentVolume) bool {
	if pv.Annotations == nil {
		return false
	}
	return pv.Annotations["hwameistor.io/acceleration-dataset"] == "true"
}

func (ctr *pvController) pvUpdated(oldObj, newObj interface{}) {
	ctr.pvAdded(newObj)
}

func (ctr *pvController) pvDeleted(obj interface{}) {
	ctr.pvAdded(obj)
}

func (ctr *pvController) Run(stopCh <-chan struct{}) {
	defer ctr.pvQueue.Shutdown()

	klog.Infof("Starting PersistentVolume controller")
	defer klog.Infof("Shutting PersistentVolume controller")

	if !cache.WaitForCacheSync(stopCh, ctr.pvListerSynced) {
		klog.Fatalf("Cannot sync caches")
	}

	go wait.Until(ctr.syncPersistentVolume, 0, stopCh)
	<-stopCh
}

func (ctr *pvController) syncPersistentVolume() {
	pvName, quiet := ctr.pvQueue.Get()
	if quiet {
		return
	}
	defer ctr.pvQueue.Done(pvName)

	klog.V(4).Infof("Started PersistentVolume porcessing %q", pvName)

	// get PersistentVolume to process
	ds, err := ctr.pvLister.Get(pvName)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.V(4).Infof("PersistentVolume %q has been deleted, ignoring", pvName)
			return
		}
		klog.Errorf("Error getting PersistentVolume %q: %v", pvName, err)
		ctr.pvQueue.AddRateLimited(pvName)
		return
	}
	ctr.SyncNewOrUpdatedPersistentVolume(ds)
}

func (ctr *pvController) SyncNewOrUpdatedPersistentVolume(pv *v1.PersistentVolume) {
	klog.V(4).Infof("Processing PersistentVolume %s", pv.Name)

	var err error
	if pv.DeletionTimestamp != nil {
		if err = ctr.deleteRelatedLocalVolume(pv.Name); err == nil {
			klog.V(4).Infof("Async Delete LocalVolume %s", pv.Name)
		}
	} else {
		// check if LV created for this PersistentVolume
		if _, err = ctr.hmClientset.HwameistorV1alpha1().LocalVolumes().Get(context.Background(), pv.Name, metav1.GetOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				klog.Errorf("Error getting PV for PersistentVolume %s: %v", pv.Name, err)
				ctr.pvQueue.AddRateLimited(pv.Name)
				return
			}
			// LV isn't found, create it
			if pv.Spec.CSI == nil || pv.Spec.CSI.VolumeAttributes == nil || len(pv.Spec.CSI.VolumeAttributes["poolClass"]) == 0 {
				klog.Errorf("Error getting poolClass for PersistentVolume from CSI.VolumeAttributes %s: %v", pv.Name, err)
				ctr.pvQueue.AddRateLimited(pv.Name)
				return
			}
			poolClass, _ := pv.Spec.CSI.VolumeAttributes["poolClass"]
			if err = ctr.createRelatedLocalVolume(pv.Name, poolClass, pv.Spec.Capacity.Storage().Value()); err == nil {
				klog.V(4).Infof("Created LocalVolume %s", pv.Name)
			}
		} else {
			// LV exists, sync lv or pv status
			err = ctr.syncPVStatus(pv)
		}
	}

	if err != nil {
		klog.V(4).Infof("Error processing PersistentVolume %s: %v", pv.Name, err)
		ctr.pvQueue.AddRateLimited(pv.Name)
		return
	}

	ctr.pvQueue.Forget(pv.Name)
	klog.V(4).Infof("Finished processing PersistentVolume %s", pv.Name)
}

func (ctr *pvController) deleteRelatedLocalVolume(lvName string) (err error) {
	if _, err = ctr.hmClientset.HwameistorV1alpha1().LocalVolumes().Get(context.Background(), lvName, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		// LV isn't found, may be deleted already
		return nil
	}
	_, err = ctr.hmClientset.HwameistorV1alpha1().LocalVolumes().Patch(context.Background(), lvName, types.MergePatchType, []byte(`{"spec":{"delete":true}}`), metav1.PatchOptions{})
	return
}

func (ctr *pvController) createRelatedLocalVolume(pvName, poolClass string, capacityBytes int64) (err error) {
	newLV := localVolumeTemplate.DeepCopy()

	newLV.Name = pvName
	newLV.Spec.RequiredCapacityBytes = capacityBytes
	newLV.Spec.PoolName = "LocalStorage_Pool" + poolClass

	_, err = ctr.hmClientset.HwameistorV1alpha1().LocalVolumes().Create(context.Background(), newLV, metav1.CreateOptions{})
	return
}

func (ctr *pvController) syncPVStatus(pv *v1.PersistentVolume) (err error) {
	switch pv.Status.Phase {
	case v1.VolumeReleased:
		// clean reclaim reference to make pv Available
		_, err = ctr.kubeClient.CoreV1().PersistentVolumes().Patch(context.Background(), pv.Name, types.MergePatchType, []byte(`{"spec":{"claimRef":null}}`), metav1.PatchOptions{})
		if err == nil {
			klog.V(4).Infof("Cleaned reclaim reference for PersistentVolume %s", pv.Name)
		}
	case v1.VolumeAvailable, v1.VolumePending, v1.VolumeBound:
		klog.V(4).Infof("PersistentVolume %s is %s, no need to sync", pv.Name, pv.Status.Phase)
		return nil
	case v1.VolumeFailed:
		klog.V(4).Infof("PersistentVolume %s is %s, can not be recycled", pv.Name, pv.Status.Phase)
		return nil
	default:
		klog.V(4).Infof("unsupported PersistentVolume %s status %s", pv.Name, pv.Status.Phase)
		return nil
	}

	return err
}
