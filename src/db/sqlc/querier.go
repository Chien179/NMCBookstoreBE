// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAddress(ctx context.Context, arg CreateAddressParams) (Address, error)
	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	CreateBookGenre(ctx context.Context, arg CreateBookGenreParams) (BooksGenre, error)
	CreateBookSubgenre(ctx context.Context, arg CreateBookSubgenreParams) (BooksSubgenre, error)
	CreateCart(ctx context.Context, arg CreateCartParams) (Cart, error)
	CreateGenre(ctx context.Context, name string) (Genre, error)
	CreateLike(ctx context.Context, arg CreateLikeParams) (Like, error)
	CreateOrder(ctx context.Context, username string) (Order, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error)
	CreateResetPassword(ctx context.Context, arg CreateResetPasswordParams) (ResetPassword, error)
	CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateShipping(ctx context.Context, arg CreateShippingParams) (Shipping, error)
	CreateSubgenre(ctx context.Context, arg CreateSubgenreParams) (Subgenre, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	CreateWishlist(ctx context.Context, arg CreateWishlistParams) (Wishlist, error)
	CreatedDislike(ctx context.Context, arg CreatedDislikeParams) (Dislike, error)
	DeleteAddress(ctx context.Context, id int64) error
	DeleteBook(ctx context.Context, id int64) error
	DeleteBookGenreByBooksID(ctx context.Context, booksID int64) error
	DeleteBookSubgenreByBooksID(ctx context.Context, booksID int64) error
	DeleteCart(ctx context.Context, arg DeleteCartParams) error
	DeleteGenre(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, id int64) error
	DeleteReview(ctx context.Context, id int64) error
	DeleteSubgenre(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, username string) error
	DeleteWishlist(ctx context.Context, arg DeleteWishlistParams) error
	GetAddress(ctx context.Context, id int64) (Address, error)
	GetBestBookByUser(ctx context.Context, username string) (GetBestBookByUserRow, error)
	GetBook(ctx context.Context, id int64) (Book, error)
	GetBookGenre(ctx context.Context, id int64) (BooksGenre, error)
	GetBookSubgenre(ctx context.Context, id int64) (BooksSubgenre, error)
	GetCart(ctx context.Context, id int64) (Cart, error)
	GetCity(ctx context.Context, id int64) (City, error)
	GetCountLikeByUser(ctx context.Context, username string) (int64, error)
	GetCountReviewByUser(ctx context.Context, username string) (int64, error)
	GetDislike(ctx context.Context, arg GetDislikeParams) (Dislike, error)
	GetDistrict(ctx context.Context, id int64) (District, error)
	GetGenre(ctx context.Context, id int64) (Genre, error)
	GetLike(ctx context.Context, arg GetLikeParams) (Like, error)
	GetOrder(ctx context.Context, id int64) (Order, error)
	GetOrderToPayment(ctx context.Context, username string) (Order, error)
	GetRank(ctx context.Context, score int32) (GetRankRow, error)
	GetReview(ctx context.Context, id int64) (Review, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetSubgenre(ctx context.Context, id int64) (Subgenre, error)
	GetTransaction(ctx context.Context, id int64) (Transaction, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetWishlist(ctx context.Context, id int64) (Wishlist, error)
	ListAddresses(ctx context.Context, username string) ([]ListAddressesRow, error)
	ListAllBooks(ctx context.Context) ([]Book, error)
	ListAllOders(ctx context.Context) ([]Order, error)
	ListAllSubgenres(ctx context.Context) ([]Subgenre, error)
	ListBookFollowGenre(ctx context.Context, arg ListBookFollowGenreParams) ([]ListBookFollowGenreRow, error)
	ListBooks(ctx context.Context, arg ListBooksParams) (ListBooksRow, error)
	ListBooksGenresByBookID(ctx context.Context, booksID int64) ([]BooksGenre, error)
	ListBooksGenresByGenreID(ctx context.Context, genresID int64) ([]BooksGenre, error)
	ListBooksGenresIDByBookID(ctx context.Context, booksID int64) ([]int64, error)
	ListBooksSubgenresByBookID(ctx context.Context, booksID int64) ([]BooksSubgenre, error)
	ListBooksSubgenresBySubgenreID(ctx context.Context, subgenresID int64) ([]BooksSubgenre, error)
	ListBooksSubgenresIDByBookID(ctx context.Context, booksID int64) ([]int64, error)
	ListCartsByUsername(ctx context.Context, username string) ([]Cart, error)
	ListCities(ctx context.Context) ([]City, error)
	ListDislike(ctx context.Context, username string) ([]Dislike, error)
	ListDistricts(ctx context.Context, cityID int64) ([]District, error)
	ListGenres(ctx context.Context) ([]Genre, error)
	ListLike(ctx context.Context, username string) ([]Like, error)
	ListNewestBooks(ctx context.Context) ([]Book, error)
	ListOders(ctx context.Context, arg ListOdersParams) (ListOdersRow, error)
	ListOdersByUserName(ctx context.Context, username string) ([]Order, error)
	ListReviews(ctx context.Context) ([]Review, error)
	ListReviewsByBookID(ctx context.Context, arg ListReviewsByBookIDParams) (ListReviewsByBookIDRow, error)
	ListSubgenres(ctx context.Context, genresID int64) ([]Subgenre, error)
	ListSubgenresNoticeable(ctx context.Context) ([]ListSubgenresNoticeableRow, error)
	ListTheBestBooks(ctx context.Context) ([]Book, error)
	ListTransactionsByOrderID(ctx context.Context, ordersID int64) ([]Transaction, error)
	ListUsers(ctx context.Context) ([]User, error)
	ListWishlistsByUsername(ctx context.Context, username string) ([]Wishlist, error)
	RevenueDays(ctx context.Context) ([]RevenueDaysRow, error)
	RevenueMonths(ctx context.Context) ([]RevenueMonthsRow, error)
	RevenueQuarters(ctx context.Context) ([]RevenueQuartersRow, error)
	RevenueYears(ctx context.Context) ([]RevenueYearsRow, error)
	SoftDeleteBook(ctx context.Context, id int64) (Book, error)
	SoftDeleteGenre(ctx context.Context, id int64) (Genre, error)
	SoftDeleteReview(ctx context.Context, id int64) (Review, error)
	SoftDeleteSubgenre(ctx context.Context, id int64) (Subgenre, error)
	SoftDeleteUser(ctx context.Context, id string) (User, error)
	UpdateAddress(ctx context.Context, arg UpdateAddressParams) (Address, error)
	UpdateAmount(ctx context.Context, arg UpdateAmountParams) (Cart, error)
	UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error)
	UpdateDislike(ctx context.Context, arg UpdateDislikeParams) (Dislike, error)
	UpdateGenre(ctx context.Context, arg UpdateGenreParams) (Genre, error)
	UpdateLike(ctx context.Context, arg UpdateLikeParams) (Like, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateResetPassword(ctx context.Context, arg UpdateResetPasswordParams) (ResetPassword, error)
	UpdateReview(ctx context.Context, arg UpdateReviewParams) (Review, error)
	UpdateSubgenre(ctx context.Context, arg UpdateSubgenreParams) (Subgenre, error)
	UpdateTransaction(ctx context.Context, id int64) (Transaction, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)
