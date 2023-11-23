package fs

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3GetBucketObjectsAPI interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func GetBucketObjects(c context.Context, api S3GetBucketObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

// S3GetObjectAPI defines the interface for the GetObject function.
// We use this interface to test the function using a mocked service.
type S3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// / GetObject gets the access control list (ACL) for an Amazon Simple Storage Service (Amazon S3) bucket object Inputs:
//
//	c is the context of the method call, which includes the AWS Region
//	api is the interface that defines the method call
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a GetObjectOutput object containing the result of the service call and nil
//	Otherwise, nil and an error from the call to GetObject
func GetObject(c context.Context, api S3GetObjectAPI, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return api.GetObject(c, input)
}

type S3GetbucktesAPI interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

func GetBuckets(c context.Context, api S3GetbucktesAPI, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return api.ListBuckets(c, input)
}

type S3PostObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PostObject(c context.Context, api S3PostObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func Test() {

	bucket := flag.String("b", "", "The bucket containing the object")
	objectName := flag.String("o", "", "The bucket object to get ACL from")
	flag.Parse()

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"../.aws/credentials"}),
		config.WithSharedConfigFiles([]string{"../.aws/config"}),
		config.WithLogConfigurationWarnings(true))

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(*bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s || size=%d || last-modified=%v", aws.ToString(object.Key), object.Size, object.LastModified)
	}

	req := &s3.GetObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*objectName),
	}

	res, err := GetObject(context.TODO(), client, req)
	if err != nil {
		fmt.Println("Got an error getting ACL for " + *objectName)
		return
	}

	body := res.Body
	fmt.Println("ContentType:", *res.ContentType)
	createObjectFrom(*objectName, body)
}
