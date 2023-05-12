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
		Username: user.Username,
		Address:  util.RandomString(15),
		District: util.RandomString(6),
		City:     util.RandomString(6),
	}

	address, err := testQueries.CreateAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, address)

	require.Equal(t, arg.Username, address.Username)
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
	require.Equal(t, address1.Username, address2.Username)
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
	oldAddress := createRandomAddress(t, user)

	newAddress := util.RandomString(6)
	newDistrict := util.RandomString(6)
	newCity := util.RandomString(6)

	arg := UpdateAddressParams{
		ID: oldAddress.ID,
		Address: sql.NullString{
			String: newAddress,
			Valid:  true,
		},
		District: sql.NullString{
			String: newDistrict,
			Valid:  true,
		},
		City: sql.NullString{
			String: newCity,
			Valid:  true,
		},
	}

	updateAddress, err := testQueries.UpdateAddress(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateAddress)

	require.Equal(t, oldAddress.Username, updateAddress.Username)
	require.Equal(t, newAddress, updateAddress.Address)
	require.NotEqual(t, oldAddress.Address, updateAddress.Address)
	require.Equal(t, newDistrict, updateAddress.District)
	require.NotEqual(t, oldAddress.District, updateAddress.District)
	require.Equal(t, newCity, updateAddress.City)
	require.NotEqual(t, oldAddress.City, updateAddress.City)
}

func TestListAddresses(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomAddress(t, user)
	}

	arg := ListAddressesParams{
		Username: user.Username,
		Limit:    5,
		Offset:   0,
	}

	addresss, err := testQueries.ListAddresses(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, addresss)

	for _, address := range addresss.Address {
		require.NotEmpty(t, address)
	}
}
