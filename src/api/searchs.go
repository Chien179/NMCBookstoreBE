package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) fullSearch(ctx *gin.Context) {
	var req models.FullSearchRequest
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

func (server *Server) recommend(ctx *gin.Context) {
	var req models.RecommedRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.RecommendParams{
		BooksID:     req.BookID,
		GenresID:    req.GenresID,
		SubgenresID: req.SubgenresID,
	}

	books, err := server.store.Recommend(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (server *Server) justForYou(ctx *gin.Context) {
	books, err := server.store.JustForYou(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
}
