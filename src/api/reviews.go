package api

import (
	"database/sql"
	"encoding/json"
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
		if transaction.BooksID == req.BookID && !transaction.Reviewed {
			isReviewed = true
			_, err := server.store.UpdateTransaction(ctx, transaction.ID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			break
		}
	}

	if !isReviewed && err != nil {
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

	rsp := models.ReviewsResponse{}
	rsp.TotalPage = int64(reviews.TotalPage)
	result := []models.ReviewData{}
	bt, _ := reviews.Reviews.MarshalJSON()
	_ = json.Unmarshal(bt, &result)

	for i, _ := range result {
		arg := db.GetLikeParams{
			ReviewID: result[i].ID,
			Username: req.Username,
		}

		isLike := false
		like, err := server.store.GetLike(ctx, arg)
		if err == nil {
			isLike = like.IsLike
		}
		result[i].Islike = isLike

		dislikeArg := db.GetDislikeParams{
			ReviewID: result[i].ID,
			Username: req.Username,
		}

		isDislike := false
		dislike, err := server.store.GetDislike(ctx, dislikeArg)
		if err == nil {
			isDislike = dislike.IsDislike
		}
		result[i].IsDislike = isDislike
	}

	rsp.Data = result

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) report(ctx *gin.Context) {
	var req models.ReportReviewRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateReviewParams{
		ID: req.ID,
		Reported: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}
	_, err := server.store.UpdateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Reported success")
}

func (server *Server) listAllReview(ctx *gin.Context) {
	reviews, err := server.store.ListReviews(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
