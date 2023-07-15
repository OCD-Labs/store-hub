package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashedPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
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
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
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

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func compareUsers(t *testing.T, user1 User, user2 User) {
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.ProfileImageUrl.String, user2.ProfileImageUrl.String)
	require.Equal(t, user1.ProfileImageUrl.Valid, user2.ProfileImageUrl.Valid)
	require.Equal(t, user1.IsActive, user2.IsActive)
	require.Equal(t, user1.IsEmailVerified, user2.IsEmailVerified)
	require.Equal(t, user1.About, user2.About)
	require.Equal(t, user1.Socials, user2.Socials)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	compareUsers(t, user1, user2)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	compareUsers(t, user1, user2)
}

func TestUpdateUserFullNameOnly(t *testing.T) {
	oldUser := createRandomUser(t)

	newFirstName := util.RandomOwner()
	newLastName := util.RandomOwner()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: sql.NullInt64{
			Int64: oldUser.ID,
			Valid: true,
		},
		FirstName: sql.NullString{
			String: newFirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: newLastName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FirstName, updatedUser.FirstName)
	require.NotEqual(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, newFirstName, updatedUser.FirstName)
	require.Equal(t, newLastName, updatedUser.LastName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserEmailOnly(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: sql.NullInt64{
			Int64: oldUser.ID,
			Valid: true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserPasswordOnly(t *testing.T) {
	oldUser := createRandomUser(t)

	newHashedPassword, err := util.HashedPassword(util.RandomString(10))
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: sql.NullInt64{
			Int64: oldUser.ID,
			Valid: true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}
