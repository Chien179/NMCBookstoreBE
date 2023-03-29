package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomSubgenre(t *testing.T, genre Genre) Subgenre {
	arg := CreateSubgenreParams{
		GenresID: genre.ID,
		Name:     util.RandomString(5),
	}

	subgenre, err := testQueries.CreateSubgenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, subgenre)

	require.Equal(t, arg.GenresID, subgenre.GenresID)
	require.Equal(t, arg.Name, subgenre.Name)

	require.NotZero(t, subgenre.ID)
	require.NotZero(t, subgenre.CreatedAt)

	return subgenre
}

func TestCreateSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	createRandomSubgenre(t, genre)
}

func TestGetSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	subgenre1 := createRandomSubgenre(t, genre)
	subgenre2, err := testQueries.GetSubgenre(context.Background(), subgenre1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, subgenre2)

	require.Equal(t, subgenre1.ID, subgenre2.ID)
	require.Equal(t, subgenre1.GenresID, subgenre2.GenresID)
	require.Equal(t, subgenre1.Name, subgenre2.Name)

	require.WithinDuration(t, subgenre1.CreatedAt, subgenre2.CreatedAt, time.Second)
}

func TestDeleteSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	subgenre1 := createRandomSubgenre(t, genre)

	err := testQueries.DeleteSubgenre(context.Background(), subgenre1.ID)
	require.NoError(t, err)

	subgenre2, err := testQueries.GetSubgenre(context.Background(), subgenre1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, subgenre2)
}

func TestUpdateSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	oldSubgenre := createRandomSubgenre(t, genre)

	newGenreID := util.RandomInt64(1, 2)
	newSubgenreName := util.RandomString(6)

	arg := UpdateSubgenreParams{
		ID: oldSubgenre.ID,
		GenresID: sql.NullInt64{
			Int64: newGenreID,
			Valid: true,
		},
		Name: sql.NullString{
			String: newSubgenreName,
			Valid:  true,
		},
	}

	updateSubgenre, err := testQueries.UpdateSubgenre(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateSubgenre)

	require.Equal(t, newGenreID, updateSubgenre.GenresID)
	require.NotEqual(t, oldSubgenre.GenresID, updateSubgenre.GenresID)
	require.Equal(t, newSubgenreName, updateSubgenre.Name)
	require.NotEqual(t, oldSubgenre.Name, updateSubgenre.Name)
}

func TestListSubgenres(t *testing.T) {
	genre := createRandomGenre(t)
	for i := 0; i < 10; i++ {
		createRandomSubgenre(t, genre)
	}

	subgenres, err := testQueries.ListSubgenres(context.Background(), genre.ID)

	require.NoError(t, err)
	require.NotEmpty(t, subgenres)

	for _, subgenre := range subgenres {
		require.NotEmpty(t, subgenre)
	}
}
