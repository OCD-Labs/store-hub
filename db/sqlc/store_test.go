package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Data struct {
	StoreOwners struct {
		AddedAt time.Time `json:"added_at"`
		StoreID int64     `json:"store_id"`
		UserID  int64     `json:"user_id"`
	} `json:"store_owners"`
	User struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		ID        int64  `json:"id"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

func TestGetStoreByID(t *testing.T) {
	res, user := createStoreAndOwners(t)
	store, err := testQueries.GetStoreByID(context.Background(), res.Store.ID)
	require.NoError(t, err)
	require.NotEmpty(t, store)

	require.Equal(t, res.Store.ID, store.ID)
	require.Equal(t, res.Store.Name, store.Name)
	require.Equal(t, res.Store.Description, store.Description)
	require.Equal(t, res.Store.ProfileImageUrl, store.ProfileImageUrl)
	require.Equal(t, res.Store.IsVerified, store.IsVerified)
	require.Equal(t, res.Store.Category, store.Category)
	require.WithinDuration(t, res.Store.CreatedAt, store.CreatedAt, time.Second)

	buf, err := store.Owners.MarshalJSON()
	require.NoError(t, err)

	data := []Data{}
	err = json.Unmarshal(buf, &data)
	require.NoError(t, err)

	require.WithinDuration(t, res.Owners[0].AddedAt, data[0].StoreOwners.AddedAt, time.Second)
	require.Equal(t, res.Owners[0].StoreID, data[0].StoreOwners.StoreID)
	require.Equal(t, res.Owners[0].UserID, data[0].StoreOwners.UserID)
	require.Equal(t, user.Email, data[0].User.Email)
	require.Equal(t, user.FirstName, data[0].User.FirstName)
	require.Equal(t, user.LastName, data[0].User.LastName)
	require.Equal(t, user.ID, data[0].User.ID)
}
