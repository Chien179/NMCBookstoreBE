package models

type RankRequest struct {
	Email string `form:"email" binding:"required"`
}

type RankReponse struct {
	Rank    string `json:"rank"`
	Vote    int    `josn:"vote"`
	Reviews int    `json:"reveiws"`
}
