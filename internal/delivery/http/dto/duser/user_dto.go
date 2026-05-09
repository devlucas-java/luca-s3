package duser

import (
	"github.com/devlucas-java/luca-s3/internal/domain/enums"
)

type UserResponse struct {
	ID        string       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string       `json:"name" example:"John Doe"`
	Email     string       `json:"email" example:"john@example.com"`
	Username  string       `json:"username" example:"johndoe"`
	IsSeller  bool         `json:"is_seller" example:"false"`
	Roles     []enums.Role `json:"roles" example:"["USER"]"`
	CreatedAt string       `json:"created_at,omitempty"`
	UpdatedAt string       `json:"updated_at,omitempty"`
}
