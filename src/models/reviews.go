package models

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
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListReviewRequest struct {
	BookID int64 `uri:"book_id" binding:"required,min=1"`
	ListReviewFormdata
}
