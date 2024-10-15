package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/oauth/model"

	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
)

type OAuthUsecase struct {
	log *logrus.Logger
}

type IOAuthUsecase interface {
	HandleOAuthUser(user goth.User) (string, error)
}

// NewOAuthUsecase creates a new OAuth service with the logger
func NewOAuthUsecase(log *logrus.Logger) IOAuthUsecase {
	return &OAuthUsecase{
		log: log,
	}
}

// HandleOAuthUser processes the OAuth user, creating them if they don't exist, and returns a JWT token
func (o *OAuthUsecase) HandleOAuthUser(user goth.User) (string, error) {

	// Call the user service to check if the user exists by their email
	url := constant.USER_SERVICE + "/get-user"

	// Create the request body
	requestBody := fmt.Sprintf(`{"email": "%s"}`, user.Email)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return "", err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to fetch user from user service: %v", err)
		return "", err
	}
	defer res.Body.Close()

	// Parse the response body into a User struct
	var existingUser model.User
	if err := json.NewDecoder(res.Body).Decode(&existingUser); err != nil {
		o.log.Errorf("Failed to decode user response: %v", err)
		return "", err
	}

	// If the user does not exist, create a new one
	if existingUser.UserID == 0 { // Ensure the check is correct, use == 0 for non-existent UserID
		// Create a new user account using the OAuth user data
		newUser := model.CreateUserRequest{
			Email: user.Email,
		}

		// Marshal the new user data to JSON
		userData, err := json.Marshal(newUser)
		if err != nil {
			o.log.Errorf("Failed to marshal new user data: %v", err)
			return "", err
		}

		// Create the new user in the user service via HTTP POST
		res, err := http.Post(constant.USER_SERVICE, "application/json", bytes.NewBuffer(userData))
		if err != nil {
			o.log.Errorf("Error creating new user: %v", err)
			return "", err
		}
		defer res.Body.Close()

		// Check for a successful status code
		if res.StatusCode != http.StatusOK {
			o.log.Errorf("Error creating user: received status %v", res.StatusCode)
			return "", errors.New("error creating user")
		}

		// Decode the response body into the `createdUser` struct
		var createdUser model.User
		if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
			o.log.Errorf("Failed to decode created user response: %v", err)
			return "", err
		}

		// Generate a JWT token for the new user
		token, err := util.GenerateJWT(createdUser.UserID, createdUser.Role, createdUser.Email)
		if err != nil {
			o.log.Errorf("Error generating token for new user: %v", err)
			return "", err
		}

		return token, nil
	}

	// If the user exists, generate a JWT token for them
	token, err := util.GenerateJWT(existingUser.UserID, existingUser.Role, existingUser.Email)
	if err != nil {
		o.log.Errorf("Error generating token for existing user: %v", err)
		return "", err
	}

	return token, nil
}
