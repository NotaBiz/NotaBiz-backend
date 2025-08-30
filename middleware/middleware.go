package middleware

import (
	// "NotaBiz-backend/app"
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSigningMethod = jwt.SigningMethodHS256

func GenerateTokenJwt(Id, name, email, role string, expiredAt int64) (string, error) {
	loginExpDuration := time.Duration(expiredAt) * time.Hour
	issuedAt := time.Now()
	myExpiresAt := issuedAt.Add(loginExpDuration).Unix()
	claims := model.JwtClaims{
		Id:       Id,
		Username: name,
		Email:    email,
		Role:     role,
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

func JwtAuthWithRoles(userId ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response.NewResponseUnauthorized(c, "Invalid authorization header token")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
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
		if len(userId) > 0 {
			for _, role := range userId {
				if role == claims.Role {
					validRole = true
					break
				}
			}
		}
		if !validRole {
			response.NewResponseUnauthorized(c, "Invalid role")
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
