package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type addToWishlistRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) addToWishlist(ctx *gin.Context) {
	var req addToWishlistRequest
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
	arg := db.CreateWishlistParams{
		BooksID:  book.ID,
		Username: authPayLoad.Username,
	}

	bookWishlist, err := server.store.CreateWishlist(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookWishlist)
}

type deleteBookInWishlistRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBookInWishlist(ctx *gin.Context) {
	var req deleteBookInWishlistRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteWishlist(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Book in wishlist deleted successfully")
}

func (server *Server) listBookInWishlistByUsername(ctx *gin.Context) ([]db.Wishlist, error) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	books, err := server.store.ListWishlistsByUsername(ctx, authPayLoad.Username)
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
