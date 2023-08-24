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

func (server *Server) addToWishlist(ctx *gin.Context) {
	var req models.AddToWishlistRequest
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

func (server *Server) deleteBookInWishlist(ctx *gin.Context) {
	var req models.DeleteMultiBookInWishlistRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	for _, wishlistID := range req.IDs {
		wishlist, err := server.store.GetWishlist(ctx, wishlistID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if wishlist.Username != authPayLoad.Username {
			err := errors.New("cart doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		arg := db.DeleteWishlistParams{
			ID:       wishlistID,
			Username: authPayLoad.Username,
		}

		err = server.store.DeleteWishlist(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "Book in wishlist deleted successfully")
}

func (server *Server) listBookInWishlist(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	wishlist, err := server.store.ListWishlistsByUsername(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := []models.ListBooksInWishlistRespone{}

	for _, wish := range wishlist {
		book, err := server.store.GetBook(ctx, wish.BooksID)
		if err != nil {
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		rsp = append(rsp, models.ListBooksInWishlistRespone{
			WishlistID: wish.ID,
			Book:       book,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}
