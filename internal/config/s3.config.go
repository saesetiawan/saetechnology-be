package config

import "github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/storage"

func LoadS3Config(cfg *Config) storage.S3Config {
	return storage.S3Config{
		Region:          cfg.S3AwsRegion,
		AccessKeyID:     cfg.S3AccessKeyID,
		SecretAccessKey: cfg.S3SecretAccessKey,
		Bucket:          cfg.S3BucketName,
		Endpoint:        cfg.S3Endpoint,
		PublicBaseURL:   cfg.S3PublicBaseURL,
	}
}
