package main

import (
	"blogoproducer/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	corsConfig := cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost",
			"http://localhost:3000/write-post",
			"http://localhost:3000",
			"http://20.127.128.101",
			"http://20.127.128.101:3000",
		},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept", "Authorization", "Access-Control-Request-Credentials", "Access-Control-Request-Origin", "Access-Control-Request-Methods"},
		ExposeHeaders:    []string{"Cookie"},
		AllowCredentials: true,
		MaxAge:           60 * 60 * time.Hour,
	})
	router.Use(corsConfig)
	router.POST("/parse", handlers.ParseHandler)
	router.Run(":5000")
}
