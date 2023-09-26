package storage

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageS3Option struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Token           string
	UseSSL          bool
}

func NewS3Storage(serviceName string, option StorageS3Option) (client *minio.Client, err error) {
	if client, err = minio.New(option.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(option.AccessKeyID, option.SecretAccessKey, option.Token),
		Secure: option.UseSSL,
	}); err != nil {
		err = fmt.Errorf("MinioStorage connect %s", err.Error())
		return
	}
	return
}
