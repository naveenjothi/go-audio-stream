package storage

import "os"

// Config holds the Backblaze B2 configuration
type Config struct {
	KeyID          string
	ApplicationKey string
	BucketName     string
	Region         string
	Endpoint       string
}

// LoadConfig loads B2 configuration from environment variables
func LoadConfig() Config {
	return Config{
		KeyID:          os.Getenv("B2_KEY_ID"),
		ApplicationKey: os.Getenv("B2_APPLICATION_KEY"),
		BucketName:     os.Getenv("B2_BUCKET_NAME"),
		Region:         os.Getenv("B2_REGION"),
		Endpoint:       os.Getenv("B2_ENDPOINT"),
	}
}

// Validate checks if all required configuration values are set
func (c Config) Validate() error {
	if c.KeyID == "" {
		return ErrMissingKeyID
	}
	if c.ApplicationKey == "" {
		return ErrMissingApplicationKey
	}
	if c.BucketName == "" {
		return ErrMissingBucketName
	}
	if c.Region == "" {
		return ErrMissingRegion
	}
	if c.Endpoint == "" {
		return ErrMissingEndpoint
	}
	return nil
}
