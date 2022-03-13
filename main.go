package main

import (
	"blogoproducer/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/parse", handlers.ParseHandler)
	router.Run(":5000")
}
