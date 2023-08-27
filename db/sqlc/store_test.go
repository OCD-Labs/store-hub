package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type StoreOwnerData struct {
	AccountID       string    `json:"account_id"`
	ProfileImgURL   string    `json:"profile_img_url"`
	AccessLevels    []int     `json:"access_levels"`
	IsOriginalOwner bool      `json:"original_owner"`
	AddedAt         time.Time `json:"added_at"`
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

	buf, err := store.StoreOwners.MarshalJSON()
	require.NoError(t, err)

	ownersData := []StoreOwnerData{}
	err = json.Unmarshal(buf, &ownersData)
	require.NoError(t, err)

	require.WithinDuration(t, res.StoreOwners[0].AddedAt, ownersData[0].AddedAt, time.Second)
	require.Equal(t, user.AccountID, ownersData[0].AccountID)
	require.Equal(t, user.ProfileImageUrl.String, ownersData[0].ProfileImgURL)

	// require.ElementsMatch(t, res.Owners[0].AccessLevels, ownersData[0].AccessLevels)
	// require.Equal(t, res.Owners[0].OriginalOwner, ownersData[0].IsOriginalOwner)
}
