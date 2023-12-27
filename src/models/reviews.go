package models

import "time"

type CreateReviewData struct {
	Comments string `json:"comments" binding:"required"`
	Ratings  int32  `json:"rating" binding:"required"`
}

type CreateReviewRequest struct {
	BookID int64 `uri:"book_id" binding:"required,min=1"`
	CreateReviewData
}

type DeleteReviewRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListReviewFormdata struct {
	Username string `form:"page_id" binding:"required"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
}

type ListReviewRequest struct {
	BookID int64 `uri:"book_id" binding:"required,min=1"`
	ListReviewFormdata
}

type ReviewsResponse struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Image     string    `json:"image"`
	BooksId   string    `json:"books_id"`
	Comments  string    `json:"comments"`
	Rating    float64   `json:"rating"`
	Islike    bool      `json:"is_like"`
	IsDislike bool      `json:"is_dislike"`
	CreatedAt time.Time `json:"create_at"`
}
