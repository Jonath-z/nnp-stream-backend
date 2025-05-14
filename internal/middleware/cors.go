package middleware

import "github.com/gin-contrib/cors"

var Cors = cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:3000", "https://nnp-stream.vercel.app", "https://nnp-dashboard.vercel.app"},
	AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	AllowCredentials: true,
})
