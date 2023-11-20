package router

import (
	"context"
	"log"
	"net/http"

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

func Buckets(client fs.S3GetbucktesAPI) func(ctx *gin.Context) {
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

func Objects(client fs.S3GetBucketObjectsAPI) func(ctx *gin.Context) {

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
