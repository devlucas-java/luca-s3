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

func (s *AuthService) Login(req *dauth.LoginRequest) (*dauth.JWTResponse, error) {
	user, err := s.userRepository.FindByEmailOrUsername(req.Login)
	if err != nil {
		return nil, errors.ErrInvalidCredentials(err)
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials(nil)
	}

	match, err := password_encoder.Match(req.Password, user.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify password", err)
	}
	if !match {
		return nil, errors.ErrInvalidCredentials(nil)
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, s.mapper.UserToUserDTO(user)), nil
}

func (s *AuthService) Register(dto *dauth.RegisterDTO) (*dauth.JWTResponse, error) {
	exists, err := s.userRepository.ExistsByEmailOrUsername(dto.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email availability", err)
	}
	if exists {
		return nil, errors.ErrConflict("email", nil)
	}

	exists, err = s.userRepository.ExistsByEmailOrUsername(dto.Username)
	if err != nil {
		return nil, errors.ErrInternal("failed to check username availability", err)
	}
	if exists {
		return nil, errors.ErrConflict("username", nil)
	}

	hash, err := password_encoder.Encoder(dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to encode password", err)
	}

	user := s.mapper.RegisterDTOToUser(dto)
	user.Password = hash
	user.Roles = []string{enums.RoleUser.String()}

	user, err = s.userRepository.Create(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create user", err)
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, s.mapper.UserToUserDTO(user)), nil
}

func (s *AuthService) UpdatePassword(dto *dauth.UpdatePasswordRequest, auth *entity.User) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return errors.ErrNotFound("user", nil)
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

	if _, err = s.userRepository.Updates(user); err != nil {
		return errors.ErrDatabase("failed to update password", err)
	}

	return nil
}

func (s *AuthService) VerifyPassword(dto *dauth.VerifyPasswordRequest, auth *entity.User) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return errors.ErrNotFound("user", nil)
	}

	match, err := password_encoder.Match(dto.Password, user.Password)
	if err != nil {
		return errors.ErrInternal("failed to verify password", err)
	}
	if !match {
		return errors.ErrInvalidCredentials(nil)
	}

	return nil
}
