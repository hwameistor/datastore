package metadatacontroller

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newMinIOClient(config *datastorev1alpha1.MinIOSpec) (*minio.Client, error) {
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
}

func (mgr *storageBackendManager) _checkConnectionForMinIO(backend *datastorev1alpha1.StorageBackend) (bool, error) {
	log.WithFields(log.Fields{"backend": backend.Name}).Debug("Checking connection for a MinIO storage backend")
	if backend.Spec.MinIO == nil {
		return false, fmt.Errorf("invaild MinIO spec info")
	}

	// Initialize minio client object.
	minioClient, err := newMinIOClient(backend.Spec.MinIO)
	if err != nil {
		log.WithField("endpoint", backend.Spec.MinIO.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return false, err
	}

	hcancel, _ := minioClient.HealthCheck(1 * time.Second)
	defer hcancel()

	time.Sleep(3 * time.Second)

	return minioClient.IsOnline(), nil

}

func (mgr *storageBackendManager) handleStorageBackendForMinIO(backend *datastorev1alpha1.StorageBackend) error {

	log.WithFields(log.Fields{"backend": backend.Name}).Debug("Handling a MinIO storage backend ...")
	// Initialize minio client object.
	minioClient, err := newMinIOClient(backend.Spec.MinIO)
	if err != nil {
		log.WithField("endpoint", backend.Spec.MinIO.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}

	log.WithFields(log.Fields{"bucket": backend.Spec.MinIO.Bucket}).Debug("Checking for the bucket ...")
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, backend.Spec.MinIO.Bucket)
	if err != nil {
		log.WithField("bucket", backend.Spec.MinIO.Bucket).WithError(err).Error("Failed to check if the bucket exists or not")
		if backend.Status.Error != err.Error() {
			backend.Status.Error = err.Error()
			if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
				log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
			}
		}
		return err
	}
	if !exists {
		msg := "Not found bucket"
		log.WithField("bucket", backend.Spec.MinIO.Bucket).WithError(err).Error(msg)
		if backend.Status.Error != msg {
			backend.Status.Error = msg
			if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
				log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
			}

		}
		return fmt.Errorf(msg)
	}
	if backend.Status.Error != "" {
		backend.Status.Error = ""
		if _, err := mgr.dsClientSet.DatastoreV1alpha1().StorageBackends().UpdateStatus(ctx, backend, metav1.UpdateOptions{}); err != nil {
			log.WithField("Message", backend.Status.Error).WithError(err).Error("Failed to update storage backend")
			return err
		}
	}
	log.WithFields(log.Fields{"bucket": backend.Spec.MinIO.Bucket}).Debug("The bucket exists")

	return mgr.refreshDataFromStorageBackendForMinIO(minioClient, backend)
}

func (mgr *storageBackendManager) refreshDataFromStorageBackendForMinIO(minioClient *minio.Client, backend *datastorev1alpha1.StorageBackend) error {
	ctx := context.Background()

	for obj := range minioClient.ListObjects(ctx, backend.Spec.MinIO.Bucket, minio.ListObjectsOptions{Prefix: backend.Spec.MinIO.Prefix, Recursive: false}) {
		log.WithFields(log.Fields{
			"key":  obj.Key,
			"size": obj.Size,
		}).Debug("Got an object")
	}

	return nil
}
