package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// random secret at least 32 bytes long
	secret := "12345678901234567890123456789012"
	maker, err := NewPasetoMaker([]byte(secret))
	require.NoError(t, err)

	userID := "user-id"
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, userID, payload.UserID)
	require.NotZero(t, payload.PayloadID)
	require.Equal(t, issuedAt.Unix(), payload.IssuedAt.Unix(), time.Second)
	require.Equal(t, expiredAt.Unix(), payload.ExpiredAt.Unix(), time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	secret := "12345678901234567890123456789012"
	maker, err := NewPasetoMaker([]byte(secret))
	require.NoError(t, err)

	userID := "user-id"
	duration := -time.Minute

	token, err := maker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.ValidateToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
}
