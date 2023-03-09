package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createTransactionRequest struct {
	OrderID int64 `json:"order_id" binding:"required,min=1"`
	BookID  int64 `json:"book_id" binding:"required,min=1"`
}

func (server *Server) createTransaction(ctx *gin.Context) {
	var req createTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTransactionParams{
		OrdersID: req.OrderID,
		BooksID:  req.BookID,
	}

	transaction, err := server.store.CreateTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

type getTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransaction(ctx *gin.Context) {
	var req getTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := server.store.GetTransaction(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

type listTransactionRequest struct {
	OrderID int64 `uri:"order_id" binding:"required,min=1"`
}

func (server *Server) listTransaction(ctx *gin.Context) {
	var req listTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := server.store.ListTransactionsByOrderID(ctx, req.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}
