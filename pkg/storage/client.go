package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Client wraps the S3 client for Backblaze B2 operations
type Client struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	region        string
}

// FileInfo represents metadata about a stored file
type FileInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ContentType  string
}

// NewClient creates a new B2 storage client
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create custom endpoint resolver for Backblaze B2
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s", cfg.Endpoint),
			}, nil
		},
	)

	// Load AWS config with custom credentials and endpoint
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.KeyID, cfg.ApplicationKey, ""),
		),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &Client{
		s3Client:      s3Client,
		presignClient: s3.NewPresignClient(s3Client),
		bucketName:    cfg.BucketName,
		region:        cfg.Region,
	}, nil
}

// Upload uploads a file to B2 and returns the object key
func (c *Client) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	return key, nil
}

// UploadWithACL uploads a file with specified ACL (public-read or private)
func (c *Client) UploadWithACL(ctx context.Context, key string, reader io.Reader, contentType string, acl string) (string, error) {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACL(acl),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	return key, nil
}

// Delete removes a file from B2
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	return nil
}

// GetPresignedURL generates a time-limited download URL
func (c *Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignedReq, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPresignFailed, err)
	}

	return presignedReq.URL, nil
}

// GetPublicURL returns the public URL for a file (only works if file has public-read ACL)
func (c *Client) GetPublicURL(key string) string {
	return fmt.Sprintf("https://%s/%s/%s", c.bucketName+".s3."+c.region+".backblazeb2.com", c.bucketName, key)
}

// ListFiles lists files with the given prefix
func (c *Client) ListFiles(ctx context.Context, prefix string) ([]FileInfo, error) {
	output, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]FileInfo, 0, len(output.Contents))
	for _, obj := range output.Contents {
		files = append(files, FileInfo{
			Key:          aws.ToString(obj.Key),
			Size:         *obj.Size,
			LastModified: aws.ToTime(obj.LastModified),
		})
	}

	return files, nil
}

// GetBucketName returns the configured bucket name
func (c *Client) GetBucketName() string {
	return c.bucketName
}
