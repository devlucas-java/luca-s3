package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/application/service"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/middleware"
	"github.com/devlucas-java/luca-s3/internal/delivery/http/response"
	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/domain/errors"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req dauth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.authService.Login(&req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req dauth.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.authService.Register(&req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	var req dauth.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	user := r.Context().Value(middleware.AuthKey).(*entity.User)
	if err := h.authService.UpdatePassword(&req, user); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
	return nil
}
