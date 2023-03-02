package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomWishlist(t *testing.T, user User) Wishlist {
	wishlists, err := testQueries.CreateWishlist(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, wishlists)

	require.Equal(t, user.ID, wishlists.UsersID)

	require.NotZero(t, wishlists.ID)
	require.NotZero(t, wishlists.CreatedAt)

	return wishlists
}

func TestCreateWishlist(t *testing.T) {
	user := createRandomUser(t)
	createRandomWishlist(t, user)
}

func TestGetWishlist(t *testing.T) {
	user := createRandomUser(t)
	wishlist1 := createRandomWishlist(t, user)
	wishlist2, err := testQueries.GetWishlist(context.Background(), wishlist1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, wishlist2)

	require.Equal(t, wishlist1.ID, wishlist1.ID)
	require.Equal(t, wishlist1.UsersID, wishlist2.UsersID)

	require.WithinDuration(t, wishlist1.CreatedAt, wishlist2.CreatedAt, time.Second)
}
