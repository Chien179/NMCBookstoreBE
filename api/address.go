package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAddressRequest struct {
	Address  string `json:"address" binding:"required"`
	District string `json:"district" binding:"required"`
	City     string `json:"city" binding:"required"`
}

// @Summary      Create address
// @Description  Use this API to create address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        Request body createAddressRequest  true  "Create address"
// @Success      200  {object}  db.Address
// @failure	 	 400
// @failure		 403
// @failure		 500
// @Router       /users/addresses [post]
func (server *Server) createAddress(ctx *gin.Context) {
	var req createAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAddressParams{
		Username: authPayLoad.Username,
		Address:  req.Address,
		District: req.District,
		City:     req.City,
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

type getAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Get address
// @Description  Use this API to get address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "Get address"
// @Success      200  {object}  db.Address
// @failure	 	 400
// @failure		 401
// @failure		 403
// @failure		 404
// @failure		 500
// @Router       /users/addresses/{id} [get]
func (server *Server) getAddress(ctx *gin.Context) {
	var req getAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

type updateAddressData struct {
	Address  string `json:"address"`
	District string `json:"district"`
	City     string `json:"city"`
}

type updateAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	updateAddressData
}

// @Summary      Update address
// @Description  Use this API to update address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "Update address id"
// @Param        Request body updateAddressData  false  "Update address data"
// @Success      200  {object}  db.Address
// @failure	 	 400
// @failure		 401
// @failure		 403
// @failure		 404
// @failure		 500
// @Router       /users/addresses/update/{id} [put]
func (server *Server) updateAddress(ctx *gin.Context) {
	var req updateAddressRequest
	if err := ctx.ShouldBindJSON(&req.updateAddressData); err != nil {
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
			String: req.updateAddressData.Address,
			Valid:  req.updateAddressData.Address != "",
		},
		District: sql.NullString{
			String: req.updateAddressData.District,
			Valid:  req.updateAddressData.District != "",
		},
		City: sql.NullString{
			String: req.updateAddressData.City,
			Valid:  req.updateAddressData.District != "",
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

type deleteAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete address
// @Description  Use this API to delete address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "Delete address"
// @Success      200
// @failure	 	 400
// @failure		 401
// @failure		 404
// @failure		 500
// @Router       /users/addresses/delete/{id} [delete]
func (server *Server) deleteAddress(ctx *gin.Context) {
	var req deleteAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	err = server.store.DeleteAddress(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Address deleted successfully")
}

type listAddressRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// @Summary      List address
// @Description  Use this API to list address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        Query query listAddressRequest  true  "List address"
// @Success      200 {object}  []db.Address
// @failure	 	 400
// @failure		 401
// @failure		 404
// @failure		 500
// @Router       /users/addresses [get]
func (server *Server) listAddress(ctx *gin.Context) {
	var req listAddressRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAddressesParams{
		Username: authPayLoad.Username,
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
	}

	addresses, err := server.store.ListAddresses(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		if addresses.Address == nil {
			ctx.JSON(http.StatusOK, addresses)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, addresses)
}
