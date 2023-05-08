package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomWishlist(t *testing.T, book Book, user User) Wishlist {
	arg := CreateWishlistParams{
		BooksID:  book.ID,
		Username: user.Username,
	}

	Wishlist, err := testQueries.CreateWishlist(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Wishlist)

	require.Equal(t, arg.BooksID, Wishlist.BooksID)
	require.Equal(t, arg.Username, Wishlist.Username)

	require.NotZero(t, Wishlist.ID)
	require.NotZero(t, Wishlist.CreatedAt)

	return Wishlist
}

func TestCreateBookWishlist(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	createRandomWishlist(t, book, user)
}

func TestDeleteBookWishlist(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	Wishlist1 := createRandomWishlist(t, book, user)

	arg := DeleteWishlistParams{
		ID:       Wishlist1.ID,
		Username: user.Username,
	}

	err := testQueries.DeleteWishlist(context.Background(), arg)
	require.NoError(t, err)

	Wishlist2, err := testQueries.GetWishlist(context.Background(), Wishlist1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, Wishlist2)
}

func TestListBookWishlistsByUsername(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomWishlist(t, book, user)
	}

	books, err := testQueries.ListWishlistsByUsername(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}
