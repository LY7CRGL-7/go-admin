package minio

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client MinIO 客户端
type Client struct {
	client *minio.Client
	config *Config
}

// Config MinIO 配置
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// NewClient 创建 MinIO 客户端
func NewClient(cfg *Config) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	c := &Client{
		client: client,
		config: cfg,
	}

	// 确保 bucket 存在
	if err := c.ensureBucketExists(); err != nil {
		return nil, err
	}

	return c, nil
}

// ensureBucketExists 确保 bucket 存在
func (c *Client) ensureBucketExists() error {
	ctx := context.Background()
	exists, err := c.client.BucketExists(ctx, c.config.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// UploadFile 上传文件
func (c *Client) UploadFile(ctx context.Context, objectName string, filePath string, contentType string) (string, error) {
	_, err := c.client.FPutObject(ctx, c.config.BucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// 生成访问 URL
	url, err := c.client.PresignedGetObject(ctx, c.config.BucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// UploadBytes 上传字节数据
func (c *Client) UploadBytes(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := c.client.PutObject(ctx, c.config.BucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload bytes: %w", err)
	}

	// 生成访问 URL
	url, err := c.client.PresignedGetObject(ctx, c.config.BucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// DownloadFile 下载文件
func (c *Client) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
	obj, err := c.client.GetObject(ctx, c.config.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return obj, nil
}

// DeleteFile 删除文件
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.config.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileURL 获取文件访问 URL
func (c *Client) GetFileURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.config.BucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
