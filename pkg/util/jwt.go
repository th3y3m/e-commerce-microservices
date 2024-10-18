package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// GenerateJWT generates a JWT token for a given user ID, role, and email.
func GenerateJWT(userID int64, role, email string) (string, error) {
	// Get environment variables
	key := viper.GetString("JWT_SECRET")
	if key == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	var jwtSecret = []byte(key) // Replace with your actual secret key

	// Create the JWT claims, including user ID and expiration time
	claims := jwt.MapClaims{
		"Id":    userID,
		"Role":  role,
		"Email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeJWT decodes a JWT token and returns the user ID.
func DecodeJWT(tokenString string) (int64, error) {
	// Get the JWT secret from the environment variables
	jwtSecret := []byte(viper.GetString("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		return 0, errors.New("JWT_SECRET is not set")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used to sign the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, errors.New("invalid token: " + err.Error())
	}

	// Extract claims and verify them
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user ID from the claims
		userIDFloat, ok := claims["Id"].(float64)
		if !ok {
			return 0, errors.New("user ID not found in token")
		}
		return int64(userIDFloat), nil
	}

	return 0, errors.New("invalid token")
}
