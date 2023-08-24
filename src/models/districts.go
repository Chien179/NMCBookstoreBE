package models

type GetDistrictRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListDistrictsRequest struct {
	CityID int64 `uri:"city_id" binding:"required,min=1"`
}
