package repository

import (
	"context"
	"database/sql"
	"errors"
	"gophemart/internal/helper"
	"gophemart/internal/model"
	"time"

	"github.com/rs/zerolog"
)

var ErrUserAlreadyExists = errors.New("user already exists")

// UserRepository репозиторий модели model.User
type UserRepository struct {
	db              *sql.DB
	logger          *zerolog.Logger
	retryableHelper *helper.RetryableHelper
}

// NewUserRepository создает и возвращает ссылку на экземпляр репозитория
func NewUserRepository(db *sql.DB, logger *zerolog.Logger, retryableHelper *helper.RetryableHelper) *UserRepository {
	return &UserRepository{
		db:              db,
		logger:          logger,
		retryableHelper: retryableHelper,
	}
}

// FindByCredential ищет и возвращает ссылку на пользователя, если найдено совпадение - иначе nil
func (rep *UserRepository) FindByCredential(ctx context.Context, login, passwordHash string) (*model.User, error) {
	user := &model.User{
		Id: new(uint64),
	}
	err := rep.retryableHelper.Run(ctx, func() error {
		row := rep.db.QueryRowContext(ctx, "SELECT id, username, password_hash, created_at FROM users WHERE username = $1 AND password_hash = $2", login, passwordHash)

		return row.Scan(user.Id, &user.Username, &user.PasswordHash, &user.CreatedAt)
	})

	return user, err
}

// Save сохраняет модель пользователя в базе
func (rep *UserRepository) Save(ctx context.Context, model *model.User) (*model.User, error) {
	var id uint64
	if model.Id != nil {
		id = *model.Id
	}

	err := rep.retryableHelper.Run(ctx, func() error {
		var err error

		if id != 0 { // обновляем данные
			_, err = rep.db.ExecContext(ctx, "UPDATE users SET password_hash=$1 WHERE id=$2", model.PasswordHash, id)
		} else {
			row := rep.db.QueryRowContext(ctx, "INSERT INTO users (username, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id", model.Username, model.PasswordHash, time.Now())
			err = row.Scan(&id)
			if helper.IsDuplicateError(err) {
				return ErrUserAlreadyExists
			}
		}

		return err
	})

	if err == nil {
		model.Id = &id
	}

	return model, err
}
