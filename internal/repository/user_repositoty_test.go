package repository

import (
	"database/sql"
	"gophemart/internal/helper"
	"gophemart/internal/model"
	"gophemart/migrations"
	"io"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openAndPrepareDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URI")
	dsn = "user=user password=password dbname=gophemart_test sslmode=disable host=localhost port=5433"
	if dsn == "" {
		t.Skip("DATABASE_URI environment variable not set")
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)

	require.NoError(t, db.PingContext(t.Context()))

	require.NoError(t, migrations.RunMigrations(db))

	_, err = db.Exec("TRUNCATE TABLE users")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func getRepo(t *testing.T, db *sql.DB) *UserRepository {
	logger := zerolog.New(io.Discard)

	return NewUserRepository(
		db,
		&logger,
		helper.NewRetryableHelper(
			&helper.StandardSleeper{},
			func(err error) bool {
				return false
			},
			[]time.Duration{},
			&logger,
		),
	)
}

func TestUserRepositorySave(t *testing.T) {
	db := openAndPrepareDB(t)
	repo := getRepo(t, db)

	savedUser, err := repo.Save(t.Context(), &model.User{
		Username:     "test",
		PasswordHash: "password",
		CreatedAt:    time.Time{},
	})
	require.NoError(t, err)
	require.NotNil(t, savedUser)
	assert.NotNil(t, savedUser.Id)

	var has bool
	res := db.QueryRowContext(t.Context(), "SELECT 1 FROM users WHERE id = $1", savedUser.Id)
	err = res.Scan(&has)
	require.NoError(t, err)
	assert.True(t, has)

	savedUser.PasswordHash = "password2"
	_, err = repo.Save(t.Context(), savedUser)

	res = db.QueryRowContext(t.Context(), "SELECT 1 FROM users WHERE id = $1 and password_hash=$2", savedUser.Id, savedUser.PasswordHash)
	err = res.Scan(&has)
	require.NoError(t, err)
	assert.True(t, has)

	user, err := repo.Save(t.Context(), &model.User{
		Username:     "test",
		PasswordHash: "password",
		CreatedAt:    time.Time{},
	})
	assert.ErrorIs(t, ErrUserAlreadyExists, err)
	assert.Nil(t, user.Id)
}

func TestUserRepositoryFind(t *testing.T) {
	db := openAndPrepareDB(t)
	repo := getRepo(t, db)

	_, err := repo.Save(t.Context(), &model.User{
		Username:     "test",
		PasswordHash: "password",
		CreatedAt:    time.Time{},
	})
	require.NoError(t, err)

	user, err := repo.FindByCredential(t.Context(), "test", "password")
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.NotEmpty(t, user.Id)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "password", user.PasswordHash)
	assert.NotEmpty(t, user.CreatedAt)

	user, err = repo.FindByCredential(t.Context(), "test", "password1")
	assert.ErrorIs(t, ErrUserNotFound, err)
}
