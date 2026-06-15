package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"cinema/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var (
	authClient     *auth.Client
	authClientOnce sync.Once
	authClientErr  error
)

func getFirebaseAuthClient(ctx context.Context, projectID string) (*auth.Client, error) {
	authClientOnce.Do(func() {
		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, option.WithoutAuthentication())
		if err != nil {
			authClientErr = err
			return
		}
		authClient, authClientErr = app.Auth(ctx)
	})
	return authClient, authClientErr
}

func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		userID, email, err := verifyToken(c, cfg, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_email", email)
		c.Set("role", cfg.UserRole(userID))
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

func verifyToken(c *gin.Context, cfg *config.Config, token string) (string, string, error) {
	if cfg.DevMode && strings.HasPrefix(token, "dev:") {
		parts := strings.Split(token, ":")
		if len(parts) >= 2 {
			email := ""
			if len(parts) >= 3 {
				email = parts[2]
			}
			return parts[1], email, nil
		}
	}

	if cfg.FirebaseProjectID == "" {
		return "", "", errors.New("auth provider not configured")
	}

	client, err := getFirebaseAuthClient(c.Request.Context(), cfg.FirebaseProjectID)
	if err != nil {
		return "", "", err
	}

	authToken, err := client.VerifyIDToken(c.Request.Context(), token)
	if err != nil {
		return "", "", err
	}

	email, _ := authToken.Claims["email"].(string)
	return authToken.UID, email, nil
}
