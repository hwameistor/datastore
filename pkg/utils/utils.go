package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	"math"
	"os"
	"strconv"
	"unicode"

	datastoreclientsetv1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/typed/datastore/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func BuildInClusterDataStoreClientset() *datastoreclientsetv1alpha1.DatastoreV1alpha1Client {
	config, err := createKubeConfig("", "")
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return datastoreclientsetv1alpha1.NewForConfigOrDie(config)
}

// BuildInClusterClientset builds a kubernetes in-cluster clientset
func BuildInClusterClientset() *kubernetes.Clientset {
	config, err := createKubeConfig("", "")
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

var (
	unitArray = []string{"B", "KB", "MB", "GB", "TB"}
	unitMap   = map[string]int64{
		"b":  1,
		"B":  1,
		"k":  1000,
		"K":  1024,
		"KB": 1024,
		"m":  1000 * 1000,
		"M":  1024 * 1024,
		"MB": 1024 * 1024,
		"g":  1000 * 1000 * 1000,
		"G":  1024 * 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"t":  1000 * 1000 * 1000 * 1000,
		"T":  1024 * 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
)

// ParseBytes parse size from string into bytes
func ParseBytes(sizeStr string) (int64, error) {
	numStr := ""
	unitStr := ""
	for i := range sizeStr {
		if !unicode.IsDigit(rune(sizeStr[i])) {
			numStr = sizeStr[:i]
			unitStr = sizeStr[i:]
			break
		} else {
			numStr = sizeStr[:i+1]
		}
	}

	if numStr == "" {
		return -1, fmt.Errorf("wrong number: %s", sizeStr)
	}
	if unitStr == "" {
		unitStr = "B"
	}
	multiple, has := unitMap[unitStr]
	if !has {
		return -1, fmt.Errorf("wrong unit: %s", sizeStr)
	}

	size, err := strconv.ParseInt(numStr, 10, 32)
	if err != nil {
		return -1, err
	}
	return size * multiple, nil
}

// ConvertBytesToStr convert size into string
func ConvertBytesToStr(size int64) string {
	unitIndex := 0
	for size > 1024 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%d%s", size, unitArray[unitIndex])
}

// ConvertBytesToIEC converts bytes to IEC units (Ki, Mi, Gi, etc.) and rounds up
func ConvertBytesToIEC(bytes uint64) string {
	if bytes == 0 {
		return "0"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d", bytes)
	}
	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	pre := "KMGTPE"[exp-1]
	value := math.Ceil(float64(bytes) / math.Pow(unit, float64(exp)))
	return fmt.Sprintf("%.0f%ci", value, pre)
}

func CapacityRoundUp(originSize int64) int64 {
	const unit = 1024
	exp := int(math.Log(float64(originSize)) / math.Log(unit)) // "KMGTPE"
	if exp > 3 {
		return originSize
	}

	base := int64(1)
	for i := 0; i < exp; i++ {
		base *= 1024
	}

	if originSize%base == 0 {
		return originSize
	}

	return ((originSize + base) / base) * base
}

const (
	defaultQPS   = 50.0
	defaultBurst = 100
)

func createKubeConfig(masterUrl, kubeconfig string) (*restclient.Config, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)
	if err != nil {
		return nil, err
	}

	kubeConfig.QPS = defaultQPS
	kubeConfig.Burst = defaultBurst
	return kubeConfig, nil
}
