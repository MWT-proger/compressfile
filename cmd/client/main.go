package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MWT-proger/compressfile/internal/transform"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName      string
	objectPrefix    string
	objectDelimiter string
	maxKeys         int
)

type BucketBasics struct {
	S3Client *s3.Client
}

func init() {
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the S3 bucket to list objects from.")
	flag.StringVar(&objectPrefix, "prefix", "", "The optional `object prefix` of the S3 Object keys to list.")
	flag.StringVar(&objectDelimiter, "delimiter", "",
		"The optional `object key delimiter` used by S3 List objects to group object keys.")
	flag.IntVar(&maxKeys, "max-keys", 0,
		"The maximum number of `keys per page` to retrieve at once.")
}

func (basics BucketBasics) GetList() {
	// Set the parameters based on the CLI flag inputs.
	params := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}
	if len(objectPrefix) != 0 {
		params.Prefix = &objectPrefix
	}
	if len(objectDelimiter) != 0 {
		params.Delimiter = &objectDelimiter
	}

	// Create the Paginator for the ListObjectsV2 operation.
	p := s3.NewListObjectsV2Paginator(basics.S3Client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})

	// Iterate through the S3 object pages, printing each object returned.
	var i int
	var a int
	log.Println("Objects:")
	for p.HasMorePages() {
		i++

		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			a++
			fmt.Println("Object:", *obj.Key)
		}
	}
	fmt.Println("Objects count:", a)

}

func (basics BucketBasics) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

func (basics BucketBasics) GetBodyObject(bucketName string, objectKey string) (io.Reader, error) {
	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}

	return result.Body, nil

}

func (basics BucketBasics) UploadFile(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

func (basics BucketBasics) UploadFileToServer(bucketName string, objectKey string) error {
	// file, err := os.Open(fileName)
	// if err != nil {
	// 	log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)

	result, errGet := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if errGet != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, errGet)
		return errGet
	}

	defer result.Body.Close()

	b, errReady := io.ReadAll(result.Body)
	if errReady != nil {
		return errReady
	}
	opt := transform.ParseOptions("100x150")

	img, _ := transform.Transform(b, opt)

	buf := new(bytes.Buffer)
	buf.Write(img)
	// fmt.Fprintf(buf, "Content-Type: %d\n\n", "image/png")
	// fmt.Fprintf(buf, "Content-Length: %d\n\n", int64(len(img)))

	// log.Printf("buf %v\n\n", buf)

	// fmt.Fprintf(buf, "%s %s\n", result., result.Status)
	// if err := result.Header.WriteSubset(buf, map[string]bool{
	// 	"Content-Length": true,
	// 	// exclude Content-Type header if the format may have changed during transformation
	// 	"Content-Type": opt.Format != "" || result.Header.Get("Content-Type") == "image/webp" || result.Header.Get("Content-Type") == "image/tiff",
	// }); err != nil {
	// 	t.log("error copying headers: %v", err)
	// }
	// fmt.Fprintf(buf, "Content-Length: %d\n\n", len(img))

	log.Printf("result.ContentLength : %v\n", result.ContentLength)
	_, err := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String("GoCompress/" + objectKey),
		Body:          bufio.NewReader(buf),
		ContentLength: int64(len(img)),
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			"fileName", bucketName, objectKey, err)
	}
	return nil
}

// Lists all objects in a bucket using pagination
func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://gateway.storjshare.io",
			// SigningRegion: "us-west-2",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	bucketBasics := BucketBasics{S3Client: s3Client}

	testKey := "collections/c547ffa8-9d26-4141-a54d-f2f4ae4d8153/28d1d5ec-ddc2-4edb-8adf-37e75031e109/tokens/cfad535a28a84b198811d225a61f566a.png"
	// bucketBasics.DownloadFile(bucketName, testKey, "test_image.png")

	// bodyObject, err := bucketBasics.GetBodyObject(bucketName, testKey)
	// if err != nil {
	// 	log.Fatalf("bodyObjec err: , %v", err)
	// }
	bucketBasics.UploadFileToServer(bucketName, testKey)

	log.Print("Operation successfull!")

}
