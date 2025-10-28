package controller

import (
	// "NotaBiz-backend/model/request"
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/service"
	"NotaBiz-backend/service/serviceimpl"
	"NotaBiz-backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// "NotaBiz-backend/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

var authService service.AuthService = serviceimpl.NewAuthService()

func (AuthController) GoogleLogin(c *gin.Context) {

}

func (AuthController) GoogleCallback(c *gin.Context) {

}

func (AuthController) RegisterOwner(c *gin.Context) {
	var req request.RegisterOwnerRequest
	if !utils.ValidateRequest(c, &req) {
		return
	}

	otp := utils.GenerateOTP()
	otpKey := utils.GetOTPKey(req.Email, req.PhoneNumber, req.Method)
	otpData := model.OTPData{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Attempts:  0,
		Method:    req.Method,
	}

	ctx := context.Background()
	otpJSON, _ := json.Marshal(otpData)
	err := config.Redis.Set(ctx, otpKey, otpJSON, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Redis error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	message := model.OTPMessage{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Method:      req.Method,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		OTP:         otp,
		Timestamp:   time.Now(),
		Retries:     0,
	}

	queueName := fmt.Sprintf("otp.%s", req.Method)
	if err := utils.PublishMessage(queueName, message); err != nil {
		log.Printf("Queue error: %v", err)
		// Fallback to direct sending
		go utils.ProcessOTPMessage(message)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("OTP sent successfully via %s", req.Method),
		"expires_at": otpData.ExpiresAt,
	})
}

func verifyOTPHandler(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get OTP from Redis
	otpKey := utils.GetOTPKey(req.Email, req.PhoneNumber, req.Method)
	ctx := context.Background()

	otpJSON, err := config.Redis.Get(ctx, otpKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service error"})
		return
	}

	var otpData model.OTPData
	if err := json.Unmarshal([]byte(otpJSON), &otpData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data error"})
		return
	}

	// Check attempts limit
	if otpData.Attempts >= 3 {
		config.Redis.Del(ctx, otpKey)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many failed attempts"})
		return
	}

	// Verify OTP
	if otpData.Code != req.OTP {
		otpData.Attempts++
		otpJSON, _ := json.Marshal(otpData)
		config.Redis.Set(ctx, otpKey, otpJSON, time.Until(otpData.ExpiresAt))

		c.JSON(http.StatusBadRequest, gin.H{
			"error":              "Invalid OTP",
			"attempts_remaining": 3 - otpData.Attempts,
		})
		return
	}

	// Create user in database
	user, err := createUser(req.Email, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	// Clean up OTP
	config.Redis.Del(ctx, otpKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user":    user,
	})
}

func createUser(email, phoneNumber string) (*response.UserResponse, error) {
	// Todo implement create user service
	return nil, nil
}

func NewAuthController() AuthController {

	return AuthController{}
}
