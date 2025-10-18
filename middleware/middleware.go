// Package middleware provides middleware functions for handling JWT authentication
// and authorization in the application. It includes functionality for generating
// JWT tokens and validating them based on roles and subscriptions.
package middleware

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/response"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var jwtSigningMethod = jwt.SigningMethodHS256

// GenerateTokenJwt generates a JWT token for a user with the specified details.
//
// Parameters:
//   - Id: The unique identifier of the user.
//   - email: The email address of the user.
//   - role: The role of the user (e.g., Admin, User).
//   - subscription: The subscription level of the user.
//   - expiredAt: The expiration time of the token in hours.
//
// Returns:
//   - A string containing the signed JWT token.
//   - An error if the token generation fails.
func GenerateTokenJwt(Id uuid.UUID, email, phoneNumber *string, role entity.RoleName, subscription entity.SubscriptionName) (string, error) {
	loginExpDuration := time.Duration(config.Data.AppConfig.JwtExpiration) * time.Second
	issuedAt := time.Now()
	myExpiresAt := issuedAt.Add(loginExpDuration).Unix()
	claims := model.JwtClaims{
		Id:           Id,
		Email:        email,
		PhoneNumber:  phoneNumber,
		Role:         string(role),
		Subscription: string(subscription),
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.Data.AppConfig.Name,
			ExpiresAt: myExpiresAt,
			IssuedAt:  issuedAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwtSigningMethod, claims)

	signedToken, err := token.SignedString(config.Data.AppConfig.JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateJwtAuth is a middleware function that validates the JWT token in the
// Authorization header of incoming requests. It checks the token's validity,
// role, and subscription level.
//
// Parameters:
//   - roles: A list of allowed roles for the request.
//   - subscription: A list of allowed subscription levels for the request.
//
// Returns:
//   - A Gin middleware handler function that validates the JWT token and
//     aborts the request with an unauthorized response if validation fails.
func ValidateJwtAuth(roles []entity.RoleName, subscription []entity.SubscriptionName) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response.NewResponseUnauthorized(c, "Invalid authorization header token")
			c.Abort()
			return
		}

		tokenString := strings.ReplaceAll(authHeader, "Bearer ", "")
		claims := &model.JwtClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Data.AppConfig.JwtSecret, nil
		})
		if err != nil {
			response.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}
		if !token.Valid {
			response.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		validRole := false
		if len(roles) > 0 {
			if slices.Contains(roles, entity.RoleName(claims.Role)) {
				validRole = true
			}
		}
		if !validRole {
			response.NewResponseUnauthorized(c, "Invalid role")
			c.Abort()
			return
		}

		validSubscription := false
		if len(subscription) > 0 {
			if slices.Contains(subscription, entity.SubscriptionName(claims.Subscription)) {
				validSubscription = true
			}
		}
		if !validSubscription {
			response.NewResponseUnauthorized(c, "Invalid subscription")
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}