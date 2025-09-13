package controller

import (
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/service"
	"NotaBiz-backend/service/serviceimpl"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthController struct{}

var authService service.AuthService = serviceimpl.NewAuthService()

func GoogleLogin(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GoogleCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("Goole auth failed: ", err)
		response.NewResponseUnauthorized(c, err.Error())
		return
	}

	res, err := authService.GoogleCallback(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NewResponseUnauthorized(c, "user not found")
			return
		}
		response.NewResponseError(c, err.Error())
		return
	}
	response.NewResponseSuccess(c, res)
}