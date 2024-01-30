package dataloader

import (
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/minio"

	log "github.com/sirupsen/logrus"
)

type minioDataLoader struct {
}

func newDataLoaderForMinIO() DataLoader {
	return &minioDataLoader{}
}

func (dlr *minioDataLoader) Load(request *datastorev1alpha1.DataLoadRequest, localRootDir string) error {

	spec := request.Spec.MinIO
	log.WithFields(log.Fields{"request": request.Name}).Debug("Handling a MinIO data loading request ...")
	// Initialize minio client object.
	minioClient, err := minio.NewClient(spec)
	if err != nil {
		log.WithField("endpoint", spec.Endpoint).WithError(err).Error("Failed to setup the minio client")
		return err
	}

	return minio.LoadObjectToLocal(minioClient, spec, localRootDir)

}
