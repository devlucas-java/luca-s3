package mapper

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) UserToUserDTO(user *entity.User) *duser.UserResponse {
	return &duser.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Username:  user.Username,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func (m *UserMapper) RegisterDTOToUser(dto *dauth.RegisterDTO) *entity.User {
	return &entity.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Username: dto.Username,
	}
}
