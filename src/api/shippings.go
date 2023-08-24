package api

import (
	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) createShipping(ctx *gin.Context, toAddress string, total float64) (*db.Shipping, error) {
	arg := db.CreateShippingParams{
		ToAddress: toAddress,
		Total:     total,
	}

	shipping, err := server.store.CreateShipping(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &shipping, nil
}
