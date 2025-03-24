package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nnp-stream-backend/handlers"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func ginCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173/", "https://80a9-2c0f-fe30-4a8b-0-69da-5711-ad27-3202.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func main() {
	r := gin.Default()

	r.GET("/", handlers.HealthCheck)
	r.GET("/mux-signed-url", handlers.HandleMuxSignedUploadUrl)
	r.POST("/mux-web-hook", handlers.HandleMuxWebhook)

	r.Run()
}
