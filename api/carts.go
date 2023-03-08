package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type addToCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) addToCart(ctx *gin.Context) {
	var req addToCartRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	book, err := server.store.GetBook(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateCartParams{
		BooksID:  book.ID,
		Username: authPayLoad.Username,
	}

	bookCart, err := server.store.CreateCart(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookCart)
}

type deleteBookInCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBookInCart(ctx *gin.Context) {
	var req deleteBookInCartRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCart(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Book in cart deleted successfully")
}

func (server *Server) listBookInCartByUsername(ctx *gin.Context) ([]db.Cart, error) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	books, err := server.store.ListCartsByUsername(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return nil, err
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return nil, err
	}

	return books, nil
}
