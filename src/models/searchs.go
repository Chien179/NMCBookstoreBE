package models

type FullSearchRequest struct {
	PageID      int32   `form:"page_id" binding:"required,min=1"`
	PageSize    int32   `form:"page_size" binding:"required,min=24,max=100"`
	Text        string  `form:"text"`
	GenresID    int64   `form:"genres_id"`
	SubgenresID int64   `form:"subgenres_id"`
	MinPrice    float64 `form:"min_price"`
	MaxPrice    float64 `form:"max_price"`
	Rating      float64 `form:"rating"`
}

type RecommedRequest struct {
	BookID      int64   `form:"book_id"  binding:"required"`
	GenresID    []int64 `form:"genres_id"  binding:"required"`
	SubgenresID []int64 `form:"subgenres_id"  binding:"required"`
}
