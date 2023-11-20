package main

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/akarshippili/gin-examples/fs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{".aws/credentials"}),
		config.WithSharedConfigFiles([]string{".aws/config"}),
		config.WithLogConfigurationWarnings(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// adding a middleware
	r.Use(func(ctx *gin.Context) {
		log.Println("Params")

		request := ctx.Request
		values, _ := url.ParseQuery(request.URL.RawQuery)

		for key, value := range values {
			log.Printf("Key: %v, Value: %v", key, value)
		}
		ctx.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRouter := r.Group("/auth", gin.BasicAuth(gin.Accounts{
		"akarsh": "ippili",
		"mike":   "ross",
		"harvey": "specter",
	}))

	authRouter.GET("/ping", func(ctx *gin.Context) {
		log.Default().Println(ctx.Request)
		// get user, it was set by the BasicAuth middleware
		user := ctx.MustGet(gin.AuthUserKey).(string)
		ctx.IndentedJSON(200, gin.H{
			"message": "Hey! " + user,
		})
	})

	authRouter.GET("/index", func(ctx *gin.Context) {
		log.Default().Println(ctx.Request)
		ctx.HTML(200, "form.html", nil)
	})

	authRouter.POST("/profile", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.Error(err)
		}

		ctx.SaveUploadedFile(file, "./data/"+file.Filename)
		ctx.JSON(201, gin.H{
			"message": "accepted",
		})
	})

	authRouter.GET("/buckets", func(ctx *gin.Context) {
		listBucketsOutput, err := fs.GetBuckets(context.TODO(), client, nil)
		log.Default().Printf("buckets %v \n", listBucketsOutput)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.HTML(http.StatusOK, "buckets.html", gin.H{"buckets": listBucketsOutput.Buckets})
	})

	authRouter.GET("/buckets/:bucketid", func(ctx *gin.Context) {
		bucketid := ctx.Param("bucketid")
		log.Default().Printf("listing objects in bucket %s", bucketid)

		listObjectsV2Output, err := fs.GetBucketObjects(context.TODO(), client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketid),
		})
		if err != nil {
			ctx.Error(err)
		}

		ctx.HTML(http.StatusOK, "objects.html", gin.H{"objects": listObjectsV2Output.Contents})
	})

	err = r.Run("localhost:2620")
	if err != nil {
		log.Fatal(err.Error())
	}
}
