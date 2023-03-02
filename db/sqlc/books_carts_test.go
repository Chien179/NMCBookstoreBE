package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBookCart(t *testing.T, book Book, cart Cart) BooksCart {
	arg := CreateBookCartParams{
		BooksID: book.ID,
		CartsID: cart.ID,
	}

	BookCart, err := testQueries.CreateBookCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, BookCart)

	require.Equal(t, arg.BooksID, BookCart.BooksID)
	require.Equal(t, arg.CartsID, BookCart.CartsID)

	require.NotZero(t, BookCart.ID)
	require.NotZero(t, BookCart.CreatedAt)

	return BookCart
}

func TestCreateBookCart(t *testing.T) {
	user := createRandomUser(t)
	cart := createRandomCart(t, user)
	book := createRandomBook(t)
	createRandomBookCart(t, book, cart)
}

func TestDeleteBookCart(t *testing.T) {
	user := createRandomUser(t)
	cart := createRandomCart(t, user)
	book := createRandomBook(t)
	bookCart1 := createRandomBookCart(t, book, cart)

	err := testQueries.DeleteBookCart(context.Background(), bookCart1.ID)
	require.NoError(t, err)

	bookCart2, err := testQueries.GetBookCart(context.Background(), bookCart1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, bookCart2)
}

func TestListBookCartsByCartID(t *testing.T) {
	user := createRandomUser(t)
	cart := createRandomCart(t, user)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomBookCart(t, book, cart)
	}

	books, err := testQueries.ListBooksCartsByCartID(context.Background(), cart.ID)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}

func TestListBookCartsByBookID(t *testing.T) {
	book := createRandomBook(t)
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		cart := createRandomCart(t, user)
		createRandomBookCart(t, book, cart)
	}

	carts, err := testQueries.ListBooksCartsByBookID(context.Background(), book.ID)

	require.NoError(t, err)
	require.NotEmpty(t, carts)

	for _, cart := range carts {
		require.NotEmpty(t, cart)
	}
}
