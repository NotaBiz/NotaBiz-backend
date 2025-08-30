package controller

import (
	"NotaBiz-backend/model/response"

	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	response.NewResponseSuccessPaging(ctx, nil, nil)
}

func RegisterUser(ctx *gin.Context) {
	response.NewResponseCreated(ctx, "register users")
}