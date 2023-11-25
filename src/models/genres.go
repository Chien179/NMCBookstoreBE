package models

type GetGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	CreateGenreRequest
}

type DeleteGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
