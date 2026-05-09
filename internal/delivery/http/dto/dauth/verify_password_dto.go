package dauth

import "github.com/devlucas-java/luca-s3/internal/domain/errors"

type VerifyPasswordRequest struct {
	Password string `json:"password"`
}

func (r *VerifyPasswordRequest) Validate() error {
	if r.Password == "" {
		return errors.ErrBadRequest("password is required", nil)
	}
	return nil
}
