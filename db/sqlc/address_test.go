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
		Username:   user.Username,
		Address:    util.RandomString(15),
		DistrictID: util.RandomInt64(1, 10),
		CityID:     util.RandomInt64(1, 10),
	}

	address, err := testQueries.CreateAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, address)

	require.Equal(t, arg.Username, address.Username)
	require.Equal(t, arg.Address, address.Address)
	require.Equal(t, arg.DistrictID, address.DistrictID)
	require.Equal(t, arg.CityID, address.CityID)

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
	require.Equal(t, address1.DistrictID, address2.DistrictID)
	require.Equal(t, address1.CityID, address2.CityID)

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
	newDistrictID := util.RandomInt64(1, 10)
	newCityID := util.RandomInt64(1, 10)

	arg := UpdateAddressParams{
		ID: oldAddress.ID,
		Address: sql.NullString{
			String: newAddress,
			Valid:  true,
		},
		DistrictID: sql.NullInt64{
			Int64: newDistrictID,
			Valid: true,
		},
		CityID: sql.NullInt64{
			Int64: newCityID,
			Valid: true,
		},
	}

	updateAddress, err := testQueries.UpdateAddress(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateAddress)

	require.Equal(t, oldAddress.Username, updateAddress.Username)
	require.Equal(t, newAddress, updateAddress.Address)
	require.NotEqual(t, oldAddress.Address, updateAddress.Address)
	require.Equal(t, newDistrictID, updateAddress.DistrictID)
	require.NotEqual(t, oldAddress.DistrictID, updateAddress.DistrictID)
	require.Equal(t, newCityID, updateAddress.CityID)
	require.NotEqual(t, oldAddress.CityID, updateAddress.CityID)
}

func TestListAddresses(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomAddress(t, user)
	}

	addresss, err := testQueries.ListAddresses(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, addresss)

	for _, address := range addresss {
		require.NotEmpty(t, address)
	}
}
