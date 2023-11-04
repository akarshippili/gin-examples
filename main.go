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

	err := r.Run("localhost:2620")
	if err != nil {
		log.Fatal(err.Error())
	}
}
