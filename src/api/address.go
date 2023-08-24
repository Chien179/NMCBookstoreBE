package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (server *Server) createAddress(ctx *gin.Context) {
	var req models.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAddressParams{
		Username:   authPayLoad.Username,
		Address:    req.Address,
		DistrictID: req.DistrictID,
		CityID:     req.CityID,
	}

	address, err := server.store.CreateAddress(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (server *Server) getAddress(ctx *gin.Context) {
	var req models.GetAddressRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	address, err := server.store.GetAddress(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if address.Username != authPayLoad.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (server *Server) updateAddress(ctx *gin.Context) {
	var req models.UpdateAddressRequest
	if err := ctx.ShouldBindJSON(&req.UpdateAddressData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	address, err := server.store.GetAddress(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if address.Username != authPayLoad.Username {
		err := errors.New("address doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	arg := db.UpdateAddressParams{
		ID: req.ID,
		Address: sql.NullString{
			String: req.Address,
			Valid:  req.Address != "",
		},
		DistrictID: sql.NullInt64{
			Int64: req.DistrictID,
			Valid: req.DistrictID != address.DistrictID,
		},
		CityID: sql.NullInt64{
			Int64: req.CityID,
			Valid: req.CityID != address.CityID,
		},
	}

	address, err = server.store.UpdateAddress(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (server *Server) deleteAddress(ctx *gin.Context) {
	var req models.DeleteAddressRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, addressID := range req.IDs {
		address, err := server.store.GetAddress(ctx, addressID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if address.Username != authPayLoad.Username {
			err := errors.New("account doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		err = server.store.DeleteAddress(ctx, addressID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	}

	ctx.JSON(http.StatusOK, "Addresses deleted successfully")
}

func (server *Server) listAddress(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	addresses, err := server.store.ListAddresses(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		if addresses == nil {
			ctx.JSON(http.StatusOK, addresses)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, addresses)
}
