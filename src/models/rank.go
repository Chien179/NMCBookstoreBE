package models

type RankRequest struct {
	Email string `json:"email" binding:"required"`
}

type RankReponse struct {
	Rank    string `json:"rank"`
	Vote    int    `json:"vote"`
	Reviews int    `json:"review"`
}
