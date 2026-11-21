package auth

//go:generate mockgen -source=user.go -destination=../../../mocks/handler/auth/mock_user.go -package=auth

import (
	"context"
	"encoding/json/v2"
	"errors"
	"gophemart/internal/http/handler"
	"gophemart/internal/repository"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/rs/zerolog"
)

type Handler struct {
	authService AuthService
	validate    *validator.Validate
	logger      *zerolog.Logger
}

type loginRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
}

type registerRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
}

type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}

func NewHandler(authService AuthService, logger *zerolog.Logger) *Handler {
	validate := validator.New()
	_ = validate.RegisterValidation("notblank", validators.NotBlank)

	return &Handler{
		authService: authService,
		validate:    validate,
		logger:      logger,
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handler.InternalError(w, h.logger, "read body", err)
			return
		}

		requestParam := &loginRequest{}
		err = json.Unmarshal(body, &requestParam)
		if err != nil {
			handler.BadRequest(w)
			return
		}

		err = h.validate.Struct(requestParam)
		if err != nil {
			handler.BadRequest(w)
			return
		}

		signedToken, err := h.authService.Login(r.Context(), requestParam.Username, requestParam.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				handler.NotAuthorized(w)
				return
			}

			handler.InternalError(w, h.logger, "login", err)
			return
		}

		handler.WithAuthToken(w, signedToken)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handler.InternalError(w, h.logger, "read body", err)
			return
		}

		requestParam := &registerRequest{}
		err = json.Unmarshal(body, &requestParam)
		if err != nil {
			handler.BadRequest(w)
			return
		}

		err = h.validate.Struct(requestParam)
		if err != nil {
			handler.BadRequest(w)
			return
		}

		signedToken, err := h.authService.Register(r.Context(), requestParam.Username, requestParam.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				handler.AlreadyExists(w)
				return
			}

			handler.InternalError(w, h.logger, "register", err)
			return
		}

		handler.WithAuthToken(w, signedToken)
		w.WriteHeader(http.StatusOK)
	}
}
