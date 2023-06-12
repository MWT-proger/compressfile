package s3storage

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/MWT-proger/compressfile/internal/transform"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
}

type OperationStorager interface {
	UploadFileToServer(bucketName string, objectKey string) ([]byte, error)
	GetClientS3() (*s3.Client, error)
}

// Load the SDK's configuration from environment and shared config, and
// create the client with this.
func (s Storage) GetClientS3() (*s3.Client, error) {

	config := configs.GetConfig()
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: config.EndpointURLS3Storage,
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		// TODO: вставить тут свою ошибку
		return nil, err
	}
	// TODO: Исправить этот момент
	return s3.NewFromConfig(cfg), nil
}

// func (s Storage) GetList() {
// 	// Set the parameters based on the CLI flag inputs.
// 	// params := &s3.ListObjectsV2Input{
// 	// 	Bucket: &bucketName,
// 	// }
// 	// if len(objectPrefix) != 0 {
// 	// 	params.Prefix = &objectPrefix
// 	// }
// 	// if len(objectDelimiter) != 0 {
// 	// 	params.Delimiter = &objectDelimiter
// 	// }

// 	// // Create the Paginator for the ListObjectsV2 operation.
// 	// p := s3.NewListObjectsV2Paginator(s.S3Client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
// 	// 	if v := int32(maxKeys); v != 0 {
// 	// 		o.Limit = v
// 	// 	}
// 	// })

// 	// // Iterate through the S3 object pages, printing each object returned.
// 	// var i int
// 	// var a int
// 	// log.Println("Objects:")
// 	// for p.HasMorePages() {
// 	// 	i++

// 	// 	// Next Page takes a new context for each page retrieval. This is where
// 	// 	// you could add timeouts or deadlines.
// 	// 	page, err := p.NextPage(context.TODO())
// 	// 	if err != nil {
// 	// 		log.Fatalf("failed to get page %v, %v", i, err)
// 	// 	}

// 	// 	// Log the objects found
// 	// 	for _, obj := range page.Contents {
// 	// 		a++
// 	// 		fmt.Println("Object:", *obj.Key)
// 	// 	}
// 	// }
// 	// fmt.Println("Objects count:", a)

// }

// func (s Storage) DownloadFile(bucketName string, objectKey string, fileName string) error {
// 	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
// 		return err
// 	}
// 	defer result.Body.Close()
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
// 		return err
// 	}
// 	defer file.Close()
// 	body, err := io.ReadAll(result.Body)
// 	if err != nil {
// 		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
// 	}
// 	_, err = file.Write(body)
// 	return err
// }

// func (s Storage) GetBodyObject(bucketName string, objectKey string) (io.Reader, error) {
// 	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
// 		return nil, err
// 	}

// 	return result.Body, nil

// }

// func (s Storage) UploadFile(bucketName string, objectKey string, fileName string) error {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
// 	} else {
// 		defer file.Close()
// 		_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
// 			Bucket: aws.String(bucketName),
// 			Key:    aws.String(objectKey),
// 			Body:   file,
// 		})
// 		if err != nil {
// 			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
// 				fileName, bucketName, objectKey, err)
// 		}
// 	}
// 	return err
// }

func (s Storage) UploadFileToServer(bucketName string, objectKey string) ([]byte, error) {

	log.Printf("++++++++++++++++++++++++++++++++++++++++++++++++=")
	client, errClient := s.GetClientS3()

	if errClient != nil {
		log.Print("Заглушка")
	}
	result, errGet := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if errGet != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, errGet)
		return nil, errGet
	}
	log.Printf("++++++++++++++++++++++++++++++++++++++++++++++++=")

	defer result.Body.Close()

	b, errReady := io.ReadAll(result.Body)
	if errReady != nil {
		return nil, errReady
	}
	opt := transform.ParseOptions("100x150")

	img, _ := transform.Transform(b, opt)

	buf := new(bytes.Buffer)
	buf.Write(img)

	log.Printf("result.ContentLength : %v\n", result.ContentLength)
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String("GoCompress/" + objectKey),
		Body:          bufio.NewReader(buf),
		ContentLength: int64(len(img)),
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			"fileName", bucketName, objectKey, err)
	}
	return img, nil
}
