package s3storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/logging"
	"github.com/nessai1/gophkeeper/internal/service/config"
	"github.com/nessai1/gophkeeper/internal/service/mediastorage"
	"github.com/nessai1/gophkeeper/pkg/bytesize"
	"go.uber.org/zap"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const keeperBucketID = "keeper"

type S3Storage struct {
	client *s3.Client

	logger *s3Logger

	config config.S3Config
}

func NewStorage(storageConfig config.S3Config, logger *zap.Logger) (*S3Storage, error) {
	storage := S3Storage{config: storageConfig}
	customResolver := aws.EndpointResolverWithOptionsFunc(func(serviceID, region string, options ...interface{}) (aws.Endpoint, error) {
		if serviceID == s3.ServiceID && region == storage.config.SigningRegion {
			return aws.Endpoint{
				PartitionID:   storage.config.PartitionID,
				URL:           storage.config.URL,
				SigningRegion: storage.config.SigningRegion,
			}, nil
		}

		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg := aws.NewConfig()

	cfg.EndpointResolverWithOptions = customResolver
	cfg.Region = storage.config.SigningRegion

	if storage.config.Credentials == nil {
		return nil, errors.New("cannot create S3 storage without credentials")
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(
		storage.config.Credentials.AccessKeyID,
		storage.config.Credentials.SecretAccessKey,
		"",
	)

	l := s3Logger{logger: logger}
	cfg.Logger = l
	storage.logger = &l

	storage.client = s3.NewFromConfig(*cfg)
	if storage.getBucket(keeperBucketID) == nil {
		return nil, fmt.Errorf("cannot find bucket for keeper service (%s) in S3 storage", keeperBucketID)
	}

	return &storage, nil
}

func (s *S3Storage) getBucket(name string) *types.Bucket {
	result, err := s.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		if *bucket.Name == name {
			return &bucket
		}
	}

	return nil
}

func (s *S3Storage) StartDownload(ctx context.Context, key string) (io.ReadCloser, error) {
	bucketID := keeperBucketID

	object, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucketID, Key: &key})
	if err != nil {
		return nil, fmt.Errorf("cannot get object from S3 storage: %w", err)
	}

	return object.Body, nil
}

func (s *S3Storage) StartUpload(ctx context.Context, key string) (mediastorage.MultipartUpload, error) {
	u, err := newS3MultipartUpload(ctx, key, keeperBucketID, s.client)
	if err != nil {
		return nil, fmt.Errorf("cannot create new media upload: %w", err)
	}

	return u, nil
}

type s3Logger struct {
	logger *zap.Logger
}

func (s s3Logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	logMsg := fmt.Sprintf(format, v)
	logMsg = fmt.Sprintf("[S3 LOG] %s", logMsg)

	s.logger.Info(logMsg, zap.String("classification", string(classification)))
}

type S3MultipartUpload struct {
	key      string
	bucketID string
	uploadID string

	activeContent []byte

	partsCount int32
	client     *s3.Client

	parts []types.CompletedPart
}

func newS3MultipartUpload(ctx context.Context, key string, bucketID string, client *s3.Client) (*S3MultipartUpload, error) {
	u := S3MultipartUpload{
		key:        key,
		bucketID:   bucketID,
		partsCount: 1,
		client:     client,
	}

	o, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{Bucket: &bucketID, Key: &key})
	if err != nil {
		return nil, fmt.Errorf("cannot create new S3 upload: %w", err)
	}

	u.uploadID = *o.UploadId

	return &u, nil
}

func (u *S3MultipartUpload) Upload(ctx context.Context, content []byte) error {
	if u.activeContent != nil {
		u.activeContent = append(u.activeContent, content...)
	} else {
		u.activeContent = content
	}

	// Yandex cloud docs: content part must be 5MB or more
	if len(u.activeContent) > int(bytesize.MB)*5 {
		return u.sendActiveContent(ctx)
	}

	return nil
}

func (u *S3MultipartUpload) sendActiveContent(ctx context.Context) error {
	l := int64(len(u.activeContent))
	b := bytes.NewBuffer(u.activeContent)

	resp, err := u.client.UploadPart(ctx, &s3.UploadPartInput{Bucket: &u.bucketID, Key: &u.key, UploadId: &u.uploadID, PartNumber: &u.partsCount, ContentLength: &l, Body: b})
	if err != nil {
		return fmt.Errorf("cannot send S3 data part: %w", err)
	}

	var cnt = u.partsCount
	u.parts = append(u.parts, types.CompletedPart{
		ETag:       resp.ETag,
		PartNumber: &cnt,
	})

	u.partsCount++
	u.activeContent = nil

	return nil
}

func (u *S3MultipartUpload) Complete(ctx context.Context) error {
	if u.activeContent != nil {
		err := u.sendActiveContent(ctx)
		if err != nil {
			return fmt.Errorf("cannot send active content while complete S3 upload: %w", err)
		}
	}

	_, err := u.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{Bucket: &u.bucketID, Key: &u.key, UploadId: &u.uploadID, MultipartUpload: &types.CompletedMultipartUpload{Parts: u.parts}})
	if err != nil {
		return fmt.Errorf("cannot complete s3 data upload: %w", err)
	}

	return nil
}

func (u *S3MultipartUpload) Abort(ctx context.Context) error {
	_, err := u.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{Bucket: &u.bucketID, Key: &u.key, UploadId: &u.uploadID})
	if err != nil {
		return fmt.Errorf("cannot abort s3 data upload: %w", err)
	}

	return nil
}
