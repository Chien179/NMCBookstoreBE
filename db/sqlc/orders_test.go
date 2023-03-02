package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T, user User) Order {
	order, err := testQueries.CreateOrder(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.Equal(t, user.ID, order.UsersID)

	require.NotZero(t, order.ID)
	require.NotZero(t, order.CreatedAt)

	return order
}

func TestCreateOrder(t *testing.T) {
	user := createRandomUser(t)
	createRandomOrder(t, user)
}

func TestGetOrder(t *testing.T) {
	user := createRandomUser(t)
	order1 := createRandomOrder(t, user)
	order2, err := testQueries.GetOrder(context.Background(), order1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, order1.UsersID, order2.UsersID)

	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
}

func TestDeleteOrder(t *testing.T) {
	user := createRandomUser(t)
	order1 := createRandomOrder(t, user)

	err := testQueries.DeleteOrder(context.Background(), order1.ID)
	require.NoError(t, err)

	order2, err := testQueries.GetOrder(context.Background(), order1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, order2)
}

func TestListOrdersByUserID(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomOrder(t, user)
	}

	arg := ListOdersByUserIDParams{
		UsersID: user.ID,
		Limit:   5,
		Offset:  0,
	}

	orders, err := testQueries.ListOdersByUserID(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orders)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}
