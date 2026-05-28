package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nnp-stream-backend/internal/handlers"
	"github.com/nnp-stream-backend/internal/handlers/webhook"
	"github.com/nnp-stream-backend/internal/middleware"

	_ "github.com/nnp-stream-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func validateEnvVars() error {
	requiredEnvVars := []string{
		"MUX_TOKEN_ID",
		"MUX_TOKEN_SECRET",
		"SHWARY_MERCHANT_ID",
		"SHWARY_MERCHANT_KEY",
	}
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

// @title NNP Stream Backend API
// @version 1.0
// @description API for managing video streams, uploads, and analytics powered by Mux
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()
	r.Use(middleware.Cors)
	r.GET("/", handlers.HealthCheck)
	r.GET("/mux-signed-url", handlers.HandleMuxSignedUploadUrl)
	r.POST("/mux-web-hook", webhook.HandleMuxWebhook)
	r.POST("/clerk-web-hook", webhook.HandleCleckWebhook)
	r.POST("/shwary-web-hook", webhook.HandleShwaryWebhook)

	// Subscriptions / payments
	r.POST("/subscriptions", handlers.HandleSubscribeToPlan)
	r.POST("/mux-live-stream", handlers.HandleMuxLiveStream)
	r.DELETE("/mux-live-stream/:liveStreamID", handlers.HandleMuxDeleteLiveStream)

	// Analytics endpoints
	r.GET("/analytics/assets", handlers.HandleListAssets)
	r.GET("/analytics/assets/:assetID", handlers.HandleGetAsset)
	r.GET("/analytics/views", handlers.HandleListVideoViews)
	r.GET("/analytics/views/:videoViewID", handlers.HandleGetVideoView)
	r.GET("/analytics/metrics/overall", handlers.HandleGetOverallMetrics)
	r.GET("/analytics/metrics/timeseries", handlers.HandleGetMetricTimeseries)
	r.GET("/analytics/metrics/breakdown", handlers.HandleGetMetricBreakdown)

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(func(c *gin.Context) {
		log.Printf("Unhandled route: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{"message": "Not found"})
	})
	r.Run()
}
