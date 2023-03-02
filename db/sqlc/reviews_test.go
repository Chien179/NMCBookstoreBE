package db

import (
	"context"
	"testing"
	"time"

	"github.com/Chien179/NMCBookstoreBE/util"
	"github.com/stretchr/testify/require"
)

func createRandomReview(t *testing.T, user User, book Book) Review {
	arg := CreateReviewParams{
		UsersID:  user.ID,
		BooksID:  book.ID,
		Comments: util.RandomString(100),
		Rating:   util.RandomInt(0, 100),
	}

	review, err := testQueries.CreateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, review)

	require.Equal(t, arg.UsersID, review.UsersID)
	require.Equal(t, arg.BooksID, review.BooksID)
	require.Equal(t, arg.Comments, review.Comments)
	require.Equal(t, arg.Rating, review.Rating)

	require.NotZero(t, review.ID)
	require.NotZero(t, review.CreatedAt)

	return review
}

func TestCreateReview(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	createRandomReview(t, user, book)
}

func TestGetReview(t *testing.T) {
	user := createRandomUser(t)
	book := createRandomBook(t)
	review1 := createRandomReview(t, user, book)
	review2, err := testQueries.GetReview(context.Background(), review1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, review2)

	require.Equal(t, review1.ID, review1.ID)
	require.Equal(t, review1.UsersID, review2.UsersID)
	require.Equal(t, review1.BooksID, review2.BooksID)
	require.Equal(t, review1.Comments, review2.Comments)
	require.Equal(t, review1.Comments, review2.Comments)

	require.WithinDuration(t, review1.CreatedAt, review2.CreatedAt, time.Second)
}

func TestGetReviewByBookId(t *testing.T) {
	book := createRandomBook(t)
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		createRandomReview(t, user, book)
	}

	arg := GetReviewsByBookIDParams{
		BooksID: book.ID,
		Limit:   5,
		Offset:  0,
	}

	reviews, err := testQueries.GetReviewsByBookID(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, reviews)

	for _, book := range reviews {
		require.NotEmpty(t, book)
	}
}
