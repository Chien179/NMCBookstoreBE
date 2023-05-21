package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getCityRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCity(ctx *gin.Context) {
	var req getCityRequest
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
