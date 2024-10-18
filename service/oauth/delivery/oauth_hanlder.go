package delivery

import (
	"log"
	"net/http"
	"os"
	"th3y3m/e-commerce-microservices/service/oauth/dependency_injection"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
)

// JWTResponse represents the response containing a JWT token.
type JWTResponse struct {
	Token string `json:"token"`
}

// InitializeOAuth initializes the OAuth providers and sets up the session store.
func InitializeOAuth() {
	// Set up environment variables
	clientID := viper.GetString("GOOGLE_CLIENT_ID")
	clientSecret := viper.GetString("GOOGLE_CLIENT_SECRET")
	facebookID, facebookSecret := viper.GetString("FACEBOOK_CLIENT_ID"), viper.GetString("FACEBOOK_CLIENT_SECRET")

	// Check if client ID and secret are set
	if clientID == "" || clientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET environment variables are not set")
	}

	if facebookID == "" || facebookSecret == "" {
		log.Fatal("FACEBOOK_CLIENT_ID or FACEBOOK_CLIENT_SECRET environment variables are not set")
	}

	sessionSecret := viper.GetString("SESSION_SECRET")

	// Check if session secret is set
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable is not set")
	}

	// Set SESSION_SECRET for Gothic
	os.Setenv("SESSION_SECRET", sessionSecret)

	// Initialize Goth with the Google and Facebook providers
	goth.UseProviders(
		google.New(clientID, clientSecret, "http://localhost:8080/auth/google/callback"),
		facebook.New(facebookID, facebookSecret, "http://localhost:8080/auth/facebook/callback"),
	)

	// Explicitly set Gothic store to use the session store from Gin
	key := []byte(sessionSecret)
	store := cookie.NewStore(key, key)
	gothic.Store = store
}

func GoogleLogin(c *gin.Context) {
	c.Request.URL.RawQuery = "provider=google"
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GoogleCallback(c *gin.Context) {
	// Complete the user authentication with Gothic
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Error during authentication: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Google"})
		return
	}
	service := dependency_injection.NewOAuthUsecaseProvider()

	// Handle Google user and generate JWT token
	token, err := service.HandleOAuthUserGoogle(user)
	if err != nil {
		log.Printf("Error handling Google user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle Google user"})
		return
	}

	// Respond with the generated JWT token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GoogleLogout(c *gin.Context) {
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		log.Printf("Error logging out: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func FacebookLogin(c *gin.Context) {
	c.Request.URL.RawQuery = "provider=facebook"
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func FacebookCallback(c *gin.Context) {
	// Complete the user authentication with Gothic
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Error during authentication: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Facebook"})
		return
	}
	service := dependency_injection.NewOAuthUsecaseProvider()

	// Handle Facebook user and generate JWT token
	token, err := service.HandleOAuthUserFacebook(user)
	if err != nil {
		log.Printf("Error handling Facebook user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle Facebook user"})
		return
	}

	// Respond with the generated JWT token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func FacebookLogout(c *gin.Context) {
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		log.Printf("Error logging out: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
