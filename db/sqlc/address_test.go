package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomAddress(t *testing.T, user User) Address {
	arg := CreateAddressParams{
		UsersID:  user.ID,
		Address:  util.RandomString(15),
		District: util.RandomString(6),
		City:     util.RandomString(6),
	}

	address, err := testQueries.CreateAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, address)

	require.Equal(t, arg.UsersID, address.UsersID)
	require.Equal(t, arg.Address, address.Address)
	require.Equal(t, arg.District, address.District)
	require.Equal(t, arg.City, address.City)

	require.NotZero(t, address.ID)
	require.NotZero(t, address.CreatedAt)

	return address
}

func TestCreateAddress(t *testing.T) {
	user := createRandomUser(t)
	createRandomAddress(t, user)
}

func TestGetAddress(t *testing.T) {
	user := createRandomUser(t)
	address1 := createRandomAddress(t, user)
	address2, err := testQueries.GetAddress(context.Background(), address1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, address2)

	require.Equal(t, address1.ID, address2.ID)
	require.Equal(t, address1.UsersID, address2.UsersID)
	require.Equal(t, address1.Address, address2.Address)
	require.Equal(t, address1.District, address2.District)
	require.Equal(t, address1.City, address2.City)

	require.WithinDuration(t, address1.CreatedAt, address2.CreatedAt, time.Second)
}

func TestDeleteAddress(t *testing.T) {
	user := createRandomUser(t)
	address1 := createRandomAddress(t, user)

	err := testQueries.DeleteAddress(context.Background(), address1.ID)
	require.NoError(t, err)

	address2, err := testQueries.GetAddress(context.Background(), address1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, address2)
}

func TestUpdateAddress(t *testing.T) {
	user := createRandomUser(t)
	address1 := createRandomAddress(t, user)

	arg := UpdateAddressParams{
		ID:       address1.ID,
		Address:  address1.Address,
		District: address1.District,
		City:     address1.City,
	}

	address2, err := testQueries.UpdateAddress(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, address2)

	require.Equal(t, address1.ID, address2.ID)
	require.Equal(t, address1.UsersID, address2.UsersID)
	require.Equal(t, address1.Address, address2.Address)
	require.Equal(t, address1.District, address2.District)
	require.Equal(t, address1.City, address2.City)

	require.WithinDuration(t, address1.CreatedAt, address2.CreatedAt, time.Second)
}

func TestListAddresses(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomAddress(t, user)
	}

	arg := ListAddressesParams{
		Limit:  5,
		Offset: 0,
	}

	addresss, err := testQueries.ListAddresses(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, addresss)

	for _, address := range addresss {
		require.NotEmpty(t, address)
	}
}
