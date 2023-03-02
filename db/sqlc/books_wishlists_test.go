package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBookWishlist(t *testing.T, book Book, wishlist Wishlist) BooksWishlist {
	arg := CreateBookWishlistParams{
		BooksID:     book.ID,
		WishlistsID: wishlist.ID,
	}

	BookWishlist, err := testQueries.CreateBookWishlist(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, BookWishlist)

	require.Equal(t, arg.BooksID, BookWishlist.BooksID)
	require.Equal(t, arg.WishlistsID, BookWishlist.WishlistsID)

	require.NotZero(t, BookWishlist.ID)
	require.NotZero(t, BookWishlist.CreatedAt)

	return BookWishlist
}

func TestCreateBookWishlist(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	wishlist := createRandomWishlist(t, user)
	createRandomBookWishlist(t, book, wishlist)
}

func TestDeleteBookWishlist(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	wishlist := createRandomWishlist(t, user)
	bookWishlist1 := createRandomBookWishlist(t, book, wishlist)

	err := testQueries.DeleteBookWishlist(context.Background(), bookWishlist1.ID)
	require.NoError(t, err)

	bookWishlist2, err := testQueries.GetBookWishlist(context.Background(), bookWishlist1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, bookWishlist2)
}

func TestListBookWishlists(t *testing.T) {
	user := createRandomUser(t)
	wishlist := createRandomWishlist(t, user)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomBookWishlist(t, book, wishlist)
	}

	bookWishlists, err := testQueries.ListBooksWishlists(context.Background(), wishlist.ID)

	require.NoError(t, err)
	require.NotEmpty(t, bookWishlists)

	for _, bookWishlist := range bookWishlists {
		require.NotEmpty(t, bookWishlist)
	}
}
