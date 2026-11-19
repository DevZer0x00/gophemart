package service_test

import (
	mock "gophemart/internal/mocks/service"
	"gophemart/internal/model"
	"gophemart/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserAuthServiceLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := mock.NewMockUserRepository(ctrl)
	hasher := mock.NewMockHasher(ctrl)
	tokenGenerator := mock.NewMockTokenGenerator(ctrl)

	userId := uint64(1)

	userRepository.EXPECT().
		FindByCredential(t.Context(), gomock.Eq("test"), gomock.Eq("password_hash")).
		MinTimes(1).MaxTimes(1).
		Return(&model.User{
			Id:           &userId,
			Username:     "test",
			PasswordHash: "password_hash",
			CreatedAt:    time.Now(),
		}, nil)
	hasher.EXPECT().Hash(gomock.Eq([]byte("password"))).Return("password_hash", nil).MinTimes(1).MaxTimes(1)
	tokenGenerator.EXPECT().GenerateToken(gomock.Any()).MaxTimes(1).Return("token", nil)

	userAuthService := service.NewUserAuthService(userRepository, hasher, tokenGenerator)
	signedToken, err := userAuthService.Login(t.Context(), "test", "password")
	require.NoError(t, err)
	assert.Equal(t, signedToken, "token")
}

func TestUserAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := mock.NewMockUserRepository(ctrl)
	hasher := mock.NewMockHasher(ctrl)
	tokenGenerator := mock.NewMockTokenGenerator(ctrl)

	userId := uint64(1)

	userRepository.EXPECT().
		Save(t.Context(), gomock.Cond(func(user *model.User) bool {
			assert.Equal(t, "test", user.Username)
			assert.Equal(t, "password_hash", user.PasswordHash)
			assert.NotEmpty(t, user.CreatedAt)

			return true
		})).
		MinTimes(1).MaxTimes(1).
		Return(&model.User{
			Id:           &userId,
			Username:     "test",
			PasswordHash: "password_hash",
			CreatedAt:    time.Now(),
		}, nil)
	hasher.EXPECT().Hash(gomock.Eq([]byte("password"))).Return("password_hash", nil).MinTimes(1).MaxTimes(1)
	tokenGenerator.EXPECT().GenerateToken(gomock.Any()).MaxTimes(1).Return("token", nil)

	userAuthService := service.NewUserAuthService(userRepository, hasher, tokenGenerator)
	signedToken, err := userAuthService.Register(t.Context(), "test", "password")
	require.NoError(t, err)
	assert.Equal(t, signedToken, "token")
}
