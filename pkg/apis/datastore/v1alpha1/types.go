package v1alpha1

import "time"

const (
	StorageBackendTypeMinIO   = "minio"
	StorageBackendTypeAWSS3   = "aws-s3"
	StorageBackendTypeNFS     = "nfs"
	StorageBackendTypeFTP     = "ftp"
	StorageBackendTypeUnknown = "unknown"
)

const (
	StorageBackendConnectionCheckInterval = 1 * time.Minute
)
