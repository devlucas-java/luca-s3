package repository

import (
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/gocql/gocql"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
	Updates(user *entity.User) (*entity.User, error)
	FindByID(id gocql.UUID) (*entity.User, error)
	FindByEmailOrUsername(str string) (*entity.User, error)
	DeleteByID(id gocql.UUID) error
}
