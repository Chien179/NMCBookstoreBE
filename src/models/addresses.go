package models

type CreateAddressRequest struct {
	Address    string `json:"address" binding:"required"`
	DistrictID int64  `json:"district_id" binding:"required"`
	CityID     int64  `json:"city_id" binding:"required"`
}

type GetAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateAddressData struct {
	Address    string `json:"address"`
	DistrictID int64  `json:"district_id" binding:"required"`
	CityID     int64  `json:"city_id" binding:"required"`
}

type UpdateAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	UpdateAddressData
}

type DeleteAddressRequest struct {
	IDs []int64 `form:"ids" binding:"required"`
}
