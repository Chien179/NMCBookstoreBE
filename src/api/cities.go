package api

import (
	"database/sql"
	"net/http"

	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) getCity(ctx *gin.Context) {
	var req models.GetCityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	city, err := server.store.GetCity(ctx, req.ID)
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

func (server *Server) listCities(ctx *gin.Context) {
	cities, err := server.store.ListCities(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cities)
}
