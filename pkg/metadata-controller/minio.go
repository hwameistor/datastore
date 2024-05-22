package metadatacontroller

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/minio"
)

func (mgr *dataSourceManager) _checkConnectionForMinIO(ds *datastorev1alpha1.DataSet) (bool, error) {
	if ds.Spec.MinIO == nil {
		return false, fmt.Errorf("invaild MinIO spec info")
	}

	return minio.IsConnected(ds.Spec.MinIO)

}

func (mgr *dataSourceManager) handleDataSourceForMinIO(ds *datastorev1alpha1.DataSet) error {

	log.WithFields(log.Fields{"ds": ds.Name}).Debug("Handling a MinIO storage ds ...")

	objs, err := minio.LoadObjectMetadata(ds.Spec.MinIO)
	if err != nil {
		if ds.Status.Error != err.Error() {
			ds.Status.Error = err.Error()
			if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSets(mgr.namespace).UpdateStatus(context.Background(), ds, metav1.UpdateOptions{}); err != nil {
				log.WithField("Message", ds.Status.Error).WithError(err).Error("Failed to update status of ds")
			}
		}
		return err
	}

	mgr.globalView.UpdateDataObjects(ds, objs)

	return nil
}
