package models

type VerifyEmailRequest struct {
	EmailID    int64  `form:"email_id" binding:"required,min=1"`
	SecretCode string `form:"secret_code" binding:"required"`
}
