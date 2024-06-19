package dataset

import (
	"context"
	"fmt"
	dsclientset "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned"
	dsinformers "github.com/hwameistor/datastore/pkg/apis/client/informers/externalversions/datastore/v1alpha1"
	dslisters "github.com/hwameistor/datastore/pkg/apis/client/listers/datastore/v1alpha1"
	datastore "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"strings"

	smino "github.com/hwameistor/datastore/pkg/storage/minio"
	hmclientset "github.com/hwameistor/hwameistor/pkg/apis/client/clientset/versioned"
	"github.com/hwameistor/hwameistor/pkg/common"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	persistentVolumeTemplate = &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"hwameistor.io/acceleration-dataset": "true", // to identify the dataset volume
			},
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany},
			Capacity: v1.ResourceList{
				v1.ResourceStorage: resource.MustParse("1Gi"),
			},
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource: v1.PersistentVolumeSource{
				CSI: &v1.CSIPersistentVolumeSource{
					Driver: "lvm.hwameistor.io",
					FSType: "xfs",
				},
			},
		},
	}
	minStorageCapacity = int64(4194304)
	minStorageQuantity = resource.NewQuantity(minStorageCapacity, resource.BinarySI)
)

type DSController interface {
	Run(stopCh <-chan struct{})
}

// dsController is a controller to manage DataSet
type dsController struct {
	dsClientset *dsclientset.Clientset
	hmClientset *hmclientset.Clientset
	kubeClient  *kubernetes.Clientset

	dsLister       dslisters.DataSetLister
	dsListerSynced cache.InformerSynced
	dsQueue        *common.TaskQueue
}

func New(kubeClientset *kubernetes.Clientset, dsClientset *dsclientset.Clientset, hmClientset *hmclientset.Clientset, dsInformer dsinformers.DataSetInformer) DSController {
	ctr := &dsController{
		dsClientset: dsClientset,
		kubeClient:  kubeClientset,
		hmClientset: hmClientset,
		dsQueue:     common.NewTaskQueue("DataSetTask", 0),
	}

	dsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctr.dsAdded,
		UpdateFunc: ctr.dsUpdated,
		DeleteFunc: ctr.dsDeleted,
	})
	ctr.dsLister = dsInformer.Lister()
	ctr.dsListerSynced = dsInformer.Informer().HasSynced

	return ctr
}

func (ctr *dsController) dsAdded(obj interface{}) {
	ds := obj.(*datastore.DataSet)
	ctr.dsQueue.Add(ds.Namespace + "/" + ds.Name)
}

func (ctr *dsController) dsUpdated(oldObj, newObj interface{}) {
	ctr.dsAdded(newObj)
}

func (ctr *dsController) dsDeleted(obj interface{}) {
	ctr.dsAdded(obj)
}

func (ctr *dsController) Run(stopCh <-chan struct{}) {
	defer ctr.dsQueue.Shutdown()

	klog.Infof("Starting DataSet controller")
	defer klog.Infof("Shutting DataSet controller")

	if !cache.WaitForCacheSync(stopCh, ctr.dsListerSynced) {
		klog.Fatalf("Cannot sync caches")
	}

	go wait.Until(ctr.syncDataSource, 0, stopCh)
	<-stopCh
}

func (ctr *dsController) syncDataSource() {
	key, quiet := ctr.dsQueue.Get()
	if quiet {
		return
	}
	defer ctr.dsQueue.Done(key)

	klog.V(4).Infof("Started DataSet porcessing %q", key)
	dsNamespace := strings.Split(key, "/")[0]
	dsName := strings.Split(key, "/")[1]

	// get DataSet to process
	ds, err := ctr.dsLister.DataSets(dsNamespace).Get(dsName)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.V(4).Infof("DataSet %q has been deleted, ignoring", key)
			return
		}
		klog.Errorf("Error getting DataSet %q: %v", key, err)
		ctr.dsQueue.AddRateLimited(key)
		return
	}
	ctr.SyncNewOrUpdatedDatasource(ds)
}

func (ctr *dsController) SyncNewOrUpdatedDatasource(ds *datastore.DataSet) {
	klog.V(4).Infof("Processing DataSet %s/%s", ds.Namespace, ds.Name)

	var err error
	defer func() {
		if err != nil {
			klog.V(4).Infof("Error processing DataSet %s/%s: %v", ds.Namespace, ds.Name, err)
			ctr.dsQueue.AddRateLimited(ds.Namespace + "/" + ds.Name)
			return
		}

		ctr.dsQueue.Forget(ds.Namespace + "/" + ds.Name)
		klog.V(4).Infof("Finished processing DataSet %s/%s", ds.Namespace, ds.Name)
	}()

	// DS is deleting, release relevant pv
	if ds.DeletionTimestamp != nil {
		if err = ctr.deleteRelatedPersistentVolume(ds.Name); err == nil {
			klog.V(4).Infof("Async Delete PersistentVolume %s", ds.Name)
		}
		return
	}

	if ds.Spec.CapacityBytes == 0 {
		if err = ctr.updateDatasetCapacity(ds); err != nil {
			klog.Errorf("Failed to update capacity for DataSet %s/%s: %v", ds.Namespace, ds.Name, err)
			return
		}

		klog.V(4).Infof("Updated capacity for DataSet %s/%s", ds.Namespace, ds.Name)
		return
	}

	// check if PV created for this DataSet
	_, err = ctr.kubeClient.CoreV1().PersistentVolumes().Get(context.Background(), ds.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			klog.Errorf("Error getting PV for DataSet %s/%s: %v", ds.Namespace, ds.Name, err)
			return
		}
		// PV not found, create it
		poolClass := "HDD"
		if err = ctr.createRelatedPersistentVolume(ds.Name, poolClass, ds.Spec.CapacityBytes); err == nil {
			klog.V(4).Infof("Created PersistentVolume %s", ds.Name)
		}
	}
}

func (ctr *dsController) deleteRelatedPersistentVolume(pvName string) error {
	deletePolicy := metav1.DeletePropagationBackground
	return ctr.kubeClient.CoreV1().PersistentVolumes().Delete(context.Background(), pvName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}

func (ctr *dsController) createRelatedPersistentVolume(pvName, poolClass string, capacityBytes int64) (err error) {
	newPV := persistentVolumeTemplate.DeepCopy()

	volumeAttr := make(map[string]string)
	volumeAttr["convertible"] = "false"
	volumeAttr["csi.storage.k8s.io/pv/name"] = pvName
	volumeAttr["volumeKind"] = "LVM"
	volumeAttr["poolClass"] = poolClass

	volumeMode := v1.PersistentVolumeFilesystem
	newPV.Name = pvName
	newPV.Spec.VolumeMode = &volumeMode
	newPV.Spec.CSI.VolumeAttributes = volumeAttr
	newPV.Spec.PersistentVolumeSource.CSI.VolumeHandle = pvName

	storageQuantity := resource.NewQuantity(utils.CapacityRoundUp(capacityBytes), resource.BinarySI)
	if storageQuantity.CmpInt64(minStorageCapacity) < 0 {
		storageQuantity = minStorageQuantity
	}
	newPV.Spec.Capacity[v1.ResourceStorage] = *storageQuantity

	_, err = ctr.kubeClient.CoreV1().PersistentVolumes().Create(context.Background(), newPV, metav1.CreateOptions{})
	return
}

func (ctr *dsController) getDatasetCapacity(ds *datastore.DataSet) (int64, error) {
	switch ds.Spec.Type {
	case datastore.DataSourceTypeMinIO:
		mspec := ds.Spec.MinIO
		if mspec == nil {
			return 0, fmt.Errorf("MinIO spec is nil")
		}

		mc, err := smino.NewClientFor(mspec.Endpoint, mspec.AccessKey, mspec.SecretKey, false)
		if err != nil {
			return 0, err
		}
		return mc.GetBucketCapacity(mspec.Bucket)
	case datastore.DataSourceTypeAWSS3:
		fallthrough
	case datastore.DataSourceTypeNFS:
		fallthrough
	case datastore.DataSourceTypeFTP:
		fallthrough
	case datastore.DataSourceTypeUnknown:
		klog.V(4).Infof("Unsupported capacity update for DataSourceType %s", ds.Spec.Type)
	}

	return 0, fmt.Errorf("unsupported capacity update for DataSourceType %s", ds.Spec.Type)
}

func (ctr *dsController) updateDatasetCapacity(ds *datastore.DataSet) error {
	capacity, err := ctr.getDatasetCapacity(ds)
	if err != nil {
		return err
	}
	patchCapacity := fmt.Sprintf("{\"spec\":{\"capacityBytes\":%d}}", capacity)
	_, err = ctr.dsClientset.DatastoreV1alpha1().DataSets(ds.Namespace).Patch(context.Background(), ds.Name, types.MergePatchType, []byte(patchCapacity), metav1.PatchOptions{})
	return err
}
