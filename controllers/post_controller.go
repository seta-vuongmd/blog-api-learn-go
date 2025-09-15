package controllers

import (
	"blog-api-learn-go/config"
	"blog-api-learn-go/models"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	// ...existing code...
	// Index to Elasticsearch
	var post models.Post
	go func(post models.Post) {
		doc := map[string]interface{}{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
		}
		jsonDoc, _ := json.Marshal(doc)
		config.ESClient.Index("posts", bytes.NewReader(jsonDoc))
	}(post)
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		log := models.ActivityLog{
			Action:   "new_post",
			PostID:   post.ID,
			LoggedAt: post.CreatedAt.Unix(),
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Find(&posts)
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	// Try to get from Redis
	val, err := config.RedisClient.Get(config.Ctx, "post:"+id).Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}
	// Cache miss, get from DB
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	// Save to Redis

	jsonData, _ := json.Marshal(post)
	config.RedisClient.Set(config.Ctx, "post:"+id, jsonData, 300000000000) // 5 phút
	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&post)
	// Invalidate Redis cache
	config.RedisClient.Del(config.Ctx, "post:"+id)
	// Đồng bộ lên Elasticsearch
	go func(p models.Post) {
		doc := map[string]interface{}{
			"id":      p.ID,
			"title":   p.Title,
			"content": p.Content,
			"tags":    p.Tags,
		}
		jsonDoc, _ := json.Marshal(doc)
		config.ESClient.Index("posts", bytes.NewReader(jsonDoc))
	}(post)
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Post{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// Search posts by tag using GIN index
// Search posts full-text via Elasticsearch
func SearchPostsES(c *gin.Context) {
	query := c.Query("q")
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "content"},
			},
		},
	}
	jsonQuery, _ := json.Marshal(esQuery)
	res, err := config.ESClient.Search(
		config.ESClient.Search.WithIndex("posts"),
		config.ESClient.Search.WithBody(bytes.NewReader(jsonQuery)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	var esRes map[string]interface{}
	json.NewDecoder(res.Body).Decode(&esRes)
	c.JSON(http.StatusOK, esRes)
}
func SearchPostsByTag(c *gin.Context) {
	tag := c.Query("tag")
	var posts []models.Post
	err := config.DB.Where("tags @> ARRAY[?]::text[]", tag).Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}
func GetRelatedPosts(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	// Tìm các bài viết liên quan qua Elasticsearch
	should := []map[string]interface{}{}
	for _, tag := range post.Tags {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{"tags": tag},
		})
	}
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
				"must_not": []map[string]interface{}{
					{"match": map[string]interface{}{"id": post.ID}},
				},
			},
		},
		"size": 5,
	}
	jsonQuery, _ := json.Marshal(esQuery)
	res, err := config.ESClient.Search(
		config.ESClient.Search.WithIndex("posts"),
		config.ESClient.Search.WithBody(bytes.NewReader(jsonQuery)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	var esRes map[string]interface{}
	json.NewDecoder(res.Body).Decode(&esRes)
	c.JSON(http.StatusOK, esRes)
}
