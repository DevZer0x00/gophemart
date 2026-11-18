package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HasherService структура для хэширования данных с секретным ключом
type HasherService struct {
	secretKey string
}

// NewHasherService создает и возвращает HasherService с переданным секретным ключом
func NewHasherService(secretKey string) *HasherService {
	return &HasherService{secretKey: secretKey}
}

// Hash возвращает хэш переданных данных
func (s *HasherService) Hash(data []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(s.secretKey))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
