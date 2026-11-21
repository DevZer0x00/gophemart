package auth_test

import (
	"bytes"
	"gophemart/internal/http/handler/auth"
	mock "gophemart/internal/mocks/handler/auth"
	"gophemart/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name                 string
		body                 string
		username             string
		password             string
		signedToken          string
		err                  error
		expectedResponseCode int
	}{
		{
			name:                 "success",
			body:                 `{"username": "test", "password": "test"}`,
			username:             "test",
			password:             "test",
			signedToken:          "signedToken1",
			err:                  nil,
			expectedResponseCode: http.StatusOK,
		},
		{
			name:                 "bad json",
			body:                 `{"username"": "test", "password": "test"}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "empty username",
			body:                 `{"username": "   ", "password": "test"}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "empty password",
			body:                 `{"username": "test", "password": "     "}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "user not found",
			body:                 `{"username": "test", "password": "test"}`,
			username:             "test",
			password:             "test",
			err:                  repository.ErrUserNotFound,
			expectedResponseCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authService := mock.NewMockAuthService(ctrl)
			authService.EXPECT().Login(gomock.Any(), tt.username, tt.password).MaxTimes(1).Return(tt.signedToken, tt.err)
			logger := zerolog.New(io.Discard)

			handler := auth.NewHandler(authService, &logger)
			request := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()

			handler.Login().ServeHTTP(w, request)
			resp := w.Result()

			assert.Equal(t, tt.expectedResponseCode, resp.StatusCode)

			if tt.signedToken != "" && tt.err == nil {
				assert.Equal(t, "Bearer "+tt.signedToken, resp.Header.Get("Authorization"))
			}
		})
	}
}

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		name                 string
		body                 string
		username             string
		password             string
		signedToken          string
		err                  error
		expectedResponseCode int
	}{
		{
			name:                 "success",
			body:                 `{"username": "test", "password": "test"}`,
			username:             "test",
			password:             "test",
			signedToken:          "signedToken1",
			err:                  nil,
			expectedResponseCode: http.StatusOK,
		},
		{
			name:                 "bad json",
			body:                 `{"username"": "test", "password": "test"}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "empty username",
			body:                 `{"username": "   ", "password": "test"}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "empty password",
			body:                 `{"username": "test", "password": "     "}`,
			err:                  nil,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "duplicate user",
			body:                 `{"username": "test", "password": "test"}`,
			username:             "test",
			password:             "test",
			err:                  repository.ErrUserAlreadyExists,
			expectedResponseCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authService := mock.NewMockAuthService(ctrl)
			authService.EXPECT().Register(gomock.Any(), tt.username, tt.password).MaxTimes(1).Return(tt.signedToken, tt.err)
			logger := zerolog.New(io.Discard)

			handler := auth.NewHandler(authService, &logger)
			request := httptest.NewRequest(http.MethodGet, "/register", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()

			handler.Register().ServeHTTP(w, request)
			resp := w.Result()

			assert.Equal(t, tt.expectedResponseCode, resp.StatusCode)

			if tt.signedToken != "" && tt.err == nil {
				assert.Equal(t, "Bearer "+tt.signedToken, resp.Header.Get("Authorization"))
			}
		})
	}
}
