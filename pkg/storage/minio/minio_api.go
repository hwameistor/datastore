package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MC struct {
	*minio.Client
}

func NewClientFor(endpoint, id, secret string, secure bool) (*MC, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(id, secret, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return &MC{cli}, nil
}

func (mc *MC) GetBucketCapacity(bucketName string) (int64, error) {
	bucketSize := int64(0)
	for object := range mc.ListObjects(context.TODO(), bucketName, minio.ListObjectsOptions{Recursive: true}) {
		if object.Err != nil {
			return bucketSize, object.Err
		}
		bucketSize += object.Size
	}
	return bucketSize, nil
}
