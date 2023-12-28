package models

type DislikeResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	ReviewId  int64  `json:"review_id"`
	IsDislike bool   `json:"is_dislike"`
}

type DislikeRequest struct {
	ReviewId int64 `json:"review_id" binding:"required"`
}

type ListdisLikeRequest struct {
	Username string `uri:"username" binding:"required"`
}
