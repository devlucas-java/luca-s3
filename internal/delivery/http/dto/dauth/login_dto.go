package dauth

import "github.com/devlucas-java/luca-s3/internal/domain/errors"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Login == "" {
		return errors.ErrBadRequest("login is required", nil)
	}
	if r.Password == "" {
		return errors.ErrBadRequest("password is required", nil)
	}
	return nil
}
