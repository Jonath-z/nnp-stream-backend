package middleware

import "github.com/gin-contrib/cors"

var Cors = cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://nnp-stream.vercel.app", "https://nnp-stream-dashboard.vercel.app/", "https://nnp-stream-dashboard.vercel.app"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	AllowCredentials: true,
})
