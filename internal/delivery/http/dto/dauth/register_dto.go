package dauth

import "github.com/devlucas-java/luca-s3/internal/domain/errors"

type RegisterDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterDTO) Validate() error {
	if r.Name == "" {
		return errors.ErrBadRequest("name is required", nil)
	}
	if len(r.Name) > 120 {
		return errors.ErrBadRequest("name must not exceed 120 characters", nil)
	}
	if r.Email == "" {
		return errors.ErrBadRequest("email is required", nil)
	}
	if len(r.Email) > 120 {
		return errors.ErrBadRequest("email must not exceed 120 characters", nil)
	}
	if r.Username == "" {
		return errors.ErrBadRequest("username is required", nil)
	}
	if len(r.Username) > 120 {
		return errors.ErrBadRequest("username must not exceed 120 characters", nil)
	}
	if r.Password == "" {
		return errors.ErrBadRequest("password is required", nil)
	}
	if len(r.Password) < 6 {
		return errors.ErrBadRequest("password must be at least 6 characters", nil)
	}
	return nil
}
