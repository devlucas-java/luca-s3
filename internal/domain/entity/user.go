package entity

import (
	"fmt"
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
	now := time.Now()
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return nil, fmt.Errorf("failed generate user uuid: %v", err)
	}
	return &User{
		ID:        uuid,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Email:     email,
		Username:  username,
		Password:  hash,
		Roles:     []string{enums.USER.ToString()},
	}, nil
}

func (u *User) HasRole(role enums.Role) bool {
	for _, r := range u.Roles {
		if r == role.ToString() {
			return true
		}
	}
	return false
}
