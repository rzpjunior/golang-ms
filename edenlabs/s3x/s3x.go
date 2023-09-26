package s3x

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type Client interface {
	List(ctx context.Context, bucketName string) (buckets []minio.BucketInfo, err error)
	UploadPublicFile(ctx context.Context, bucketName string, fileName string, filePath string, directory string) (info string, err error)
	UploadPrivateFile(ctx context.Context, bucketName string, fileName string, filePath string, directory string) (info string, err error)
}

type S3x struct {
	Client *minio.Client
}

func NewS3x(client *minio.Client) Client {
	return &S3x{
		Client: client,
	}
}

func (b *S3x) List(ctx context.Context, bucketName string) (buckets []minio.BucketInfo, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	buckets, err = b.Client.ListBuckets(ctx)
	if err != nil {
		return
	}

	return
}

func (b *S3x) UploadPublicFile(ctx context.Context, bucketName string, fileName string, filePath string, directory string) (info string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	_, err = b.Client.FPutObject(ctx, bucketName, directory+"/"+fileName, filePath, minio.PutObjectOptions{UserMetadata: userMetaData})
	if err != nil {
		return
	}

	s3Url := &url.URL{
		Scheme: "https",
		Host:   b.Client.EndpointURL().Host,
		Path:   bucketName,
	}

	info = s3Url.String() + "/" + directory + "/" + fileName

	return
}

func (b *S3x) UploadPrivateFile(ctx context.Context, bucketName string, fileName string, filePath string, directory string) (info string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userMetaData := map[string]string{}

	_, err = b.Client.FPutObject(ctx, bucketName, directory+"/"+fileName, filePath, minio.PutObjectOptions{UserMetadata: userMetaData})
	if err != nil {
		return
	}

	reqParams := make(url.Values)
	preSignedURLImage, _ := b.Client.PresignedGetObject(context.Background(), bucketName, directory+"/"+fileName, time.Second*3600, reqParams)

	info = preSignedURLImage.String()

	return
}
