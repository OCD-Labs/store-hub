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

	userID := "643a7dedbc8c7b338e50bd0f"

	userRole := util.RandomPermission()
	duration := 15 * time.Minute

	var perlvl int16 = 1

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(userID, userRole, &perlvl, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	fmt.Println(token)
	fmt.Printf("\n\n")

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, userRole, payload.UserRole)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// TestExpiredPasetoToken tests expiry of paseto tokens.
func TestExpiredPasetoToken(t *testing.T) {
	config, err := util.ParseConfigs("..")
	require.NoError(t, err)

	maker, err := NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	var perlvl int16 = 1

	token, payload, err := maker.CreateToken(
		"643a7dedbc8c7b338e50bd0f",
		util.RandomPermission(),
		&perlvl,
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