// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

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
	DeleteAddress(ctx context.Context, id int64) error
	DeleteBook(ctx context.Context, id int64) error
	DeleteBookGenre(ctx context.Context, id int64) error
	DeleteBookSubgenre(ctx context.Context, id int64) error
	DeleteCart(ctx context.Context, arg DeleteCartParams) error
	DeleteGenre(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, id int64) error
	DeleteReview(ctx context.Context, id int64) error
	DeleteSubgenre(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, username string) error
	DeleteWishlist(ctx context.Context, arg DeleteWishlistParams) error
	FullSearch(ctx context.Context, arg FullSearchParams) ([]FullSearchRow, error)
	GetAddress(ctx context.Context, id int64) (Address, error)
	GetBook(ctx context.Context, id int64) (Book, error)
	GetBookGenre(ctx context.Context, id int64) (BooksGenre, error)
	GetBookSubgenre(ctx context.Context, id int64) (BooksSubgenre, error)
	GetCart(ctx context.Context, id int64) (Cart, error)
	GetGenre(ctx context.Context, id int64) (Genre, error)
	GetOrder(ctx context.Context, id int64) (Order, error)
	GetOrderToPayment(ctx context.Context, username string) (Order, error)
	GetReview(ctx context.Context, id int64) (Review, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetSubgenre(ctx context.Context, id int64) (Subgenre, error)
	GetTransaction(ctx context.Context, id int64) (Transaction, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetWishlist(ctx context.Context, id int64) (Wishlist, error)
	ListAddresses(ctx context.Context, arg ListAddressesParams) ([]Address, error)
	ListAllBooks(ctx context.Context) ([]Book, error)
	ListBooks(ctx context.Context, arg ListBooksParams) (ListBooksRow, error)
	ListBooksGenresByBookID(ctx context.Context, booksID int64) ([]BooksGenre, error)
	ListBooksGenresByGenreID(ctx context.Context, genresID int64) ([]BooksGenre, error)
	ListBooksSubgenresByBookID(ctx context.Context, booksID int64) ([]BooksSubgenre, error)
	ListBooksSubgenresBySubgenreID(ctx context.Context, subgenresID int64) ([]BooksSubgenre, error)
	ListCartsByUsername(ctx context.Context, username string) ([]Cart, error)
	ListGenres(ctx context.Context) ([]Genre, error)
	ListOders(ctx context.Context) ([]Order, error)
	ListOdersByUserName(ctx context.Context, arg ListOdersByUserNameParams) ([]Order, error)
	ListReviewsByBookID(ctx context.Context, arg ListReviewsByBookIDParams) ([]Review, error)
	ListSubgenres(ctx context.Context, genresID int64) ([]Subgenre, error)
	ListTop10NewestBooks(ctx context.Context) ([]Book, error)
	ListTop10TheBestBooks(ctx context.Context) ([]Book, error)
	ListTransactionsByOrderID(ctx context.Context, ordersID int64) ([]Transaction, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListWishlistsByUsername(ctx context.Context, username string) ([]Wishlist, error)
	RevenueDays(ctx context.Context, arg RevenueDaysParams) ([]RevenueDaysRow, error)
	RevenueHours(ctx context.Context, arg RevenueHoursParams) ([]RevenueHoursRow, error)
	RevenueMonths(ctx context.Context, arg RevenueMonthsParams) ([]RevenueMonthsRow, error)
	RevenueQuarters(ctx context.Context, arg RevenueQuartersParams) ([]RevenueQuartersRow, error)
	RevenueYears(ctx context.Context, arg RevenueYearsParams) ([]RevenueYearsRow, error)
	UpdateAddress(ctx context.Context, arg UpdateAddressParams) (Address, error)
	UpdateAmount(ctx context.Context, arg UpdateAmountParams) (Cart, error)
	UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error)
	UpdateGenre(ctx context.Context, arg UpdateGenreParams) (Genre, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateResetPassword(ctx context.Context, arg UpdateResetPasswordParams) (ResetPassword, error)
	UpdateSubgenre(ctx context.Context, arg UpdateSubgenreParams) (Subgenre, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)
