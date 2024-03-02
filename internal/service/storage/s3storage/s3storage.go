package s3storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/logging"
	"github.com/nessai1/gophkeeper/internal/service/config"
	"go.uber.org/zap"
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

func (s *S3Storage) List() {
	// Запрашиваем список бакетов

}

type s3Logger struct {
	logger *zap.Logger
}

func (s s3Logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	logMsg := fmt.Sprintf(format, v)
	logMsg = fmt.Sprintf("[S3 LOG] %s", logMsg)

	s.logger.Info(logMsg, zap.String("classification", string(classification)))
}
