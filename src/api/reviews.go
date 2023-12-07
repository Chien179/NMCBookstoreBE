package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/token"
	"github.com/gin-gonic/gin"
)

func (server *Server) createReview(ctx *gin.Context) {
	var req models.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req.CreateReviewData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userOrders, err := server.store.ListOdersByUserName(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderTransactions := []db.Transaction{}
	for _, order := range userOrders {
		orderTransactions, err = server.store.ListTransactionsByOrderID(ctx, order.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	isReviewed := false
	for _, transaction := range orderTransactions {
		if transaction.BooksID == req.BookID && transaction.Reviewed == false {
			isReviewed = true
			_, err := server.store.UpdateTransaction(ctx, transaction.ID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			break
		}
	}

	if !isReviewed {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateReviewParams{
		Username: authPayLoad.Username,
		BooksID:  req.BookID,
		Comments: req.Comments,
		Rating:   req.Ratings,
	}

	review, err := server.store.CreateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (server *Server) deleteReview(ctx *gin.Context) {
	var req models.DeleteReviewRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.GetReview(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if review.Username != authPayLoad.Username {
		err := errors.New("review doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteReview(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Review deleted successfully")
}

func (server *Server) listReview(ctx *gin.Context) {
	var req models.ListReviewRequest
	if err := ctx.ShouldBindQuery(&req.ListReviewFormdata); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListReviewsByBookIDParams{
		BooksID: req.BookID,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	reviews, err := server.store.ListReviewsByBookID(ctx, arg)
	if err != nil {
		if reviews.Reviews == nil {
			ctx.JSON(http.StatusOK, reviews)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reviews)
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
	if dislike.IsDislike != true {
		dislikeAmount = -1
	}

	review, err := server.store.GetReview(ctx, dislike.ReviewID)
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
	}

	_, err = server.store.UpdateReview(ctx, rArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dislike)
}
