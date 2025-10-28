// Package router defines the application's routing structure and initializes
// the API endpoints. It uses the Gin framework for HTTP routing and middleware.
package router

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/controller"
	"NotaBiz-backend/middleware"
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/service"
	"NotaBiz-backend/service/serviceimpl"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

// InitRouter initializes the application's routes and API endpoints.
//
// Parameters:
//   - r: The Gin engine instance.
//   - apiVersion: The version of the API to be included in the route paths.
func InitRouter(r *gin.Engine) {
	appConfig := config.Data.AppConfig
	apiGroup := fmt.Sprintf("/api/v%s", string(appConfig.Version[0]))
	goth.UseProviders(
		google.New(
			appConfig.GoogleClientID,
			appConfig.GoogleClientSecret,
			fmt.Sprintf("http://localhost:%s%s/auth/google/callback", strconv.Itoa(appConfig.Port), apiGroup),
			"email", "profile",
		),
	)
	api := r.Group(apiGroup)
	{
		api.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "pong"})
		})
		api.GET("/health", healthHandler)

		// authController := controller.NewAuthController()
		auth := api.Group("/auth")
		{
			auth.GET("/:provider", func(c *gin.Context) {
				provider := c.Param("provider")
				// ⚡ Inject provider into request context
				req := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))
				c.Request = req
				fmt.Println("provider", c.Param("provider"))
				gothic.BeginAuthHandler(c.Writer, c.Request)
			})
			auth.GET("/:provider/callback", func(c *gin.Context) {
				provider := c.Param("provider")
				// ⚡ Inject provider into request context
				req := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))
				c.Request = req
				user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
				if err != nil {
					log.Println("Goole auth failed: ", err)
					response.NewResponseUnauthorized(c, err.Error())
					return
				}
				fmt.Println("user", user)

				var authService service.AuthService = serviceimpl.NewAuthService()
				res, err := authService.GoogleCallback(user)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusForbidden, response.Response{
							Status:    response.StatusError,
							Code:      http.StatusForbidden,
							Message:   "user not found",
							Data:      user,
							Timestamp: time.Now(),
						})
						return
					}
					response.NewResponseError(c, err.Error())
					return
				}
				response.NewResponseSuccess(c, res)
			})
		}

		users := api.Group("/users")
		{
			users.GET("/", middleware.ValidateJwtAuth([]entity.RoleName{entity.RoleAdmin, entity.RoleOwner}, []entity.SubscriptionName{}), controller.GetUsers)
		}

		products := api.Group("/products")
		{
			products.GET("/", controller.GetProducts)
		}
	}
}

func healthHandler(c *gin.Context) {
	status := gin.H{
		"status":   "healthy",
		"services": gin.H{},
	}

	// Check Redis
	ctx := context.Background()
	if err := config.Redis.Ping(ctx).Err(); err != nil {
		status["services"].(gin.H)["redis"] = "unhealthy"
	} else {
		status["services"].(gin.H)["redis"] = "healthy"
	}

	// Check Database
	conn, _ := config.DB.DB()
	if err := conn.Ping(); err != nil {
		status["services"].(gin.H)["database"] = "unhealthy"
	} else {
		status["services"].(gin.H)["database"] = "healthy"
	}

	// Check RabbitMQ
	if config.RabbitMQ.IsClosed() {
		status["services"].(gin.H)["rabbitmq"] = "unhealthy"
	} else {
		status["services"].(gin.H)["rabbitmq"] = "healthy"
	}

	c.JSON(http.StatusOK, status)
}
