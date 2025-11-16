package serviceimpl

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/middleware"
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/repository"
	"NotaBiz-backend/repository/repoimpl"
	"NotaBiz-backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/goth"
)

type AuthService struct{}

var userService UserService = UserService{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

var authRepo repository.AuthRepository = repoimpl.NewAuthRepository()

// GoogleCallback implements service.AuthService.
func (AuthService) GoogleCallback(googleUser goth.User) (res *response.LoginResponse, err error) {
	user, err := userService.GetUserByEmail(googleUser.Email)
	if err != nil {
		return nil, err
	}
	token, err := middleware.GenerateTokenJwt(user.ID, user.Email, user.PhoneNumber, entity.RoleName(user.Role.Role), entity.SubscriptionName(user.Company.Subscription.Subscription))
	if err != nil {
		return nil, err
	}

	return toLoginResponse(user, token), nil
}

// RegisterOwner implements service.AuthService.
func (AuthService) RegisterOwner(req request.RegisterOwnerRequest) (err error) {
	// register user with verified false
	user := toMasterUser(req)
	company := toMasterCompany(req)
	err = authRepo.RegisterOwner(user, company)
	if err != nil {
		return err
	}

	// send OTP for verification
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
	err = config.Redis.Set(ctx, otpKey, otpJSON, 5*time.Minute).Err()
	if err != nil {
		fmt.Printf("Redis error: %v", err)
		return err
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
		fmt.Printf("Queue error: %v", err)
		// Fallback to direct sending
		go utils.ProcessOTPMessage(message)
	}
	return nil
}

// VerifyOTP implements service.AuthService.
func (AuthService) VerifyOTP(req request.VerifyOTPRequest) (user *response.UserResponse, err error) {
	// verify OTP from redis
	otpKey := utils.GetOTPKey(req.Email, req.PhoneNumber, req.Method)
	ctx := context.Background()

	otpJSON, err := config.Redis.Get(ctx, otpKey).Result()
	if err != nil {
		return nil, fmt.Errorf("otp expired")
	}

	var otpData model.OTPData
	if err := json.Unmarshal([]byte(otpJSON), &otpData); err != nil {
		return nil, err
	}

	if otpData.Attempts >= 3 {
		config.Redis.Del(ctx, otpKey)
		return nil, fmt.Errorf("too many failed attempts")
	}

	if otpData.Code != req.OTP {
		otpData.Attempts++
		otpJSON, _ := json.Marshal(otpData)
		config.Redis.Set(ctx, otpKey, otpJSON, time.Until(otpData.ExpiresAt))

		return nil, fmt.Errorf("invalid OTP, attempts remaining: %d", 3-otpData.Attempts)
	}

	// OTP verified, activate user
	res, err := authRepo.VerifyUser(req.Email, req.PhoneNumber, req.Method)
	if err != nil {
		return nil, err
	}

	config.Redis.Del(ctx, otpKey)
	return userToUserResponse(&res), nil
}

func toLoginResponse(user *response.UserResponse, token string) *response.LoginResponse {
	return &response.LoginResponse{
		Token: token,
		User:  *user,
	}
}

func toMasterUser(req request.RegisterOwnerRequest) entity.MasterUser {
	if req.Method == "email" {
		return entity.MasterUser{
			Name:     &req.Name,
			Email:    &req.Email,
			Password: req.Password,
		}
	} else {
		return entity.MasterUser{
			Name:        &req.Name,
			PhoneNumber: &req.PhoneNumber,
			Password:    req.Password,
		}
	}
}

func toMasterCompany(req request.RegisterOwnerRequest) entity.MasterCompany {
	return entity.MasterCompany{
		Company:     req.Company,
		Description: &req.CompanyDescription,
		Address:     &req.CompanyAddress,
	}
}
