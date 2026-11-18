package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtTokenGenerateToken(t *testing.T) {
	service := NewJwtTokenService("test")
	generatedToken, err := service.GenerateToken(&CustomData{
		UserId: uint64(100),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, generatedToken)
}

func TestJwtTokenParseToken(t *testing.T) {
	service := NewJwtTokenService("test")
	generatedToken, err := service.GenerateToken(&CustomData{
		UserId: uint64(100),
	})
	require.NoError(t, err)

	data, err := service.ParseSignedToken(generatedToken)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), data.UserId)
}
