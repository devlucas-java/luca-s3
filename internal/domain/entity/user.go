package entity

import (
	"fmt"
	"slices"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/pkg/password_encoder"
	"github.com/gocql/gocql"
)

type User struct {
	ID       gocql.UUID `db:"id"`
	Name     string     `db:"name"`
	Email    string     `db:"email"`
	Username string     `db:"username"`
	Password string     `db:"password"`

	Roles []string `db:"roles"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(name, email, username, pass string) (*User, error) {
	hash, err := password_encoder.Encoder(pass)
	if err != nil {
		return nil, err
	}

	uuid, err := gocql.RandomUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user uuid: %w", err)
	}

	now := time.Now()
	return &User{
		ID:        uuid,
		Name:      name,
		Email:     email,
		Username:  username,
		Password:  hash,
		Roles:     []string{enums.RoleUser.String()},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) HasRole(role enums.Role) bool {
	return slices.Contains(u.Roles, role.String())
}
