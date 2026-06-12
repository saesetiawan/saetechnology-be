package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/storage"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type S3Storage struct {
	client     *s3.Client
	presign    *s3.PresignClient
	bucket     string
	publicBase string
	tracer     trace.Tracer
}

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	PublicBaseURL   string
}

func NewS3Storage(
	cfg S3Config,
	tracerProvider trace.TracerProvider,
) storage.Storage {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Storage{
		client:     client,
		presign:    s3.NewPresignClient(client),
		bucket:     cfg.Bucket,
		publicBase: cfg.PublicBaseURL,
		tracer:     tracerProvider.Tracer("S3Storage"),
	}
}

func (s *S3Storage) Upload(
	ctx context.Context,
	input storage.UploadInput,
) (*storage.UploadOutput, error) {
	ctx, span := s.tracer.Start(ctx, "S3Storage.Upload")
	defer span.End()

	safeFilename := strings.ToValidUTF8(input.Filename, "_")
	safeDirectory := strings.ToValidUTF8(input.Directory, "_")
	safeBucket := strings.ToValidUTF8(s.bucket, "_")

	extension := filepath.Ext(safeFilename)
	extension = strings.ToLower(extension)

	filename := fmt.Sprintf(
		"%d/%02d/%s%s",
		time.Now().Year(),
		int(time.Now().Month()),
		uuid.NewString(),
		extension,
	)

	key := fmt.Sprintf(
		"%s/%s",
		safeDirectory,
		filename,
	)

	safeKey := strings.ToValidUTF8(key, "_")

	span.SetAttributes(
		attribute.String("storage.provider", "s3"),
		attribute.String("storage.bucket", safeBucket),
		attribute.String("storage.directory", safeDirectory),
		attribute.String("storage.filename", safeFilename),
		attribute.String("storage.extension", strings.ToValidUTF8(extension, "_")),
		attribute.String("storage.key", safeKey),
		attribute.Int("storage.file_size", len(input.Body)),
	)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(input.Body),
		ACL:    types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed upload object to s3")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success upload object to s3")

	return &storage.UploadOutput{
		Key: key,
		URL: s.GetPublicURL(key),
	}, nil
}

func (s *S3Storage) Delete(
	ctx context.Context,
	key string,
) error {
	ctx, span := s.tracer.Start(ctx, "S3Storage.Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("storage.provider", "s3"),
		attribute.String("storage.bucket", s.bucket),
		attribute.String("storage.key", key),
	)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed delete object from s3")
		return err
	}

	span.SetStatus(codes.Ok, "success delete object from s3")

	return nil
}

func (s *S3Storage) GetPresignedUploadURL(
	ctx context.Context,
	input storage.PresignedUploadInput,
) (*storage.PresignedUploadOutput, error) {
	ctx, span := s.tracer.Start(ctx, "S3Storage.GetPresignedUploadURL")
	defer span.End()

	span.SetAttributes(
		attribute.String("storage.provider", "s3"),
		attribute.String("storage.bucket", s.bucket),
		attribute.String("storage.key", input.Key),
		attribute.String("storage.content_type", input.ContentType),
		attribute.Int64("storage.presign_expired_seconds", int64(input.Expired.Seconds())),
	)

	req, err := s.presign.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(input.Key),
		},
		s3.WithPresignExpires(input.Expired),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create presigned url")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success create presigned url")

	return &storage.PresignedUploadOutput{
		URL:    req.URL,
		Method: "PUT",
		Header: map[string]string{
			"Content-Type": input.ContentType,
		},
	}, nil
}

func (s *S3Storage) GetPublicURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicBase, s.bucket, key)
}
