package models

import db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"

type BookResponse struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	Sale        float64       `json:"sale"`
	Image       []string      `json:"image"`
	Description string        `json:"description"`
	Author      string        `json:"author"`
	Publisher   string        `json:"publisher"`
	Quantity    int32         `json:"quantity"`
	Rating      float64       `json:"rating"`
	Genres      []db.Genre    `json:"genres"`
	Subgenres   []db.Subgenre `json:"subgenres"`
	IsDeleted   bool          `json:"is_deleted"`
}

type GetBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type DeleteBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListBookRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=24,max=100"`
}

type ListBookFollowGenreRequest struct {
	ID    int64 `form:"genre_id" binding:"required,min=1"`
	Limit int64 `form:"limit" binding:"required,min=6"`
}
