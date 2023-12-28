package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/gin-gonic/gin"
)

func (server *Server) getLikeReview(ctx *gin.Context) {
	var req models.LikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetLikeParams{
		Username: authPayLoad.Username,
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

func (server *Server) listLike(ctx *gin.Context) {
	var req models.ListLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	likes, err := server.store.ListLike(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, likes)
}

func (server *Server) likeReview(ctx *gin.Context) {
	var req models.LikeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetLikeParams{
		Username: authPayLoad.Username,
		ReviewID: req.ReviewId,
	}

	like, err := server.store.GetLike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreateLikeParams{
				Username: authPayLoad.Username,
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
			Username: authPayLoad.Username,
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
	if !like.IsLike {
		likeAmount = -1
	}

	if likeAmount == 1 {
		dislikeArg := db.UpdateDislikeParams{
			Username:  authPayLoad.Username,
			ReviewID:  req.ReviewId,
			IsDislike: false,
		}

		_, err := server.store.GetDislike(ctx, db.GetDislikeParams{
			Username: authPayLoad.Username,
			ReviewID: req.ReviewId,
		})

		if err == nil {
			server.store.UpdateDislike(ctx, dislikeArg)
		}

	}

	review, err := server.store.GetReview(ctx, req.ReviewId)

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
		ID: req.ReviewId,
	}

	_, err = server.store.UpdateReview(ctx, rArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, like)
}
