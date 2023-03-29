package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:    util.RandomUser(),
		FullName:    util.RandomUser(),
		Email:       util.RandomEmail(),
		Password:    hashedPassword,
		Image:       util.RandomString(10),
		PhoneNumber: util.RandomString(10),
		Role:        "user",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Image, user.Image)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.Role, user.Role)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Image, user2.Image)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.Role, user2.Role)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomUser()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newImages := util.RandomString(10)
	newPhoneNumber := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Image: sql.NullString{
			String: newImages,
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: newPhoneNumber,
			Valid:  true,
		},
		Password: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	}

	updateUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)

	require.Equal(t, newEmail, updateUser.Email)
	require.NotEqual(t, oldUser.Email, updateUser.Email)
	require.Equal(t, newHashedPassword, updateUser.Password)
	require.NotEqual(t, oldUser.Password, updateUser.Password)
	require.Equal(t, newImages, updateUser.Image)
	require.NotEqual(t, oldUser.Image, updateUser.Image)
	require.Equal(t, newPhoneNumber, updateUser.PhoneNumber)
	require.NotEqual(t, oldUser.PhoneNumber, updateUser.PhoneNumber)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
