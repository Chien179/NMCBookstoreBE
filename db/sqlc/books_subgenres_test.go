package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBookSubgenre(t *testing.T, book Book, subgenre Subgenre) BooksSubgenre {
	arg := CreateBookSubgenreParams{
		BooksID:     book.ID,
		SubgenresID: subgenre.ID,
	}

	BookSubgenre, err := testQueries.CreateBookSubgenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, BookSubgenre)

	require.Equal(t, arg.BooksID, BookSubgenre.BooksID)
	require.Equal(t, arg.SubgenresID, BookSubgenre.SubgenresID)

	require.NotZero(t, BookSubgenre.ID)
	require.NotZero(t, BookSubgenre.CreatedAt)

	return BookSubgenre
}

func TestCreateBookSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	book := createRandomBook(t)
	subgenre := createRandomSubgenre(t, genre)
	createRandomBookSubgenre(t, book, subgenre)
}

func TestDeleteBookSubgenre(t *testing.T) {
	genre := createRandomGenre(t)
	book := createRandomBook(t)
	subgenre := createRandomSubgenre(t, genre)
	bookSubgenre1 := createRandomBookSubgenre(t, book, subgenre)

	err := testQueries.DeleteBookSubgenre(context.Background(), bookSubgenre1.ID)
	require.NoError(t, err)

	bookSubgenre2, err := testQueries.GetBookSubgenre(context.Background(), bookSubgenre1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, bookSubgenre2)
}

func TestListBookSubgenresBySubgenreID(t *testing.T) {
	genre := createRandomGenre(t)
	subgenre := createRandomSubgenre(t, genre)
	for i := 0; i < 10; i++ {
		book := createRandomBook(t)
		createRandomBookSubgenre(t, book, subgenre)
	}

	books, err := testQueries.ListBooksSubgenresBySubgenreID(context.Background(), subgenre.ID)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}

func TestListBookSubgenresByBookID(t *testing.T) {
	book := createRandomBook(t)
	for i := 0; i < 10; i++ {
		genre := createRandomGenre(t)
		subgenre := createRandomSubgenre(t, genre)
		createRandomBookSubgenre(t, book, subgenre)
	}

	subgenres, err := testQueries.ListBooksSubgenresByBookID(context.Background(), book.ID)

	require.NoError(t, err)
	require.NotEmpty(t, subgenres)

	for _, subgenre := range subgenres {
		require.NotEmpty(t, subgenre)
	}
}
