package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/application/service"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/middleware"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/response"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
	"github.com/go-chi/chi/v5"
	"github.com/gocql/gocql"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.userService.GetMe(auth)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) error {
	var req duser.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	res, err := h.userService.UpdateMe(auth, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	if err := h.userService.DeleteMe(auth); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
	return nil
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	id, err := gocql.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	res, err := h.userService.GetByID(auth, id)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) UpdateByID(w http.ResponseWriter, r *http.Request) error {
	var req duser.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	id, err := gocql.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	res, err := h.userService.UpdateByID(auth, id, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	id, err := gocql.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	if err := h.userService.DeleteByID(auth, id); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
	return nil
}
