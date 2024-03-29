// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID         int64     `json:"id"`
	Address    string    `json:"address"`
	Username   string    `json:"username"`
	CityID     int64     `json:"city_id"`
	DistrictID int64     `json:"district_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Book struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Image       []string  `json:"image"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Sale        float64   `json:"sale"`
	Quantity    int32     `json:"quantity"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	Rating      float64   `json:"rating"`
}

type BooksGenre struct {
	ID        int64     `json:"id"`
	BooksID   int64     `json:"books_id"`
	GenresID  int64     `json:"genres_id"`
	CreatedAt time.Time `json:"created_at"`
}

type BooksSubgenre struct {
	ID          int64     `json:"id"`
	BooksID     int64     `json:"books_id"`
	SubgenresID int64     `json:"subgenres_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Cart struct {
	ID        int64     `json:"id"`
	BooksID   int64     `json:"books_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Amount    int32     `json:"amount"`
	Total     float64   `json:"total"`
}

type City struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Dislike struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	ReviewID  int64  `json:"review_id"`
	IsDislike bool   `json:"is_dislike"`
}

type District struct {
	ID        int64     `json:"id"`
	CityID    int64     `json:"city_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Genre struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
}

type Like struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	ReviewID int64  `json:"review_id"`
	IsLike   bool   `json:"is_like"`
}

type Order struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	CreatedAt time.Time      `json:"created_at"`
	Status    string         `json:"status"`
	SubAmount int32          `json:"sub_amount"`
	SubTotal  float64        `json:"sub_total"`
	Sale      float64        `json:"sale"`
	Note      sql.NullString `json:"note"`
}

type Payment struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	OrderID    int64     `json:"order_id"`
	ShippingID int64     `json:"shipping_id"`
	Subtotal   float64   `json:"subtotal"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Rank struct {
	ID    int64  `json:"id"`
	Score int32  `json:"score"`
	Name  string `json:"name"`
}

type ResetPassword struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	ResetCode string    `json:"reset_code"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Review struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	BooksID   int64     `json:"books_id"`
	Liked     int32     `json:"liked"`
	Disliked  int32     `json:"disliked"`
	Reported  bool      `json:"reported"`
	Comments  string    `json:"comments"`
	IsDeleted bool      `json:"is_deleted"`
	Rating    int32     `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Shipping struct {
	ID        int64   `json:"id"`
	ToAddress string  `json:"to_address"`
	Total     float64 `json:"total"`
}

type Subgenre struct {
	ID        int64     `json:"id"`
	GenresID  int64     `json:"genres_id"`
	Name      string    `json:"name"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	OrdersID  int64     `json:"orders_id"`
	BooksID   int64     `json:"books_id"`
	CreatedAt time.Time `json:"created_at"`
	Amount    int32     `json:"amount"`
	Total     float64   `json:"total"`
	Reviewed  bool      `json:"reviewed"`
}

type User struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Image             string    `json:"image"`
	PhoneNumber       string    `json:"phone_number"`
	Age               int32     `json:"age"`
	Sex               string    `json:"sex"`
	Role              string    `json:"role"`
	Rank              int32     `json:"rank"`
	IsDeleted         bool      `json:"is_deleted"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type Wishlist struct {
	ID        int64     `json:"id"`
	BooksID   int64     `json:"books_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
