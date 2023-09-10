package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserTx(t *testing.T) {
	hashedPassword, err := util.HashedPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserTxParams{
		CreateUserParams: CreateUserParams{
			FirstName:      util.RandomOwner(),
			LastName:       util.RandomOwner(),
			Status:         util.RandomPermission(),
			HashedPassword: hashedPassword,
			Email:          util.RandomEmail(),
			ProfileImageUrl: sql.NullString{
				String: "",
				Valid:  false,
			},
			About:     "",
			Socials:   json.RawMessage([]byte("{}")),
			AccountID: fmt.Sprintf("%s.testnet", util.RandomOwner()),
		},
		AfterCreate: func(user User) error {
			return nil
		},
	}

	user, err := testQueries.CreateUserTx(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Status, user.Status)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.ProfileImageUrl.String, user.ProfileImageUrl.String)
	require.Equal(t, arg.ProfileImageUrl.Valid, user.ProfileImageUrl.Valid)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	require.True(t, user.IsActive)
	require.False(t, user.IsEmailVerified)
}
