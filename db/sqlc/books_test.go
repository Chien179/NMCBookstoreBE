package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomBook(t *testing.T) Book {
	arg := CreateBookParams{
		Name:        util.RandomUser(),
		Price:       util.RandomFloat(80000, 400000),
		Image:       util.RandomString(10),
		Description: util.RandomString(1000),
		Author:      util.RandomUser(),
		Publisher:   util.RandomUser(),
		Quantity:    util.RandomInt(1, 50),
	}

	book, err := testQueries.CreateBook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, book)

	require.Equal(t, arg.Name, book.Name)
	require.Equal(t, arg.Price, book.Price)
	require.Equal(t, arg.Image, book.Image)
	require.Equal(t, arg.Description, book.Description)
	require.Equal(t, arg.Author, book.Author)
	require.Equal(t, arg.Publisher, book.Publisher)
	require.Equal(t, arg.Quantity, book.Quantity)

	require.NotZero(t, book.ID)
	require.NotZero(t, book.CreatedAt)

	return book
}

func TestCreateBook(t *testing.T) {
	createRandomBook(t)
}

func TestGetBook(t *testing.T) {
	book1 := createRandomBook(t)
	book2, err := testQueries.GetBook(context.Background(), book1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, book2)

	require.Equal(t, book1.ID, book2.ID)
	require.Equal(t, book1.Name, book2.Name)
	require.Equal(t, book1.Price, book2.Price)
	require.Equal(t, book1.Image, book2.Image)
	require.Equal(t, book1.Description, book2.Description)
	require.Equal(t, book1.Author, book2.Author)
	require.Equal(t, book1.Publisher, book2.Publisher)
	require.Equal(t, book1.Quantity, book2.Quantity)

	require.WithinDuration(t, book1.CreatedAt, book2.CreatedAt, time.Second)
}

func TestDeleteBook(t *testing.T) {
	book1 := createRandomBook(t)

	err := testQueries.DeleteBook(context.Background(), book1.ID)
	require.NoError(t, err)

	book2, err := testQueries.GetBook(context.Background(), book1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, book2)
}

func TestUpdateBook(t *testing.T) {
	book1 := createRandomBook(t)

	arg := UpdateBookParams{
		ID:          book1.ID,
		Name:        book1.Name,
		Price:       book1.Price,
		Image:       book1.Image,
		Description: book1.Description,
		Author:      book1.Author,
		Publisher:   book1.Publisher,
		Quantity:    book1.Quantity,
	}

	book2, err := testQueries.UpdateBook(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, book2)

	require.Equal(t, book1.ID, book2.ID)
	require.Equal(t, book1.Name, book2.Name)
	require.Equal(t, book1.Price, book2.Price)
	require.Equal(t, book1.Image, book2.Image)
	require.Equal(t, book1.Description, book2.Description)
	require.Equal(t, book1.Author, book2.Author)
	require.Equal(t, book1.Publisher, book2.Publisher)
	require.Equal(t, book1.Quantity, book2.Quantity)

	require.WithinDuration(t, book1.CreatedAt, book2.CreatedAt, time.Second)
}

func TestListBooks(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomBook(t)
	}

	arg := ListBooksParams{
		Limit:  5,
		Offset: 0,
	}

	books, err := testQueries.ListBooks(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, books)

	for _, book := range books {
		require.NotEmpty(t, book)
	}
}
