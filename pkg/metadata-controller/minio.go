package metadatacontroller

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/minio"
)

func (mgr *dataSourceManager) _checkConnectionForMinIO(backend *datastorev1alpha1.DataSource) (bool, error) {
	if backend.Spec.MinIO == nil {
		return false, fmt.Errorf("invaild MinIO spec info")
	}

	// Initialize minio client object.
	minioClient, err := minio.NewClient(backend.Spec.MinIO)
	if err != nil {
		log.WithField("endpoint", backend.Spec.MinIO.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return false, err
	}

	return minio.IsConnected(minioClient, backend.Spec.MinIO)

}

func (mgr *dataSourceManager) handleDataSourceForMinIO(backend *datastorev1alpha1.DataSource) error {

	log.WithFields(log.Fields{"backend": backend.Name}).Debug("Handling a MinIO storage backend ...")
	// Initialize minio client object.
	minioClient, err := minio.NewClient(backend.Spec.MinIO)
	if err != nil {
		log.WithField("endpoint", backend.Spec.MinIO.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}

	log.WithFields(log.Fields{"bucket": backend.Spec.MinIO.Bucket}).Debug("Checking for the bucket ...")

	ctx := context.Background()
	exists, err := minio.IsBucketExists(minioClient, backend.Spec.MinIO.Bucket)
	if err != nil {
		log.WithField("bucket", backend.Spec.MinIO.Bucket).WithError(err).Error("Failed to check if the bucket exists or not")
		if backend.Status.Error != err.Error() {
			backend.Status.Error = err.Error()
			if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
				log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update status of storage backend")
			}
		}
		return err
	}
	if !exists {
		msg := "Not found bucket"
		log.WithField("bucket", backend.Spec.MinIO.Bucket).WithError(err).Error(msg)
		if backend.Status.Error != msg {
			backend.Status.Error = msg
			if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
				log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
			}

		}
		return fmt.Errorf(msg)
	}
	if backend.Status.Error != "" {
		backend.Status.Error = ""
		if _, err := mgr.dsClientSet.DatastoreV1alpha1().DataSources(mgr.namespace).UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
			log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
			return err
		}
	}
	log.WithFields(log.Fields{"bucket": backend.Spec.MinIO.Bucket}).Debug("The bucket exists")

	objs := minio.LoadObjectMetadata(minioClient, backend.Spec.MinIO)

	mgr.globalView.UpdateDataObjects(backend, objs)

	return nil
}
