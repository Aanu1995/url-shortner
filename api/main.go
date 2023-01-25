package main

import (
	"log"
	"os"

	"github.com/Aanu1995/url-shortner/api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	// load the environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// get the port number from environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	router := gin.Default()

	router.GET("/:url", controllers.ResolveURL)
	router.POST("/api/v1", controllers.ShortenURL)

	router.Run(":" + port)
}
