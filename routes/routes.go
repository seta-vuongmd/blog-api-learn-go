package routes

import (
	"blog-api-learn-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// User routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Post routes
	r.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPost)
	r.PUT("/posts/:id", controllers.UpdatePost)
	r.DELETE("/posts/:id", controllers.DeletePost)
	r.GET("/posts/search-by-tag", controllers.SearchPostsByTag)
	r.GET("/posts/search", controllers.SearchPostsES)
	r.GET("/posts/:id/related", controllers.GetRelatedPosts)
}
