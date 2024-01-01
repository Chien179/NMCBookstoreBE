package models

import (
	"time"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
)

type CreateOrderRequest struct {
	PaymentID     string  `json:"payment_id" binding:"required"`
	CartIDs       []int64 `json:"cart_ids" binding:"required"`
	ToAddress     string  `json:"to_address" binding:"required"`
	Note          string  `json:"note"`
	Email         string  `json:"email" binding:"required"`
	TotalShipping float64 `json:"total_shipping" binding:"required,min=1000,max=100000000"`
	Status        string  `json:"status" binding:"required"`
}

type OrderReponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	ToAddress string    `json:"to_address"`
	Note      string    `json:"note"`
	SubAmount int32     `json:"sub_amount"`
	SubTotal  float64   `json:"sub_total"`
	Sale      float64   `json:"sale"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteOrderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListOrderRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=24,max=100"`
}

type ListOrderResponse struct {
	ID           int64            `json:"id"`
	Username     string           `json:"username"`
	Books        []db.Book        `json:"books"`
	Transactions []db.Transaction `json:"transactions"`
	Status       string           `json:"status"`
	SubTotal     float64          `json:"sub_total"`
	Sale         float64          `json:"sale"`
	SubAmount    int32            `json:"sub_amount"`
}
