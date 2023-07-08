package main

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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

	err := r.Run("localhost:2620")
	if err != nil {
		log.Fatal(err.Error())
	}
}
