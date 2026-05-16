package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "videos"

type Client struct {
	client *minio.Client
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

	return &Client{client: minioClient}, nil
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

func (c *Client) PresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// ListObjects lista objetos com um prefixo específico.
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
