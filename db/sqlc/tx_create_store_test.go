package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/stretchr/testify/require"
)

func createStoreAndOwners(t *testing.T) (CreateStoreTxResult, User) {
	user := createRandomUser(t)
	arg := CreateStoreTxParams{
		CreateStoreParams: CreateStoreParams{
			Name: util.RandomString(8),
			Description: util.RandomString(20),
			ProfileImageUrl: fmt.Sprintf("https://%s.com", util.RandomString(15)),
			Category: util.RandomString(5),
		},
		OwnerID: user.ID,
	}

	res, err := testQueries.CreateStoreTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, arg.Name, res.Store.Name)
	require.Equal(t, arg.Description, res.Store.Description)
	require.Equal(t, arg.ProfileImageUrl, res.Store.ProfileImageUrl)
	require.Equal(t, arg.Category, res.Store.Category)

	require.False(t, res.Store.IsVerified)
	require.NotZero(t, res.Store.CreatedAt)

	require.Equal(t, arg.OwnerID, res.Owners[0].UserID)
	require.Equal(t, res.Store.ID, res.Owners[0].StoreID)
	require.NotZero(t, res.Owners[0].AddedAt)

	return res, user
}

func TestCreateStoreAndOwners(t *testing.T) {
	createStoreAndOwners(t)
}