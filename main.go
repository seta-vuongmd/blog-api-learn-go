package main

import (
	"blog-api-learn-go/config"
	"blog-api-learn-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.ConnectRedis()
	config.ConnectElasticsearch()
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run()
}
