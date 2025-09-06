package middleware

import (
	// "NotaBiz-backend/app"
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/response"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSigningMethod = jwt.SigningMethodHS256

func GenerateTokenJwt(Id, username, email string, role entity.RoleName, subscription entity.SubscriptionName, expiredAt int64) (string, error) {
	loginExpDuration := time.Duration(expiredAt) * time.Hour
	issuedAt := time.Now()
	myExpiresAt := issuedAt.Add(loginExpDuration).Unix()
	claims := model.JwtClaims{
		Id:           Id,
		Username:     username,
		Email:        email,
		Role:         string(role),
		Subscription: string(subscription),
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.GetConfig().AppConfig.Name,
			ExpiresAt: myExpiresAt,
			IssuedAt:  issuedAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwtSigningMethod, claims)

	signedToken, err := token.SignedString(config.GetConfig().AppConfig.JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

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
			return config.GetConfig().AppConfig.JwtSecret, nil
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

func AdminOwner() gin.HandlerFunc {
	return ValidateJwtAuth([]entity.RoleName{entity.RoleAdmin, entity.RoleOwner}, []entity.SubscriptionName{})
}