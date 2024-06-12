package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	_ "github.com/minio/minio-go/v7"
	"os"
	"strconv"
	"testing"
)

const (
	serverEndpoint = "SERVER_ENDPOINT"
	accessKey      = "ACCESS_KEY"
	secretKey      = "SECRET_KEY"
	enableHTTPS    = "ENABLE_HTTPS"
	enableKMS      = "ENABLE_KMS"
)

func TestNewClient(t *testing.T) {
	mc, err := NewClientFor(os.Getenv(serverEndpoint), os.Getenv(accessKey), os.Getenv(secretKey), mustParseBool(os.Getenv(enableHTTPS)))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Client: %v is offline: %v", mc, mc.IsOffline())
}

func TestGetBucketCapacity(t *testing.T) {
	testCases := []struct {
		desc       string
		bucketName string
		objectName string
		filePath   string
		objectSize int64
	}{
		{
			desc:       "get bucket capacity",
			bucketName: "datastore-test-bucket",
			objectName: "test.txt",
			filePath:   "test.txt",
			objectSize: 1024,
		},
	}

	for _, testCase := range testCases {
		mc, err := NewClientFor(os.Getenv(serverEndpoint), os.Getenv(accessKey), os.Getenv(secretKey), mustParseBool(os.Getenv(enableHTTPS)))
		if err != nil {
			t.Fatal(err)
		}

		t.Run("parent", func(t *testing.T) {
			if err = makeBucket(mc.Client, testCase.bucketName); err != nil {
				t.Fatal(err)
			}
			if err = createFileWithSize(testCase.objectName, testCase.objectSize); err != nil {
				t.Fatal(err)
			}
			if err = putObject(mc.Client, testCase.objectName, testCase.bucketName, testCase.filePath, "text/plain"); err != nil {
				t.Fatal(err)
			}

			t.Run(testCase.desc, func(t *testing.T) {
				capacity, err := mc.GetBucketCapacity(testCase.bucketName)
				if err != nil {
					t.Fatal(err)
				}

				if capacity != testCase.objectSize {
					t.Fatalf("expected %d, got %d", testCase.objectSize, capacity)
				}
			})
		})
	}
}

func makeBucket(mc *minio.Client, bucketName string) error {
	err := mc.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, err := mc.BucketExists(context.TODO(), bucketName)
		if err == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func putObject(mc *minio.Client, objectName, bucketName, filePath, contentType string) error {
	_, err := mc.FPutObject(context.TODO(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func createFileWithSize(filename string, size int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, size)
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// If server endpoint is not set, all tests default to
	// using https://play.min.io
	if os.Getenv(serverEndpoint) == "" {
		os.Setenv(serverEndpoint, "play.min.io")
		os.Setenv(accessKey, "Q3AM3UQ867SPQQA43P2F")
		os.Setenv(secretKey, "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
		os.Setenv(enableHTTPS, "1")
	}
}

// Convert string to bool and always return false if any error
func mustParseBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return b
}
