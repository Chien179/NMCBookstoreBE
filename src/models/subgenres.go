package models

type CreateSubgenreRequest struct {
	GenresID int64  `json:"genres_id" binding:"required,min=1"`
	Name     string `json:"name" binding:"required"`
}

type GetSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateSubgenreData struct {
	GenreID int64  `json:"genre_id"`
	Name    string `json:"name"`
}

type UpdateSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	UpdateSubgenreData
}

type DeleteSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListSubgenreRequest struct {
	GenreID int64 `uri:"genre_id" binding:"required,min=1"`
}
