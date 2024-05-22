package minio

import (
	"context"
	"fmt"
	"github.com/hwameistor/datastore/pkg/exechelper/nsexecutor"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	datastoreapis "github.com/hwameistor/datastore/pkg/apis"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/exechelper"
	"github.com/hwameistor/datastore/pkg/exechelper/basicexecutor"
)

const (
	podUIidEnvVar   = "POD_UUID"
	podNameEnvVar   = "POD_NAME"
	NameSpaceEnvVar = "NAMESPACE"
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

func LoadObjectsFromDragonfly(spec *datastorev1alpha1.MinIOSpec, localDir string, dataSourceName string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Error("Failed to create config object ")
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err).Error("Failed to create  Clientset")
		return err
	}

	podUID := os.Getenv(podUIidEnvVar)
	namespace := os.Getenv(NameSpaceEnvVar)
	podName := os.Getenv(podNameEnvVar)

	podClient := clientset.CoreV1().Pods(namespace)
	pod, err := podClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to get pod:%s", podName)
		return err
	}
	volumeName := dataSourceName
	genDir, err := findGenDirPath(pod, volumeName, clientset)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	fmt.Printf("Generated directory path: %s\n", genDir)

	if localDir[0:len(genDir)] == genDir {
		localDir = localDir[len(genDir):]
	} else {
		return fmt.Errorf("The data storage path configuration is incorrect:%s", localDir)
	}

	output := "/var/lib/kubelet/pods/" + podUID + "/volumes/kubernetes.io~csi/" + volumeName + "/mount" + localDir
	log.Debugf(fmt.Sprintf("output is %s:", output))

	spec.Prefix = strings.TrimRight(strings.TrimLeft(spec.Prefix, "/"), "/")

	var dfget []string
	dfget = append(dfget,
		"--recursive",
		"--level=100",
		"--header", "awsEndpoint: "+spec.Endpoint,
		"--header", "awsRegion: "+spec.Region,
		"--header", "awsAccessKeyID: "+spec.AccessKey,
		"--header", "awsSecretAccessKey: "+spec.SecretKey,
		"--header", "awsS3ForcePathStyle: true",
		"--url", fmt.Sprintf("s3://%s/%s/", spec.Bucket, spec.Prefix),
		"--output", output,
	)

	params := exechelper.ExecParams{
		CmdName: "dfget",
		CmdArgs: dfget,
		Timeout: int(time.Hour),
	}
	return nsexecutor.New().RunCommand(params).Error
}

// findGenDirPath 在给定的Pod中查找与特定PersistentVolumeClaim (PVC) 相关联的卷的挂载路径
func findGenDirPath(pod *v1.Pod, volumeName string, clientset kubernetes.Interface) (string, error) {
	pvClient := clientset.CoreV1().PersistentVolumes()
	pv, err := pvClient.Get(context.TODO(), volumeName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get persistent volume %s: %v", volumeName, err)
	}

	if pv == nil {
		return "", fmt.Errorf("persistent volume %s not found", volumeName)
	}

	if pv.Spec.ClaimRef != nil {
		pvcName := pv.Spec.ClaimRef.Name
		genDir, err := findMountPathForPVC(pod, pvcName)
		if err != nil {
			return "", fmt.Errorf("failed to find mount path for PVC %s: %v", pvcName, err)
		}
		return genDir, nil
	} else {
		return "", fmt.Errorf("persistent volume ClaimRef %s not found", volumeName)
	}

}

// findMountPathForPVC 在Pod中查找与指定PVC关联的卷的挂载路径
func findMountPathForPVC(pod *v1.Pod, pvcName string) (string, error) {
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			if volume.PersistentVolumeClaim.ClaimName == pvcName {
				genDir, found := findMountPathInContainers(pod.Spec.Containers, volume.Name)
				if found {
					return genDir, nil
				}
				return "", fmt.Errorf("failed to find mount path for PVC %s", pvcName)
			}
		}
	}

	return "", fmt.Errorf("PVC %s not found in pod volumes", pvcName)
}

func findMountPathInContainers(containers []v1.Container, volumeName string) (string, bool) {
	for _, container := range containers {
		for _, mount := range container.VolumeMounts {
			if mount.Name == volumeName {
				return mount.MountPath, true
			}
		}
	}
	return "", false
}
