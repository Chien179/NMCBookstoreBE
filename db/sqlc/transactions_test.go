package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransaction(t *testing.T, order Order, book Book) Transaction {
	arg := CreateTransactionParams{
		OrdersID: order.ID,
		BooksID:  book.ID,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.OrdersID, transaction.OrdersID)
	require.Equal(t, arg.BooksID, transaction.BooksID)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	user := createRandomUser(t)
	order := createRandomOrder(t, user)
	book := createRandomBook(t)
	createRandomTransaction(t, order, book)
}

func TestGetTransaction(t *testing.T) {
	user := createRandomUser(t)
	order := createRandomOrder(t, user)
	book := createRandomBook(t)
	transaction1 := createRandomTransaction(t, order, book)
	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	require.Equal(t, transaction1.ID, transaction2.ID)
	require.Equal(t, transaction1.OrdersID, transaction2.OrdersID)
	require.Equal(t, transaction1.BooksID, transaction2.BooksID)

	require.WithinDuration(t, transaction1.CreatedAt, transaction2.CreatedAt, time.Second)
}

func TestDeleteTransaction(t *testing.T) {
	user := createRandomUser(t)
	order := createRandomOrder(t, user)
	book := createRandomBook(t)
	transaction1 := createRandomTransaction(t, order, book)

	err := testQueries.DeleteTransaction(context.Background(), transaction1.ID)
	require.NoError(t, err)

	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transaction2)
}

func TestListTransactionsByOrderID(t *testing.T) {
	user := createRandomUser(t)
	order := createRandomOrder(t, user)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomTransaction(t, order, book)
	}

	transactions, err := testQueries.ListTransactionsByOrderID(context.Background(), order.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}
