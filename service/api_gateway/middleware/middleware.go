package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// AuthMiddleware validates JWT and applies Casbin authorization
func AuthMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT secret from the environment variables
		jwtSecret := []byte(viper.GetString("JWT_SECRET"))
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims and role
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["Role"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		role := claims["Role"].(string)
		obj := c.Request.URL.Path
		act := c.Request.Method

		// Casbin authorization
		allowed, err := enforcer.Enforce(role, obj, act)
		if err != nil {
			fmt.Printf("Error during Casbin enforcement: %v\n", err)
		}
		if !allowed {
			fmt.Printf("Access denied for role: %s, object: %s, action: %s\n", role, obj, act)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Proceed to the next handler if authorized
		c.Next()
	}
}
