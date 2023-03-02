package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBookGenre(t *testing.T, book Book, genre Genre) BooksGenre {
	arg := CreateBookGenreParams{
		BooksID:  book.ID,
		GenresID: genre.ID,
	}

	BookGenre, err := testQueries.CreateBookGenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, BookGenre)

	require.Equal(t, arg.BooksID, BookGenre.BooksID)
	require.Equal(t, arg.GenresID, BookGenre.GenresID)

	require.NotZero(t, BookGenre.ID)
	require.NotZero(t, BookGenre.CreatedAt)

	return BookGenre
}

func TestCreateBookGenre(t *testing.T) {
	book := createRandomBook(t)
	genre := createRandomGenre(t)
	createRandomBookGenre(t, book, genre)
}

func TestDeleteBookGenre(t *testing.T) {
	book := createRandomBook(t)
	genre := createRandomGenre(t)
	bookGenre1 := createRandomBookGenre(t, book, genre)

	err := testQueries.DeleteBookGenre(context.Background(), bookGenre1.ID)
	require.NoError(t, err)

	bookGenre2, err := testQueries.GetBookGenre(context.Background(), bookGenre1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, bookGenre2)
}

func TestListBookGenresByGenreID(t *testing.T) {
	genre := createRandomGenre(t)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomBookGenre(t, book, genre)
	}

	books, err := testQueries.ListBooksGenresByGenreID(context.Background(), genre.ID)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}

func TestListBookGenresByBookID(t *testing.T) {
	book := createRandomBook(t)
	for i := 0; i < 10; i++ {
		genre := createRandomGenre(t)
		createRandomBookGenre(t, book, genre)
	}

	genres, err := testQueries.ListBooksGenresByBookID(context.Background(), book.ID)

	require.NoError(t, err)
	require.NotEmpty(t, genres)

	for _, bookGenre := range genres {
		require.NotEmpty(t, bookGenre)
	}
}
