package main

import (
	"log"
	"th3y3m/e-commerce-microservices/service/api_gateway/handler"
	"th3y3m/e-commerce-microservices/service/api_gateway/logging"
	"th3y3m/e-commerce-microservices/service/api_gateway/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	enforcer, err := casbin.NewEnforcer("rbac/rbac_model.conf", "rbac/rbac_policy.csv")
	if err != nil {
		log.Fatalf("Failed to load Casbin model and policy: %v", err)
	}

	r := gin.Default()

	// Set up CORS middleware (optional, if needed)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set up Auth middleware with Casbin
	r.Use(middleware.AuthMiddleware(enforcer))

	// Use rate limiting, logging, and caching
	r.Use(gin.WrapF(middleware.CacheMiddleware(logging.LogRequests(middleware.RateLimit(handler.RouteHandler)))))

	log.Println("API Gateway running on port 9000...")
	if err := r.Run(":9000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
