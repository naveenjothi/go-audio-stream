package storage

import "errors"

var (
	ErrMissingKeyID          = errors.New("B2_KEY_ID is required")
	ErrMissingApplicationKey = errors.New("B2_APPLICATION_KEY is required")
	ErrMissingBucketName     = errors.New("B2_BUCKET_NAME is required")
	ErrMissingRegion         = errors.New("B2_REGION is required")
	ErrMissingEndpoint       = errors.New("B2_ENDPOINT is required")
	ErrUploadFailed          = errors.New("failed to upload file")
	ErrDeleteFailed          = errors.New("failed to delete file")
	ErrPresignFailed         = errors.New("failed to generate presigned URL")
)
