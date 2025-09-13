package controller

import (
	// "NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/service"
	"NotaBiz-backend/service/serviceimpl"
	// "NotaBiz-backend/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthController struct{}

var authService service.AuthService = serviceimpl.NewAuthService()

func (AuthController) GoogleLogin(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (AuthController) GoogleCallback(c *gin.Context) {
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

// func (AuthController) Register(c *gin.Context) {
// 	var RegisterRequest request.RegisterRequest
// 	if !utils.ValidateRequest(c, &RegisterRequest) {
// 		return
// 	}
// 	if RegisterRequest.Email == "" && RegisterRequest.PhoneNumber == "" {
// 		response.NewResponseBadRequest(c, "Email or phone number is required")
// 		return
// 	}
// 	res, err := authService.Register(RegisterRequest)
// 	if err != nil {
// 		response.NewResponseError(c, err.Error())
// 		return
// 	}
// 	response.NewResponseSuccess(c, res)
// }