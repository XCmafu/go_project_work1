// cmd/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"gocode/project_work1/internal/handlers"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", handlers.IndexHandler)
	r.GET("/topic/:id", handlers.TopicHandler)
	r.POST("/post", handlers.PostHandler)
	r.POST("/topic", handlers.NewTopicHandler)

	r.Run(":8080")
}
