package middleware

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/golang-jwt/jwt"
)

func MaxSizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
		if err := c.Request.ParseForm(); err != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request too large",
			})
			return
		}
		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var key string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &model.JwtClaims{}
			_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return config.Data.AppConfig.JwtSecret, nil
			})
			if err != nil {
				response.NewResponseUnauthorized(c, "Invalid token")
				c.Abort()
				return
			}
			key = "jwt:" + claims.Id.String()
		} else {
			key = "ip:" + c.ClientIP()
		}

		var rateLimit redis_rate.Limit
		if strings.HasPrefix(key, "ip:") {
			rateLimit = redis_rate.PerSecond(4)
		} else {
			rateLimit = redis_rate.PerSecond(1)
		}

		// Request a token from Redis bucket
		res, err := config.Limiter.Allow(ctx, key, rateLimit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error " + err.Error()})
			return
		}

		if res.Allowed == 0 {
			c.Header("Retry-After", res.ResetAfter.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": res.ResetAfter.String(),
			})
			return
		}

		c.Next()
	}
}
