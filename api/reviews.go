package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type createReviewData struct {
	Comments string `json:"comments" binding:"required"`
	Ratings  int32  `json:"rating" binding:"required"`
}

type createReviewRequest struct {
	BookID int64 `uri:"book_id" binding:"required,min=1"`
	createReviewData
}

// @Summary      Create review
// @Description  Use this API to create review
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        book_id path int  true  "Create review id"
// @Param        Request body createReviewData  true  "Create review data"
// @Success      200 {object} db.Review
// @failure	 	 400
// @failure		 500
// @Router       /usersreviews/{book_id} [post]
func (server *Server) createReview(ctx *gin.Context) {
	var req createReviewRequest
	if err := ctx.ShouldBindJSON(&req.createReviewData); err != nil {
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
		Comments: req.createReviewData.Comments,
		Rating:   req.createReviewData.Ratings,
	}

	review, err := server.store.CreateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

type deleteReviewRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete review
// @Description  Use this API to delete review
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Delete review"
// @Success      200
// @failure	 	 400
// @failure	 	 401
// @failure	 	 404
// @failure		 500
// @Router       /users/reviews/delete/{id} [delete]
func (server *Server) deleteReview(ctx *gin.Context) {
	var req deleteReviewRequest
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

type listReviewFormdata struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listReviewRequest struct {
	BookID int64 `uri:"book_id" binding:"required,min=1"`
	listReviewFormdata
}

// @Summary      List review
// @Description  Use this API to list review
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        book_id path int  true  "List review id"
// @Param        Request query listReviewFormdata  true  "List review formdata"
// @Success      200 {object} []db.Review
// @failure	 	 400
// @failure		 500
// @Router       /usersreviews/{book_id} [get]
func (server *Server) listReview(ctx *gin.Context) {
	var req listReviewRequest
	if err := ctx.ShouldBindQuery(&req.listReviewFormdata); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListReviewsByBookIDParams{
		BooksID: req.BookID,
		Limit:   req.listReviewFormdata.PageSize,
		Offset:  (req.listReviewFormdata.PageID - 1) * req.listReviewFormdata.PageSize,
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
