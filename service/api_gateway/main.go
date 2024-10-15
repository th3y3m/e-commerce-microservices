package main

import (
	"log"
	"th3y3m/e-commerce-microservices/service/api_gateway/routes"
)

func main() {
	router := routes.SetupRouter()

	// Start the API Gateway on port 8080
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("API Gateway failed to start: %v", err)
	}
}
