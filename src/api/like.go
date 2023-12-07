package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) getLikeReview(ctx *gin.Context) {
	var req models.LikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetLikeParams{
		Username: req.Username,
		ReviewID: req.ReviewId,
	}

	like, err := server.store.GetLike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, like)
}

func (server *Server) likeReview(ctx *gin.Context) {
	var req models.LikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetLikeParams{
		Username: req.Username,
		ReviewID: req.ReviewId,
	}

	like, err := server.store.GetLike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreateLikeParams{
				Username: req.Username,
				ReviewID: req.ReviewId,
				IsLike:   true,
			}
			like, err := server.store.CreateLike(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusOK, like)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else {
		uArg := db.UpdateLikeParams{
			Username: req.Username,
			ReviewID: req.ReviewId,
			IsLike:   !like.IsLike,
		}

		like, err = server.store.UpdateLike(ctx, uArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	likeAmount := 1
	if like.IsLike != true {
		likeAmount = -1
	}

	review, err := server.store.GetReview(ctx, like.ReviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rArg := db.UpdateReviewParams{
		Liked: sql.NullInt32{
			Int32: review.Liked + int32(likeAmount),
			Valid: true,
		},
	}

	_, err = server.store.UpdateReview(ctx, rArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, like)
}
