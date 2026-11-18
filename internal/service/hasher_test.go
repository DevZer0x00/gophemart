package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasherServiceHash(t *testing.T) {
	hasher := NewHasherService("test")
	hash, err := hasher.Hash([]byte("string"))
	require.NoError(t, err)
	assert.Len(t, hash, 64)
}
