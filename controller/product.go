package controller

import (
	"NotaBiz-backend/model/response"

	"github.com/gin-gonic/gin"
)

func GetProducts(ctx *gin.Context) {
	response.NewResponseSuccess(ctx, "products")
}