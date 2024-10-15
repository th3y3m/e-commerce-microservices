package main

import (
	"log"
	"os"
	"strings"
	"th3y3m/e-commerce-microservices/service/product/delivery"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("../../.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error while reading config file: %s", err.Error())
		return
	}
	log.Println("Config file loaded successfully")

	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}

	r := delivery.RegisterHandlers()

	log.Println("Starting server on port 8087")
	if err := r.Run(":8087"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
