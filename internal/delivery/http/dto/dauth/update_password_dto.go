package dauth

import "github.com/devlucas-java/luca-s3/internal/domain/errors"

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (r *UpdatePasswordRequest) Validate() error {
	if r.CurrentPassword == "" {
		return errors.ErrBadRequest("current_password is required", nil)
	}
	if r.NewPassword == "" {
		return errors.ErrBadRequest("new_password is required", nil)
	}
	if len(r.NewPassword) < 6 {
		return errors.ErrBadRequest("new_password must be at least 6 characters", nil)
	}
	if r.CurrentPassword == r.NewPassword {
		return errors.ErrBadRequest("new_password must differ from current_password", nil)
	}
	return nil
}
