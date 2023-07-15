package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/stretchr/testify/require"
)

// TestPasetoMaker tests creating of paseto tokens.
func TestPasetoMaker(t *testing.T) {
	config, err := util.ParseConfigs("..")
	require.NoError(t, err)

	maker, err := NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	var userID int64 = 1
	accountID := "643a7dedbc8c7b338e50bd0f"

	duration := 60 * time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(userID, accountID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	fmt.Println(token)
	fmt.Println(userID)
	fmt.Println(accountID)
	fmt.Printf("\n\n")

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, accountID, payload.AccountID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// TestExpiredPasetoToken tests expiry of paseto tokens.
func TestExpiredPasetoToken(t *testing.T) {
	config, err := util.ParseConfigs("..")
	require.NoError(t, err)

	maker, err := NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(
		1,
		"643a7dedbc8c7b338e50bd0f",
		-time.Minute,
	)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
