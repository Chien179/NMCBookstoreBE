package api

import (
	"database/sql"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/gin-gonic/gin"
)

func (server *Server) createPayment(ctx *gin.Context, PaymentID string, OrderID int64, ToAddress string, TotalShipping float64, SubTotal float64, Status string, email string) error {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	shipping, err := server.createShipping(ctx, ToAddress, TotalShipping)
	if err != nil {
		return err
	}

	order, err := server.store.GetOrder(ctx, OrderID)
	if err != nil {
		return err
	}

	if Status != "failed" {
		arg := db.UpdateOrderParams{
			ID: order.ID,
			Status: sql.NullString{
				String: "paid",
				Valid:  true,
			},
		}

		_, err := server.store.UpdateOrder(ctx, arg)
		if err != nil {
			return err
		}

		user, err := server.store.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		userArg := db.UpdateUserParams{
			Username: authPayLoad.Username,
			Rank: sql.NullInt32{
				Int32: user.Rank + int32(SubTotal),
				Valid: true,
			},
		}

		_, err = server.store.UpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
	}

	arg := db.CreatePaymentParams{
		ID:         PaymentID,
		Username:   authPayLoad.Username,
		OrderID:    order.ID,
		ShippingID: shipping.ID,
		Subtotal:   order.SubTotal,
		Status:     Status,
	}

	_, err = server.store.CreatePayment(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}
