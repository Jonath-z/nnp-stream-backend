package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nnp-stream-backend/handlers"
	"github.com/nnp-stream-backend/handlers/webhook"
)

func validateEnvVars() error {
	requiredEnvVars := []string{"MUX_TOKEN_ID", "MUX_TOKEN_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}
	return nil
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	if err := validateEnvVars(); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}
}

func main() {
	r := gin.Default()
	// r.Use(middleware.Cors)
	r.GET("/", handlers.HealthCheck)
	r.GET("/mux-signed-url", handlers.HandleMuxSignedUploadUrl)
	r.POST("/mux-web-hook", webhook.HandleMuxWebhook)
	r.POST("/clerk-web-hook", webhook.HandleCleckWebhook)

	r.NoRoute(func(c *gin.Context) {
		log.Printf("Unhandled route: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{"message": "Not found"})
	})
	r.Run()
}
