package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken             = errors.New("invalid token")
	ErrInvalidTokenSignedMethod = errors.New("invalid token method")
)

type claims struct {
	jwt.RegisteredClaims
	CustomData
}

// JwtTokenService структура для работы jwt токеном
type JwtTokenService struct {
	secretKey string
}

// CustomData сохраняемая структура пользовательских данных
type CustomData struct {
	UserId uint64
}

// NewJwtTokenService создает и возвращает сервис JwtTokenService с переданным secretKey
func NewJwtTokenService(secretKey string) *JwtTokenService {
	return &JwtTokenService{
		secretKey: secretKey,
	}
}

// GenerateToken генерирует jwt-token с переданным CustomData
func (service JwtTokenService) GenerateToken(customData *CustomData) (string, error) {
	customClaims := &claims{
		CustomData: *customData,
		ExpiresAt:  jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:   jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	signedToken, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseSignedToken распаршивает и валидирует signedToken и возращает CustomData
func (service JwtTokenService) ParseSignedToken(signedToken string) (*CustomData, error) {
	token, err := jwt.ParseWithClaims(signedToken, &claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenSignedMethod
		}
		return []byte(service.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if c, ok := token.Claims.(*claims); ok && token.Valid {
		return &c.CustomData, nil
	}

	return nil, ErrInvalidToken
}
