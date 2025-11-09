package controller

import (
	// "NotaBiz-backend/model/request"
	"NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/service"
	"NotaBiz-backend/service/serviceimpl"
	"NotaBiz-backend/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	// "NotaBiz-backend/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthController struct{}

var authService service.AuthService = serviceimpl.NewAuthService()

// GoogleLogin godoc
// @Summary      OAuth2 login redirect
// @Description  Initiate OAuth2 login for the specified provider (e.g. google). Redirects the client to the provider's consent page.
// @Tags         Auth
// @Success      200  
// @Router       /auth/google [get]
func (AuthController) GoogleLogin(c *gin.Context) {
	provider := c.Param("provider")
	req := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))
	c.Request = req
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// GoogleCallback godoc
// @Summary      OAuth2 callback
// @Description  OAuth2 callback endpoint. Exchanges code for user info and performs login/registration. Returns 200 on success, 401 if not authorized, 500 on server error.
// @Tags         Auth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/google/callback [get]
func (AuthController) GoogleCallback(c *gin.Context) {
	provider := c.Param("provider")
	req := c.Request.WithContext(context.WithValue(c.Request.Context(), gothic.ProviderParamKey, provider))
	c.Request = req
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("Goole auth failed: ", err)
		response.NewResponseUnauthorized(c, err.Error())
		return
	}

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
}

// RegisterOwner godoc
// @Summary      Register owner
// @Description  Register owner by email and phone number
// @Tags         Auth
// @Param        request body request.RegisterOwnerRequest true "Request body"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/register [post]
func (AuthController) RegisterOwner(c *gin.Context) {
	var req request.RegisterOwnerRequest
	if !utils.ValidateRequest(c, &req) {
		return
	}

	err := authService.RegisterOwner(req)
	if err != nil {
		response.NewResponseError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("OTP sent successfully via %s", req.Method),
	})
}

// VerifyOTPHandler godoc
// @Summary      Verify OTP
// @Description  Verify OTP code for registration
// @Tags         Auth
// @Param        request body request.VerifyOTPRequest true "Request body"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/verify-otp [post]
func (AuthController) VerifyOTPHandler(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := authService.VerifyOTP(req)
	if err != nil {
		response.NewResponseError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user":    user,
	})
}

func NewAuthController() AuthController {

	return AuthController{}
}
