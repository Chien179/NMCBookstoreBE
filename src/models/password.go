package models

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	ID        int64  `json:"id" binding:"required,min=1"`
	ResetCode string `json:"reset_code" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
}
