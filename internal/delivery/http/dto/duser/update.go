package duser

import "github.com/devlucas-java/luca-s3/internal/domain/errors"

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Name == "" && r.Email == "" && r.Username == "" {
		return errors.ErrBadRequest("at least one field (name, email or username) must be provided", nil)
	}
	if r.Name != "" && len(r.Name) > 120 {
		return errors.ErrBadRequest("name must not exceed 120 characters", nil)
	}
	if r.Email != "" && len(r.Email) > 120 {
		return errors.ErrBadRequest("email must not exceed 120 characters", nil)
	}
	if r.Username != "" && len(r.Username) > 120 {
		return errors.ErrBadRequest("username must not exceed 120 characters", nil)
	}
	return nil
}
