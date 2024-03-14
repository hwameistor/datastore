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
	"github.com/hwameistor/datastore/pkg/exechelper"
	"github.com/hwameistor/datastore/pkg/exechelper/basicexecutor"
)

func newClient(spec *datastorev1alpha1.MinIOSpec) (*minio.Client, error) {
	return minio.New(spec.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(spec.AccessKey, spec.SecretKey, ""),
		Secure: false,
	})
}

func IsConnected(spec *datastorev1alpha1.MinIOSpec) (bool, error) {
	client, err := newClient(spec)
	if err != nil {
		return false, err
	}
	hcancel, _ := client.HealthCheck(1 * time.Second)
	defer hcancel()

	time.Sleep(3 * time.Second)

	return client.IsOnline(), nil
}

func LoadObjectMetadata(spec *datastorev1alpha1.MinIOSpec) ([]*datastoreapis.DataObject, error) {
	client, err := newClient(spec)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	existed, err := client.BucketExists(ctx, spec.Bucket)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, fmt.Errorf("bucket not found")
	}

	objs := []*datastoreapis.DataObject{}
	for obj := range client.ListObjects(ctx, spec.Bucket, minio.ListObjectsOptions{Prefix: spec.Prefix, Recursive: true}) {
		items := strings.Split(obj.Key, "/")
		objs = append(objs, &datastoreapis.DataObject{Path: obj.Key, Name: items[len(items)-1], Size: obj.Size, MTime: obj.LastModified})
	}

	return objs, nil
}

func DownloadObject(spec *datastorev1alpha1.MinIOSpec, objName string, localFilePath string) error {
	client, err := newClient(spec)
	if err != nil {
		return err
	}

	ctx := context.Background()
	existed, err := client.BucketExists(ctx, spec.Bucket)
	if err != nil {
		return err
	}
	if !existed {
		return fmt.Errorf("bucket not found")
	}

	return client.FGetObject(ctx, spec.Bucket, fmt.Sprintf("%s/%s", spec.Prefix, objName), localFilePath, minio.GetObjectOptions{Checksum: true})
}

func UploadObject(spec *datastorev1alpha1.MinIOSpec, objName string, localFilePath string) error {

	client, err := newClient(spec)
	if err != nil {
		return err
	}

	ctx := context.Background()
	existed, err := client.BucketExists(ctx, spec.Bucket)
	if err != nil {
		return err
	}
	if !existed {
		return fmt.Errorf("bucket not found")
	}

	_, err = client.FPutObject(ctx, spec.Bucket, fmt.Sprintf("%s/%s", spec.Prefix, objName), localFilePath, minio.PutObjectOptions{})
	return err
}

func LoadObjects(spec *datastorev1alpha1.MinIOSpec, localDir string) error {
	spec.Prefix = strings.TrimRight(strings.TrimLeft(spec.Prefix, "/"), "/")
	url := fmt.Sprintf("minio://%s:%s@%s/%s/%s/", spec.AccessKey, spec.SecretKey, spec.Endpoint, spec.Bucket, spec.Prefix)
	params := exechelper.ExecParams{
		CmdName: "juicesync",
		CmdArgs: []string{url, strings.TrimRight(localDir, "/") + "/"},
		Timeout: int(time.Hour),
	}
	return basicexecutor.New().RunCommand(params).Error
}
