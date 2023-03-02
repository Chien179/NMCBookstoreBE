package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomCart(t *testing.T, user User) Cart {
	carts, err := testQueries.CreateCart(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, carts)

	require.Equal(t, user.ID, carts.UsersID)

	require.NotZero(t, carts.ID)
	require.NotZero(t, carts.CreatedAt)

	return carts
}

func TestCreateCart(t *testing.T) {
	user := createRandomUser(t)
	createRandomCart(t, user)
}

func TestGetCart(t *testing.T) {
	user := createRandomUser(t)
	cart1 := createRandomCart(t, user)
	cart2, err := testQueries.GetCart(context.Background(), cart1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cart2)

	require.Equal(t, cart1.ID, cart1.ID)
	require.Equal(t, cart1.UsersID, cart2.UsersID)

	require.WithinDuration(t, cart1.CreatedAt, cart2.CreatedAt, time.Second)
}
