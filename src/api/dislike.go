package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) getDislikeReview(ctx *gin.Context) {
	var req models.DislikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetDislikeParams{
		Username: req.Username,
		ReviewID: req.ReviewId,
	}

	dislike, err := server.store.GetDislike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dislike)
}

func (server *Server) listDislike(ctx *gin.Context) {
	var req models.ListdisLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dislikes, err := server.store.ListDislike(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dislikes)
}

func (server *Server) dislikeReview(ctx *gin.Context) {
	var req models.DislikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetDislikeParams{
		Username: req.Username,
		ReviewID: req.ReviewId,
	}

	dislike, err := server.store.GetDislike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreatedDislikeParams{
				Username:  req.Username,
				ReviewID:  req.ReviewId,
				IsDislike: true,
			}
			dislike, err := server.store.CreatedDislike(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusOK, dislike)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else {
		uArg := db.UpdateDislikeParams{
			Username:  req.Username,
			ReviewID:  req.ReviewId,
			IsDislike: !dislike.IsDislike,
		}

		dislike, err = server.store.UpdateDislike(ctx, uArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	dislikeAmount := 1
	if !dislike.IsDislike {
		dislikeAmount = -1
	}

	if dislikeAmount == 1 {
		likeArg := db.UpdateLikeParams{
			Username: req.Username,
			ReviewID: req.ReviewId,
			IsLike:   false,
		}

		_, err := server.store.GetLike(ctx, db.GetLikeParams{
			Username: req.Username,
			ReviewID: req.ReviewId,
		})

		if err == nil {
			server.store.UpdateLike(ctx, likeArg)
		}

	}

	review, err := server.store.GetReview(ctx, req.ReviewId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rArg := db.UpdateReviewParams{
		Disliked: sql.NullInt32{
			Int32: review.Disliked + int32(dislikeAmount),
			Valid: true,
		},
		ID: req.ReviewId,
	}

	_, err = server.store.UpdateReview(ctx, rArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dislike)
}