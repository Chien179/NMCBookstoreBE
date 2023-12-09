package models

type GoogleOauthRequest struct {
	Code string `form:"code" binding:"required"`
}
