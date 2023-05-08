package api

import (
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type FullSearchRequest struct {
	PageID   int32   `form:"page_id" binding:"required,min=1"`
	PageSize int32   `form:"page_size" binding:"required,min=24,max=100"`
	Text     string  `form:"text"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
	Rating   float64 `form:"rating"`
}

func (server *Server) fullSearch(ctx *gin.Context) {
	var req FullSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
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
