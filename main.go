package main

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
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

	err := r.Run("localhost:2620")
	if err != nil {
		log.Fatal(err.Error())
	}
}
