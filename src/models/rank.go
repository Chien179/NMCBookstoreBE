package models

type RankRequest struct {
	Email string `form:"email" binding:"required"`
}

type RankReponse struct {
	Rank string `json:"rank"`
}
