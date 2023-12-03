package router

import (
	"context"
	"log"
	"net/http"

	"github.com/akarshippili/gin-examples/fs"
	"github.com/gin-gonic/gin"
)

func GetBucketsJson(client fs.S3GetbucktesAPI) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		listBucketsOutput, err := fs.GetBuckets(context.TODO(), client, nil)
		log.Default().Printf("buckets %v \n", listBucketsOutput)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, listBucketsOutput.Buckets)
	}
}
