package dauth

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/duser"
)

type JWTResponse struct {
	Token     string              `json:"token"`
	TypeToken string              `json:"typeToken"`
	User      *duser.UserResponse `json:"user_response"`
}

func NewJWTResponse(token string, user *duser.UserResponse) *JWTResponse {
	return &JWTResponse{
		Token:     token,
		TypeToken: "Bearer",
		User:      user,
	}
}
