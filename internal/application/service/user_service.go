package service

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/gocql/gocql"
)

type UserService struct {
	userRepository repository.UserRepository
	userMapper     *mapper.UserMapper
}

func NewUserService(repo repository.UserRepository, mapper *mapper.UserMapper) *UserService {
	return &UserService{
		userRepository: repo,
		userMapper:     mapper,
	}
}

// GetMe returns the authenticated user's profile.
func (s *UserService) GetMe(auth *entity.User) (*duser.UserResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound("user", nil)
	}
	return s.userMapper.UserToUserDTO(user), nil
}

// UpdateMe updates the authenticated user's own profile.
func (s *UserService) UpdateMe(auth *entity.User, dto *duser.UpdateUserRequest) (*duser.UserResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound("user", nil)
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Email != "" {
		existing, err := s.userRepository.FindByEmailOrUsername(dto.Email)
		if err != nil {
			return nil, errors.ErrInternal("failed to check email availability", err)
		}
		if existing != nil && existing.ID != auth.ID {
			return nil, errors.ErrConflict("email", nil)
		}
		user.Email = dto.Email
	}
	if dto.Username != "" {
		existing, err := s.userRepository.FindByEmailOrUsername(dto.Username)
		if err != nil {
			return nil, errors.ErrInternal("failed to check username availability", err)
		}
		if existing != nil && existing.ID != auth.ID {
			return nil, errors.ErrConflict("username", nil)
		}
		user.Username = dto.Username
	}

	saved, err := s.userRepository.Updates(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to update user", err)
	}
	return s.userMapper.UserToUserDTO(saved), nil
}

// DeleteMe removes the authenticated user's own account.
func (s *UserService) DeleteMe(auth *entity.User) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return errors.ErrNotFound("user", nil)
	}

	if err := s.userRepository.DeleteByID(auth.ID); err != nil {
		return errors.ErrDatabase("failed to delete user", err)
	}
	return nil
}

// GetByID returns any user by ID. Admin only.
func (s *UserService) GetByID(auth *entity.User, id gocql.UUID) (*duser.UserResponse, error) {
	if !auth.HasRole(enums.RoleAdmin) {
		return nil, errors.ErrForbidden(nil)
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound("user", nil)
	}
	return s.userMapper.UserToUserDTO(user), nil
}

// UpdateByID updates any user by ID. Admin only.
func (s *UserService) UpdateByID(auth *entity.User, id gocql.UUID, dto *duser.UpdateUserRequest) (*duser.UserResponse, error) {
	if !auth.HasRole(enums.RoleAdmin) {
		return nil, errors.ErrForbidden(nil)
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound("user", nil)
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Username != "" {
		user.Username = dto.Username
	}

	saved, err := s.userRepository.Updates(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to update user", err)
	}
	return s.userMapper.UserToUserDTO(saved), nil
}

// DeleteByID removes any user by ID. Admin only.
func (s *UserService) DeleteByID(auth *entity.User, id gocql.UUID) error {
	if !auth.HasRole(enums.RoleAdmin) {
		return errors.ErrForbidden(nil)
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return errors.ErrInternal("failed to retrieve user", err)
	}
	if user == nil {
		return errors.ErrNotFound("user", nil)
	}

	if err := s.userRepository.DeleteByID(id); err != nil {
		return errors.ErrDatabase("failed to delete user", err)
	}
	return nil
}
