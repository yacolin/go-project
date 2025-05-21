package configs

import (
	"fmt"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	ossClient *oss.Client
	ossBucket *oss.Bucket
	ossOnce   sync.Once
)

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// InitOSS initializes OSS client with configuration
func InitOSS(config OSSConfig) error {
	var initErr error
	ossOnce.Do(func() {
		client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
		if err != nil {
			initErr = fmt.Errorf("failed to create OSS client: %w", err)
			return
		}

		bucket, err := client.Bucket(config.BucketName)
		if err != nil {
			initErr = fmt.Errorf("failed to get bucket: %w", err)
			return
		}

		ossClient = client
		ossBucket = bucket
	})
	return initErr
}

// GetOSSBucket returns the initialized OSS bucket
func GetOSSBucket() *oss.Bucket {
	return ossBucket
}
