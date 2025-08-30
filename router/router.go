package router

import (
	"NotaBiz-backend/controller"
	"NotaBiz-backend/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, apiVersion string) {
	api := r.Group(fmt.Sprintf("/api/v%s", apiVersion))
	{
		users := api.Group("/users")
		{
			users.GET("/", middleware.JwtAuthWithRoles("admin"), controller.GetUsers)
			users.POST("/register", controller.RegisterUser)
		}

		products := api.Group("/products")
		{
			products.GET("/", controller.GetProducts)
		}
	}
}