package repository

import (
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/pkg/id"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
	Updates(user *entity.User) (*entity.User, error)
	FindByID(id id.UUID) (*entity.User, error)
	FindByEmailOrUsername(str string) (*entity.User, error)
	DeleteByID(id id.UUID) error
}
