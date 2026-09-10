package infrastructure

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	config "go-starter/internal/config/domain"
)

type S3Adapter struct {
	client          *minio.Client
	clientExternal  *minio.Client
	bucket          string
	presignedExpiry time.Duration
	privateEndpoint string
	publicEndpoint  string
}

func NewS3Adapter(cfg config.IConfig) (*S3Adapter, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.S3Domain(), cfg.S3Port())
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey(), cfg.S3SecretKey(), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	rawPublicEndpoint := cfg.S3PublicEndpoint()
	parsedURL, _ := url.Parse(rawPublicEndpoint)
	publicEndpoint := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	clientExternal, err := minio.New(parsedURL.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey(), cfg.S3SecretKey(), ""),
		Secure: parsedURL.Scheme == "https",
		Region: cfg.S3Region(),
	})
	if err != nil {
		return nil, err
	}

	expiry := cfg.S3BucketExpiry()

	return &S3Adapter{
		client:          client,
		clientExternal:  clientExternal,
		bucket:          cfg.S3Bucket(),
		presignedExpiry: expiry,
		privateEndpoint: endpoint,
		publicEndpoint:  publicEndpoint,
	}, nil
}

func (a *S3Adapter) Init(ctx context.Context) error {
	exists, _ := a.client.BucketExists(ctx, a.bucket)
	if !exists {
		if err := a.client.MakeBucket(ctx, a.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	slog.Info("S3 adapter initialized", "bucket", a.bucket, "endpoint", a.privateEndpoint, "public", a.publicEndpoint)
	return nil
}

func (a *S3Adapter) UploadFile(ctx context.Context, key string, body []byte, contentType string) (string, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	_, err := a.client.PutObject(ctx, a.bucket, key, strings.NewReader(string(body)), int64(len(body)), opts)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return key, nil
}

func (a *S3Adapter) GetSignedUrl(ctx context.Context, key string) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := a.clientExternal.PresignedGetObject(ctx, a.bucket, key, a.presignedExpiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return presignedURL.String(), nil
}

func (a *S3Adapter) GetSignedUploadUrl(ctx context.Context, key string, contentType string) (string, error) {
	presignedURL, err := a.clientExternal.PresignedPutObject(ctx, a.bucket, key, a.presignedExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed upload URL: %w", err)
	}
	return presignedURL.String(), nil
}

func (a *S3Adapter) GetPublicUrl(key string) string {
	return fmt.Sprintf("%s/%s/%s", a.publicEndpoint, a.bucket, key)
}

func (a *S3Adapter) DeleteFile(ctx context.Context, key string) error {
	err := a.client.RemoveObject(ctx, a.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (a *S3Adapter) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := a.client.StatObject(ctx, a.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		e := minio.ToErrorResponse(err)
		if e.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *S3Adapter) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	reader, err := a.client.GetObject(ctx, a.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return data, nil
}

func (a *S3Adapter) RemoveBucket(ctx context.Context) error {
	for obj := range a.client.ListObjects(ctx, a.bucket, minio.ListObjectsOptions{Recursive: true}) {
		if obj.Err != nil {
			return fmt.Errorf("list objects: %w", obj.Err)
		}
		a.client.RemoveObject(ctx, a.bucket, obj.Key, minio.RemoveObjectOptions{})
	}

	err := a.client.RemoveBucket(ctx, a.bucket)
	if err != nil {
		return fmt.Errorf("failed to remove bucket: %w", err)
	}
	return nil
}
