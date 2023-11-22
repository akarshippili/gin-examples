package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/akarshippili/gin-examples/fs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	log.Default().Println(ctx.Request)
	// get user, it was set by the BasicAuth middleware
	user := ctx.MustGet(gin.AuthUserKey).(string)
	ctx.IndentedJSON(200, gin.H{
		"message": "Hey! " + user,
	})
}

func Index(ctx *gin.Context) {
	log.Default().Println(ctx.Request)
	ctx.HTML(200, "form.html", nil)
}

func Profile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
	}

	ctx.SaveUploadedFile(file, "./data/"+file.Filename)
	ctx.JSON(201, gin.H{
		"message": "accepted",
	})
}

func GetBuckets(client fs.S3GetbucktesAPI) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		listBucketsOutput, err := fs.GetBuckets(context.TODO(), client, nil)
		log.Default().Printf("buckets %v \n", listBucketsOutput)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.HTML(http.StatusOK, "buckets.html", gin.H{"buckets": listBucketsOutput.Buckets})
	}
}

func GetBucketObjects(client fs.S3GetBucketObjectsAPI) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		bucketid := ctx.Param("bucketid")
		log.Default().Printf("listing objects in bucket \"%s\"", bucketid)

		listObjectsV2Output, err := fs.GetBucketObjects(context.TODO(), client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketid),
		})
		if err != nil {
			ctx.Error(err)
		}

		ctx.HTML(http.StatusOK, "objects.html", gin.H{"objects": listObjectsV2Output.Contents})
	}
}

func GetObject(client fs.S3GetObjectAPI) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		bucketid := ctx.Param("bucketid")
		objectid := ctx.Param("objectid")[1:]
		log.Default().Printf("accessing object \"%s\" in bucket \"%s\"", objectid, bucketid)

		objectOutput, err := fs.GetObject(
			context.TODO(),
			client,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketid),
				Key:    aws.String(objectid),
			},
		)
		if err != nil {
			ctx.Error(err)
			return
		}

		// case 1:
		// bytes, err := fs.GetBytes(objectOutput.Body)
		// if err != nil {
		// 	ctx.Error(err)
		// 	return
		// }

		// ctx.Data(http.StatusOK, "application/octet-stream", bytes)

		// case 2:
		// ctx.DataFromReader(http.StatusOK, objectOutput.ContentLength, "application/octet-stream", objectOutput.Body, nil)

		// case 3:
		// ctx.DataFromReader(http.StatusOK, objectOutput.ContentLength, *objectOutput.ContentType, objectOutput.Body, nil)

		// // case 4:
		splits := strings.Split(objectid, "/")
		filename := splits[len(splits)-1]
		additionalHeaders := make(map[string]string)
		additionalHeaders["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s"`, filename)
		ctx.DataFromReader(http.StatusOK, objectOutput.ContentLength, *objectOutput.ContentType, objectOutput.Body, additionalHeaders)
	}
}
