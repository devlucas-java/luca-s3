package service

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/luca-s3/pkg/password_encoder"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtService     *jwt.JWTService
	mapper         *mapper.UserMapper
}

func NewAuthService(userRepository repository.UserRepository, jwtService *jwt.JWTService, mapper *mapper.UserMapper) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
		mapper:         mapper,
	}
}

func (a *AuthService) Login(login *dauth.LoginRequest) (*dauth.JWTResponse, error) {
	user, err := a.userRepository.FindByEmailOrUsername(login.Login)
	if err != nil {
		return nil, errors.ErrInvalidCredentials(err)
	}

	match, err := password_encoder.Match(login.Password, user.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify password", err)
	}

	if !match {
		return nil, errors.ErrInvalidCredentials(nil)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) Register(dto *dauth.RegisterDTO) (*dauth.JWTResponse, error) {
	user := a.mapper.RegisterDTOToUser(dto)

	pass, err := password_encoder.Encoder(dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to encode password", err)
	}

	user.Password = pass
	user.Roles = append(user.Roles, enums.USER.ToString())

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create duser", err)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) UpdatePassword(dto *dauth.UpdatePasswordRequest, auth *entity.User) error {
	user, err := a.userRepository.FindByID(auth.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve duser", err)
	}

	match, err := password_encoder.Match(dto.CurrentPassword, user.Password)
	if err != nil {
		return errors.ErrInternal("failed to verify password", err)
	}

	if !match {
		return errors.ErrInvalidCredentials(nil)
	}

	hash, err := password_encoder.Encoder(dto.NewPassword)
	if err != nil {
		return errors.ErrInternal("failed to encode password", err)
	}

	user.Password = hash

	_, err = a.userRepository.Updates(user)
	if err != nil {
		return errors.ErrDatabase("failed to update password", err)
	}

	return nil
}
