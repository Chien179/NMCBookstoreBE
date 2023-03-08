package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCart(t *testing.T, book Book, user User) Cart {
	arg := CreateCartParams{
		BooksID:  book.ID,
		Username: user.Username,
	}

	Cart, err := testQueries.CreateCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Cart)

	require.Equal(t, arg.BooksID, Cart.BooksID)
	require.Equal(t, arg.Username, Cart.Username)

	require.NotZero(t, Cart.ID)
	require.NotZero(t, Cart.CreatedAt)

	return Cart
}

func TestCreateCart(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	createRandomCart(t, book, user)
}

func TestDeleteCart(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	Cart1 := createRandomCart(t, book, user)

	err := testQueries.DeleteCart(context.Background(), Cart1.ID)
	require.NoError(t, err)

	Cart2, err := testQueries.GetCart(context.Background(), Cart1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, Cart2)
}

func TestListBookCartsByUsername(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomCart(t, book, user)
	}

	books, err := testQueries.ListCartsByUsername(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}
