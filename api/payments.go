package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type CreatePaymentRequest struct {
	FromAddress   string  `json:"from_address" binding:"required"`
	ToAddress     string  `json:"to_address" binding:"required"`
	TotalShipping float64 `json:"total_shipping" binding:"required,min=1000,max=100000000"`
	SubTotal      float64 `json:"sub_total" binding:"required,min=1000,max=100000000"`
	Status        string  `json:"status" binding:"required"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	shipping, err := server.createShipping(ctx, req.FromAddress, req.ToAddress, req.TotalShipping)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	order, err := server.store.GetOrderToPayment(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Status != "failed" {
		arg := db.UpdateOrderParams{
			ID: order.ID,
			Status: sql.NullString{
				String: "paid",
				Valid:  true,
			},
		}
		_, err := server.store.UpdateOrder(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := db.CreatePaymentParams{
		Username:   authPayLoad.Username,
		OrderID:    order.ID,
		ShippingID: shipping.ID,
		Subtotal:   req.SubTotal,
		Status:     req.Status,
	}

	payment, err := server.store.CreatePayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
