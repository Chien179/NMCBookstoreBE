package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *Server) getDislikeReview(ctx *gin.Context) {
	var req models.DislikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetDislikeParams{
		Username: authPayLoad.Username,
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

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetDislikeParams{
		Username: authPayLoad.Username,
		ReviewID: req.ReviewId,
	}

	dislike, err := server.store.GetDislike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreatedDislikeParams{
				Username:  authPayLoad.Username,
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
			Username:  authPayLoad.Username,
			ReviewID:  req.ReviewId,
			IsDislike: !dislike.IsDislike,
		}

		dislike, err = server.store.UpdateDislike(ctx, uArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	like, err := server.store.GetLike(ctx, db.GetLikeParams{
		Username: authPayLoad.Username,
		ReviewID: req.ReviewId,
	})
	log.Print(like)

	dislikeAmount := 1
	if !dislike.IsDislike {
		dislikeAmount = -1
	} else if err == nil {
		likeArg := db.UpdateLikeParams{
			Username: authPayLoad.Username,
			ReviewID: req.ReviewId,
			IsLike:   false,
		}

		server.store.UpdateLike(ctx, likeArg)
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
