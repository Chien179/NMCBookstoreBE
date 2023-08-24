package models

type GetCityRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
