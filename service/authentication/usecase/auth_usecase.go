package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/authentication/model"
	"time"

	"github.com/sirupsen/logrus"
)

type authUsecase struct {
	log *logrus.Logger
}

type IAuthUsecase interface {
	Login(email, password string) (string, error)
	RegisterCustomer(email, password, confirmPassword string) error
	VerifyUserEmail(token string) error
}

func NewAuthUsecase(log *logrus.Logger) IAuthUsecase {
	return &authUsecase{
		log: log,
	}
}

func (o *authUsecase) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	// Call the user service to check if the user exists by their email
	url := constant.USER_SERVICE + "/get-user"

	// Create the request body
	requestBody := fmt.Sprintf(`{"email": "%s"}`, email)
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
	var user model.GetUserResponse
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		o.log.Errorf("Failed to decode user response: %v", err)
		return "", err
	}

	if !util.CheckPasswordHash(user.PasswordHash, password) {
		return "", errors.New("invalid password")
	}

	token, err := util.GenerateJWT(user.UserID, user.Role, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (o *authUsecase) RegisterCustomer(email, password, confirmPassword string) error {
	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Call the user service to check if the user exists by their email
	url := constant.USER_SERVICE + "/get-user"

	// Create the request body
	requestBody := fmt.Sprintf(`{"email": "%s"}`, email)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to fetch user from user service: %v", err)
		return err
	}
	defer res.Body.Close()

	// Parse the response body into a User struct
	var user model.GetUserResponse
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		o.log.Errorf("Failed to decode user response: %v", err)
		return err
	}

	// If the user already exists, prevent registration
	if user.Email != "" {
		return errors.New("user already exists")
	}

	var defaultVerification = false

	newUser := model.CreateUserRequest{
		Email:      email,
		Password:   hashedPassword,
		Role:       "Customer",
		Provider:   "System",
		ImageURL:   constant.DEFAULT_USER_IMAGE,
		IsVerified: &defaultVerification,
	}

	// Marshal the new user data to JSON
	userData, err := json.Marshal(newUser)
	if err != nil {
		o.log.Errorf("Failed to marshal new user data: %v", err)
		return err
	}

	// Create the new user in the user service via HTTP POST
	res, err = http.Post(constant.USER_SERVICE, "application/json", bytes.NewBuffer(userData))
	if err != nil {
		o.log.Errorf("Error creating new user: %v", err)
		return err
	}
	defer res.Body.Close()

	// Check for a successful status code
	if res.StatusCode != http.StatusOK {
		o.log.Errorf("Error creating user: received status %v", res.StatusCode)
		return errors.New("error creating user")
	}

	// Decode the response body into the `createdUser` struct
	var createdUser model.GetUserResponse
	if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
		o.log.Errorf("Failed to decode created user response: %v", err)
		return err
	}

	token, err := util.GenerateJWT(createdUser.UserID, createdUser.Role, email)
	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}

	updateUserRequest := model.UpdateUserRequest{
		UserID:       createdUser.UserID,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash,
		FullName:     createdUser.FullName,
		PhoneNumber:  createdUser.PhoneNumber,
		Address:      createdUser.Address,
		Role:         createdUser.Role,
		ImageURL:     createdUser.ImageURL,
		Token:        token,
		TokenExpires: time.Now().Add(24 * time.Hour), // Example token expiry time
		IsVerified:   createdUser.IsVerified,
		IsDeleted:    createdUser.IsDeleted,
	}

	userData, err = json.Marshal(updateUserRequest)
	if err != nil {
		o.log.Errorf("Failed to marshal updated user data: %v", err)
		return err
	}

	// Create the HTTP PUT request
	req, err = http.NewRequest("PUT", constant.USER_SERVICE, bytes.NewBuffer(userData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	res, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to update user in user service: %v", err)
		return err
	}
	defer res.Body.Close()

	// Check for a successful status code
	if res.StatusCode != http.StatusOK {
		o.log.Errorf("Error updating user: received status %v", res.StatusCode)
		return errors.New("error updating user")
	}

	// Optionally, parse the response body if needed
	var updatedUser model.GetUserResponse
	if err := json.NewDecoder(res.Body).Decode(&updatedUser); err != nil {
		o.log.Errorf("Failed to decode updated user response: %v", err)
		return err
	}

	// Send a verification email to the user
	// url = constant.MAIL_SERVICE + "/send-mail?to=" + newUser.Email + "&token=" + token

	// res, err = http.Post(url, "application/json", nil)
	// if err != nil {
	// 	o.log.Errorf("Failed to send verification email: %v", err)
	// 	return err
	// }
	// defer res.Body.Close()

	// // Check for a successful status code
	// if res.StatusCode != http.StatusOK {
	// 	o.log.Errorf("Error sending verification email: received status %v", res.StatusCode)
	// 	return errors.New("error sending verification email")
	// }

	return nil
}

func (a *authUsecase) VerifyUserEmail(token string) error {

	// Extract the user ID from the token (using your existing JWT decode logic)
	userID, err := util.DecodeJWT(token)
	if err != nil {
		return fmt.Errorf("error decoding token: %w", err)
	}

	url := constant.USER_SERVICE + "/verify?token=" + token + "&user_id=" + fmt.Sprintf("%d", userID)

	// Create the HTTP PUT request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		a.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		a.log.Errorf("Failed to verify user email: %v", err)
		return err
	}

	defer res.Body.Close()

	// Check for a successful status code
	if res.StatusCode != http.StatusOK {
		a.log.Errorf("Error verifying user email: received status %v", res.StatusCode)
		return errors.New("error verifying user email")
	}

	var resp model.VerifyTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		a.log.Errorf("Failed to decode verify token response: %v", err)
		return err
	}

	if !resp.IsValid {
		return errors.New("invalid token")
	}

	return nil
}
