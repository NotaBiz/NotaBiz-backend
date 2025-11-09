package model

import (
	"time"
)

type OTPMessage struct {
	ID          string    `json:"id"`
	Method      string    `json:"method"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	OTP         string    `json:"otp"`
	Timestamp   time.Time `json:"timestamp"`
	Retries     int       `json:"retries"`
}

type OTPData struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Attempts  int       `json:"attempts"`
	Method    string    `json:"method"`
}

type WhatsAppResponse struct {
	Reason string `json:"reason"`
	Status bool `json:"status"`
}