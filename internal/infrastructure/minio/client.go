package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "luca-s3"

type Client struct {
	client   *minio.Client
	endpoint string
	useSSL   bool
}

func New(endpoint, accessKey, secretKey string, useSSL bool) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: new client: %w", err)
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("minio: check bucket: %w", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: create bucket: %w", err)
		}
	}

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
	if err := minioClient.SetBucketPolicy(ctx, bucketName, policy); err != nil {
		return nil, fmt.Errorf("minio: set public policy: %w", err)
	}

	return &Client{client: minioClient, endpoint: endpoint, useSSL: useSSL}, nil
}

func (c *Client) PublicURL(objectName string) string {
	scheme := "http"
	if c.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, c.endpoint, bucketName, objectName)
}

func (c *Client) UploadObject(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) UploadFile(ctx context.Context, objectName, filePath, contentType string) error {
	_, err := c.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) DownloadFile(ctx context.Context, objectName, filePath string) error {
	return c.client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
}

func (c *Client) ObjectExists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Client) DeleteObject(ctx context.Context, objectName string) error {
	return c.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (c *Client) DeleteFolder(ctx context.Context, prefix string) error {
	objectsCh := c.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	for obj := range objectsCh {
		if obj.Err != nil {
			return obj.Err
		}
		if err := c.DeleteObject(ctx, obj.Key); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) StatObject(ctx context.Context, objectName string) (size int64, exists bool) {
	info, err := c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return 0, false
	}
	return info.Size, true
}

func (c *Client) ListObjects(ctx context.Context, prefix string) []string {
	objectsCh := c.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	var objects []string
	for obj := range objectsCh {
		if obj.Err != nil {
			continue
		}
		objects = append(objects, obj.Key)
	}
	return objects
}
