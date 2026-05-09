package mapper

import (
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

type UserMapper struct {
}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) UserToUserDTO(user *entity.User) *duser.UserResponse {

	var roles []enums.Role

	for _, r := range user.Roles {
		roles = append(roles, enums.Role(r))
	}

	return &duser.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Roles:    roles,
	}
}

func (m *UserMapper) RegisterDTOToUser(dto *dauth.RegisterDTO) *entity.User {

	return &entity.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Username: dto.Username,
	}
}
