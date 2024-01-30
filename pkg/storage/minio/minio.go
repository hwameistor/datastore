package minio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	datastoreapis "github.com/hwameistor/datastore/pkg/apis"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

func NewClient(config *datastorev1alpha1.MinIOSpec) (*minio.Client, error) {
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
}

func IsBucketExists(client *minio.Client, bucketName string) (bool, error) {
	return client.BucketExists(context.Background(), bucketName)
}

func IsConnected(client *minio.Client, spec *datastorev1alpha1.MinIOSpec) (bool, error) {

	hcancel, _ := client.HealthCheck(1 * time.Second)
	defer hcancel()

	time.Sleep(3 * time.Second)

	return client.IsOnline(), nil
}

func LoadObjectMetadata(client *minio.Client, spec *datastorev1alpha1.MinIOSpec) []*datastoreapis.DataObject {
	objs := []*datastoreapis.DataObject{}
	for obj := range client.ListObjects(context.Background(), spec.Bucket, minio.ListObjectsOptions{Prefix: spec.Prefix, Recursive: true}) {
		items := strings.Split(obj.Key, "/")
		objs = append(objs, &datastoreapis.DataObject{Path: obj.Key, Name: items[len(items)-1], Size: obj.Size, MTime: obj.LastModified})
	}

	return objs
}

func LoadObjectToLocal(client *minio.Client, spec *datastorev1alpha1.MinIOSpec, rootDir string) error {

	ctx := context.Background()
	for obj := range client.ListObjects(ctx, spec.Bucket, minio.ListObjectsOptions{Prefix: spec.Prefix, Recursive: true}) {
		localFilePath := fmt.Sprintf("%s/%s", rootDir, obj.Key)
		if err := client.FGetObject(ctx, spec.Bucket, obj.Key, localFilePath, minio.GetObjectOptions{}); err != nil {
			return fmt.Errorf("failed to load object %s with error : %s", obj.Key, err.Error())
		}
	}
	return nil
}
