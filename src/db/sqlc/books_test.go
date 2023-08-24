package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/stretchr/testify/require"
)

func createRandomBook(t *testing.T) Book {
	arg := CreateBookParams{
		Name:        util.RandomUser(),
		Price:       util.RandomFloat(80000, 400000),
		Image:       []string{util.RandomString(10), util.RandomString(10), util.RandomString(10)},
		Description: util.RandomString(1000),
		Author:      util.RandomUser(),
		Publisher:   util.RandomUser(),
		Quantity:    util.RandomInt32(1, 50),
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
	oldBook := createRandomBook(t)

	newName := util.RandomString(6)
	newPrice := util.RandomFloat(60000, 1000000)
	newImages := []string{
		util.RandomString(20),
		util.RandomString(20),
		util.RandomString(20),
		util.RandomString(20),
		util.RandomString(20),
	}
	newDescription := util.RandomString(100)
	newAuthor := util.RandomString(6)
	newPublisher := util.RandomString(6)
	newQuantity := util.RandomInt32(1, 10)

	arg := UpdateBookParams{
		ID: oldBook.ID,
		Name: sql.NullString{
			String: newName,
			Valid:  true,
		},
		Price: sql.NullFloat64{
			Float64: newPrice,
			Valid:   true,
		},
		Image: newImages,
		Description: sql.NullString{
			String: newDescription,
			Valid:  true,
		},
		Author: sql.NullString{
			String: newAuthor,
			Valid:  true,
		},
		Publisher: sql.NullString{
			String: newPublisher,
			Valid:  true,
		},
		Quantity: sql.NullInt32{
			Int32: newQuantity,
			Valid: true,
		},
	}

	updateBook, err := testQueries.UpdateBook(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updateBook)

	require.Equal(t, newName, updateBook.Name)
	require.NotEqual(t, oldBook.Name, updateBook.Name)
	require.Equal(t, newPrice, updateBook.Price)
	require.NotEqual(t, oldBook.Price, updateBook.Price)
	require.Equal(t, newImages, updateBook.Image)
	require.NotEqual(t, oldBook.Image, updateBook.Image)
	require.Equal(t, newDescription, updateBook.Description)
	require.NotEqual(t, oldBook.Description, updateBook.Description)
	require.Equal(t, newAuthor, updateBook.Author)
	require.NotEqual(t, oldBook.Author, updateBook.Author)
	require.Equal(t, newPublisher, updateBook.Publisher)
	require.NotEqual(t, oldBook.Publisher, updateBook.Publisher)
	require.Equal(t, newQuantity, updateBook.Quantity)
	require.NotEqual(t, oldBook.Quantity, updateBook.Quantity)
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

	for _, book := range books.Books {
		require.NotEmpty(t, book)
	}
}
