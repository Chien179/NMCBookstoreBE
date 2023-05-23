package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getDistrictRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getDistrict(ctx *gin.Context) {
	var req getDistrictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	city, err := server.store.GetDistrict(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, city)
}

type listDistrictsRequest struct {
	CityID int64 `uri:"city_id" binding:"required,min=1"`
}

func (server *Server) listDistricts(ctx *gin.Context) {
	var req listDistrictsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cities, err := server.store.ListDistricts(ctx, req.CityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cities)
}
