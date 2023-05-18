package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type FullSearchRequest struct {
	PageID      int32   `form:"page_id" binding:"required,min=1"`
	PageSize    int32   `form:"page_size" binding:"required,min=24,max=100"`
	Text        string  `form:"text"`
	GenresID    int64   `form:"genres_id"`
	SubgenresID int64   `form:"subgenres_id"`
	MinPrice    float64 `form:"min_price"`
	MaxPrice    float64 `form:"max_price"`
	Rating      float64 `form:"rating"`
}

func (server *Server) fullSearch(ctx *gin.Context) {
	var req FullSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.FullSearchParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Text: sql.NullString{
			String: req.Text,
			Valid:  req.Text != "",
		},
		GenresID: sql.NullInt64{
			Int64: req.GenresID,
			Valid: req.GenresID > 0,
		},
		SubgenresID: sql.NullInt64{
			Int64: req.SubgenresID,
			Valid: req.SubgenresID > 0,
		},
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Rating: sql.NullFloat64{
			Float64: req.Rating,
			Valid:   req.Rating >= 0,
		},
	}

	results, err := server.store.FullSearch(ctx, arg)
	if err != nil {
		if results.Books == nil {
			ctx.JSON(http.StatusOK, results)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, results)
}
