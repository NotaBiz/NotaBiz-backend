// Package router defines the application's routing structure and initializes
// the API endpoints. It uses the Gin framework for HTTP routing and middleware.
package router

import (
	"NotaBiz-backend/controller"
	"NotaBiz-backend/middleware"
	"NotaBiz-backend/model/entity"
	"fmt"

	"github.com/gin-gonic/gin"
)

// InitRouter initializes the application's routes and API endpoints.
//
// Parameters:
//   - r: The Gin engine instance.
//   - apiVersion: The version of the API to be included in the route paths.
func InitRouter(r *gin.Engine, apiVersion string) {
	api := r.Group(fmt.Sprintf("/api/v%s", apiVersion))
	{
		api.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "pong"})
		})

		users := api.Group("/users")
		{
			users.GET("/", middleware.ValidateJwtAuth([]entity.RoleName{entity.RoleAdmin, entity.RoleOwner}, []entity.SubscriptionName{}), controller.GetUsers)
			users.POST("/register", controller.RegisterUser)
		}

		products := api.Group("/products")
		{
			products.GET("/", controller.GetProducts)
		}
	}
}
