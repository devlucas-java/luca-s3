package repository

import (
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/pkg/pagination"
	"github.com/gocql/gocql"
)

type MetaDataRepository interface {
	Create(metadata *entity.MetaData) (*entity.MetaData, error)
	Updates(metadata *entity.MetaData) (*entity.MetaData, error)
	FindByID(userID gocql.UUID, id gocql.UUID) (*entity.MetaData, error)
	FindAllByUserID(userID gocql.UUID, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error)
	FindByObjectKey(objectKey string, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error)
	FindAllByStatus(userID gocql.UUID, status string, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error)
	DeleteByID(userID gocql.UUID, id gocql.UUID) error
}
