package s3storage

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/MWT-proger/compressfile/internal/errors"
)

type Storage struct {
	S3Client *s3.Client
}

type OperationStorager interface {
	InitClientS3() error
	Get(bucketName string, objectKey string) ([]byte, error)
	Put(img []byte, bucketName string, objectKey string) error
}

// Load the SDK's configuration from environment and shared config, and
// create the client with this.
func (s *Storage) InitClientS3() error {

	config := configs.GetConfig()
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: config.EndpointURLS3Storage,
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		return err
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return nil
}

func (s *Storage) Get(bucketName string, objectKey string) ([]byte, error) {

	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		if index := strings.Index(err.Error(), "NoSuchKey"); index != -1 {
			return nil, &errors.ErrorNoSuchKeyInS3Storage{}
		}
		return nil, err

	}

	defer result.Body.Close()

	b, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (s *Storage) Put(img []byte, bucketName string, objectKey string) error {

	buf := new(bytes.Buffer)
	buf.Write(img)

	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          bufio.NewReader(buf),
		ContentLength: int64(len(img)),
	})

	if err != nil {
		return err
	}

	return nil
}
