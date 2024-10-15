package main

import (
	"log"
	"os"
	"strings"
	"th3y3m/e-commerce-microservices/service/api_gateway/routes"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("../../.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error while reading config file: %s", err.Error())
		return
	}
	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}
	router := routes.SetupRouter()

	// Start the API Gateway on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("API Gateway failed to start: %v", err)
	}
}
