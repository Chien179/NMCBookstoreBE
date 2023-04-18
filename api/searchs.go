package api

import (
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type FullSearchQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=24,max=100"`
}

type FullSearchData struct {
	Text     string  `json:"text" binding:"required"`
	MinPrice float64 `json:"min_price" binding:"required,min=1000,max=100000000"`
	MaxPrice float64 `json:"max_price" binding:"required,min=1000,max=100000000"`
	Rating   float64 `json:"rating" binding:"required,min=0,max=5"`
}

type FullSearchRequest struct {
	FullSearchQuery
	FullSearchData
}

func (server *Server) fullSearch(ctx *gin.Context) {
	var req FullSearchRequest
	if err := ctx.ShouldBindQuery(&req.FullSearchQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.FullSearchData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.FullSearchParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		Text:     req.Text,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Rating:   req.Rating,
	}

	results, err := server.store.FullSearch(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, results)
}
