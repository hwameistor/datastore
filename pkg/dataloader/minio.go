package dataloader

import (
	"context"
	"fmt"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type minioDataLoader struct {
}

func newDataLoaderForMinIO() DataLoader {
	return &minioDataLoader{}
}

func (dl *minioDataLoader) newClient(config *datastorev1alpha1.MinIOSpec) (*minio.Client, error) {
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
}

func (dlr *minioDataLoader) Load(request *datastorev1alpha1.DataLoadRequest, rootDir string) error {

	log.WithFields(log.Fields{"request": request.Name}).Debug("Handling a MinIO data loading request ...")
	// Initialize minio client object.
	minioClient, err := dlr.newClient(request.Spec.MinIO)
	if err != nil {
		log.WithField("endpoint", request.Spec.MinIO.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}

	ctx := context.Background()
	for obj := range minioClient.ListObjects(ctx, request.Spec.MinIO.Bucket, minio.ListObjectsOptions{Prefix: request.Spec.MinIO.Prefix, Recursive: true}) {
		localFilePath := fmt.Sprintf("%s/%s", rootDir, obj.Key)
		if err := minioClient.FGetObject(ctx, request.Spec.MinIO.Bucket, obj.Key, localFilePath, minio.GetObjectOptions{}); err != nil {
			log.WithFields(log.Fields{"request": request.Name, "obj": obj.Key}).WithError(err).Error("Failed to download an object")
			return err
		}
	}

	return nil

}
