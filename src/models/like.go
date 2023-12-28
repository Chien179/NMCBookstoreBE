package models

type LikeResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	ReviewId int64  `json:"review_id"`
	IsLike   bool   `json:"is_like"`
}

type LikeRequest struct {
	ReviewId int64 `json:"review_id" binding:"required"`
}

type ListLikeRequest struct {
	Username string `uri:"username" binding:"required"`
}
