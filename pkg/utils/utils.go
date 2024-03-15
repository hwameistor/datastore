package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	datastoreclientsetv1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/typed/datastore/v1alpha1"
)

func BuildInClusterDataStoreClientset() *datastoreclientsetv1alpha1.DatastoreV1alpha1Client {
	config, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return datastoreclientsetv1alpha1.NewForConfigOrDie(config)
}

// BuildInClusterClientset builds a kubernetes in-cluster clientset
func BuildInClusterClientset() *kubernetes.Clientset {
	config, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return kubernetes.NewForConfigOrDie(config)
}

func Checksum(fpath string) (string, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Mkdir(dirpath string) error {
	f, err := os.Stat(dirpath)
	if err != nil {
		return os.MkdirAll(dirpath, os.ModePerm)
	}
	if !f.IsDir() {
		return os.MkdirAll(dirpath, os.ModePerm)
	}
	return nil
}

func IsStringInSet(str string, strArray []string) bool {
	for _, s := range strArray {
		if str == s {
			return true
		}
	}
	return false
}
