package service

//go:generate mockgen -source=user_auth.go -destination=../mocks/service/mock_user_auth.go -package=service

import (
	"context"
	"gophemart/internal/model"
	"time"
)

// UserAuthService сервис регистрации/авторизации и обработки токенов
type UserAuthService struct {
	repo           UserRepository
	hasher         Hasher
	tokenGenerator TokenGenerator
}

type UserRepository interface {
	FindByCredential(ctx context.Context, login, passwordHash string) (*model.User, error)
	Save(ctx context.Context, model *model.User) (*model.User, error)
}

type Hasher interface {
	Hash(data []byte) (string, error)
}

type TokenGenerator interface {
	GenerateToken(customData *CustomData) (string, error)
	ParseSignedToken(signedToken string) (*CustomData, error)
}

func NewUserAuthService(repo UserRepository, hasher Hasher, tokenGenerator TokenGenerator) *UserAuthService {
	return &UserAuthService{
		repo:           repo,
		hasher:         hasher,
		tokenGenerator: tokenGenerator,
	}
}

// Register регистрирует нового пользователя по переданным username, password - возвращает авторизационный токен при успехе, либо ошибку
func (s *UserAuthService) Register(ctx context.Context, username, password string) (string, error) {
	passwordHash, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return "", err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	user, err = s.repo.Save(ctx, user)
	if err != nil {
		return "", err
	}

	signedToken, err := s.tokenGenerator.GenerateToken(&CustomData{
		UserId: *user.Id,
	})

	return signedToken, err
}

// Login авторизует нового пользователя по переданным username, password - возвращает авторизационный токен при успехе, либо ошибку
func (s *UserAuthService) Login(ctx context.Context, username, password string) (string, error) {
	passwordHash, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return "", err
	}

	user, err := s.repo.FindByCredential(ctx, username, passwordHash)
	if err != nil {
		return "", err
	}

	signedToken, err := s.tokenGenerator.GenerateToken(&CustomData{
		UserId: *user.Id,
	})

	return signedToken, err
}

// GetUserIdFromBearer возвращает id пользователя из переданного signedToken либо ошибку
func (s *UserAuthService) GetUserIdFromBearer(ctx context.Context, signedToken string) (uint64, error) {
	tokenData, err := s.tokenGenerator.ParseSignedToken(signedToken)
	if err != nil {
		return 0, err
	}

	return tokenData.UserId, nil
}
