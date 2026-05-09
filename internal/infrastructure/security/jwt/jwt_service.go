package jwt

import (
	"errors"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/go-chi/jwtauth"
)

type JWTService struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewJWTService(secret string) *JWTService {
	j := jwtauth.New("HS256", []byte(secret), nil)
	return &JWTService{tokenAuth: j}
}

func (s *JWTService) GenerateToken(user *entity.User) (string, error) {

	_, token, err := s.tokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"roles":   user.Roles,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *JWTService) Validate(tokenString string) (map[string]interface{}, error) {
	token, err := jwtauth.VerifyToken(s.tokenAuth, tokenString)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("invalid token")
	}

	claims := token.PrivateClaims()

	return claims, nil
}
