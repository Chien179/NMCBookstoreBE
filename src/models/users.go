package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username    string                `form:"username" binding:"required,alphanum"`
	Password    string                `form:"password" binding:"required,min=8"`
	Fullname    string                `form:"full_name" binding:"required"`
	Email       string                `form:"email" binding:"required,email"`
	Image       *multipart.FileHeader `form:"image"`
	Age         int32                 `form:"age" binding:"required"`
	Sex         string                `form:"sex" binding:"required"`
	PhoneNumber string                `form:"phone_number" binding:"required"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Image             string    `json:"image"`
	Age               int32     `json:"age"`
	Sex               string    `json:"sex"`
	PhoneNumber       string    `json:"phone_number"`
	Role              string    `json:"role"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

type UpdateUserRequest struct {
	Fullname    string                `form:"full_name"`
	Email       string                `form:"email"`
	Image       *multipart.FileHeader `form:"image"`
	Age         int32                 `form:"age"`
	Sex         string                `form:"sex"`
	PhoneNumber string                `form:"phone_number"`
	Password    string                `form:"password"`
}
