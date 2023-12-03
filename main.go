package main

import (
	"context"
	"log"
	"net/url"

	"github.com/akarshippili/gin-examples/middlewares"
	"github.com/akarshippili/gin-examples/router"
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

	r.Use(middlewares.CORSMiddleware())

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

	restRouter := r.Group("/rest", gin.BasicAuth(gin.Accounts{
		"akarsh": "ippili",
		"mike":   "ross",
		"harvey": "specter",
	}))

	authRouter.GET("/ping", router.Ping)
	authRouter.GET("/index", router.Index)
	authRouter.POST("/profile", router.Profile)
	authRouter.GET("/buckets", router.GetBuckets(client))
	authRouter.GET("/buckets/:bucketid", router.GetBucketObjects(client))
	authRouter.GET("/buckets/:bucketid/objects/*objectid", router.GetObject(client))
	authRouter.POST("/buckets/:bucketid/objects", router.PostObject(client))

	restRouter.GET("/buckets", router.GetBucketsJson(client))

	err = r.Run("0.0.0.0:2620")
	if err != nil {
		log.Fatal(err.Error())
	}
}
