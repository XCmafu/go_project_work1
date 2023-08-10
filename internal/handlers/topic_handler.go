// internal/handlers/topic_handler.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"gocode/project_work1/internal/models"
	"gocode/project_work1/internal/storage"
	"net/http"
	"strconv"
)

// IndexHandler 处理显示所有话题的请求
func IndexHandler(c *gin.Context) {
	topicList := storage.LoadTopics()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"topics": topicList,
	})
}

// TopicHandler 处理显示特定话题及回帖的请求
func TopicHandler(c *gin.Context) {
	topicID, _ := strconv.Atoi(c.Param("id"))
	topic := storage.LoadTopic(topicID)
	if topic == nil {
		c.String(http.StatusNotFound, "Topic not found")
		return
	}
	postList := storage.LoadPosts(topicID)
	c.HTML(http.StatusOK, "topic.tmpl", gin.H{
		"topic":  topic,
		"posts":  postList,
		"postId": 0, // Set a placeholder value for new post form
	})
}

// PostHandler 处理发布新回帖的请求
func PostHandler(c *gin.Context) {
	topicID, _ := strconv.Atoi(c.PostForm("topic_id"))
	content := c.PostForm("content")
	post := models.Post{
		TopicID: topicID,
		Content: content,
	}
	storage.SavePost(post)
	c.Redirect(http.StatusSeeOther, "/topic/"+strconv.Itoa(topicID))
}

// NewTopicHandler 处理发布新话题的请求
func NewTopicHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	topic := models.Topic{
		Title:   title,
		Content: content,
	}
	storage.SaveTopic(topic)
	c.Redirect(http.StatusSeeOther, "/")
}
