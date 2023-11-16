package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

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

func createObjectFrom(key string, reader io.ReadCloser) {
	defer reader.Close()

	err := os.MkdirAll(filepath.Dir(key), 0770)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write the data from the reader to the file.
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// TODO:
func GetListOfBuckets() []string {
	arr := make([]string, 0)

	return arr
}

func main() {

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
