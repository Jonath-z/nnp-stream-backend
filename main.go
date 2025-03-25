package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nnp-stream-backend/handlers"
	"github.com/nnp-stream-backend/handlers/webhook"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	r := gin.Default()

	r.GET("/", handlers.HealthCheck)
	r.GET("/mux-signed-url", handlers.HandleMuxSignedUploadUrl)
	r.POST("/mux-web-hook", webhook.HandleMuxWebhook)

	// Add a catch-all route to log any unhandled routes
	r.NoRoute(func(c *gin.Context) {
		log.Printf("Unhandled route: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{"message": "Not found"})
	})
	r.Run()
}
