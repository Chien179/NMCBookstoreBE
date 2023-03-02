package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomGenre(t *testing.T) Genre {
	arg := util.RandomString(6)

	Genre, err := testQueries.CreateGenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Genre)

	require.Equal(t, arg, Genre.Name)

	require.NotZero(t, Genre.ID)
	require.NotZero(t, Genre.CreatedAt)

	return Genre
}

func TestCreateGenre(t *testing.T) {
	createRandomGenre(t)
}

func TestGetGenre(t *testing.T) {
	genre1 := createRandomGenre(t)
	genre2, err := testQueries.GetGenre(context.Background(), genre1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, genre2)

	require.Equal(t, genre1.ID, genre2.ID)
	require.Equal(t, genre1.Name, genre2.Name)

	require.WithinDuration(t, genre1.CreatedAt, genre2.CreatedAt, time.Second)
}

func TestDeleteGenre(t *testing.T) {
	genre1 := createRandomGenre(t)

	err := testQueries.DeleteGenre(context.Background(), genre1.ID)
	require.NoError(t, err)

	genre2, err := testQueries.GetGenre(context.Background(), genre1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, genre2)
}

func TestUpdateGenre(t *testing.T) {
	genre1 := createRandomGenre(t)

	arg := UpdateGenreParams{
		ID:   genre1.ID,
		Name: genre1.Name,
	}

	genre2, err := testQueries.UpdateGenre(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, genre2)

	require.Equal(t, genre1.ID, genre2.ID)
	require.Equal(t, genre1.Name, genre2.Name)

	require.WithinDuration(t, genre1.CreatedAt, genre2.CreatedAt, time.Second)
}

func TestListGenres(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomGenre(t)
	}

	genres, err := testQueries.ListGenres(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, genres)

	for _, genre := range genres {
		require.NotEmpty(t, genre)
	}
}
