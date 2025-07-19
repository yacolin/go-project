package configs

import (
	"fmt"
	"log"
	"os"
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

// 初始化OSS配置
func ConnectOSS() {
	ossConfig := OSSConfig{
		Endpoint:        os.Getenv("OSS_ENDPOINT"),
		AccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		BucketName:      os.Getenv("OSS_BUCKET_NAME"),
	}

	if err := InitOSS(ossConfig); err != nil {
		log.Fatalf("Failed to initialize OSS: %v", err)
	}
}
