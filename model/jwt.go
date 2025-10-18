package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtClaims struct {
	Id           uuid.UUID
	Email        *string
	PhoneNumber  *string
	Role         string
	Subscription string
	jwt.StandardClaims
}

func (c JwtClaims) Valid() error {
	if c.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("token has expired")
	}
	return nil
}
