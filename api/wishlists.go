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

// @Summary      Add to wislist
// @Description  Use this API to add to wishlist
// @Tags         Carts
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Add to wishlist"
// @Success      200 {object} db.Wishlist
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /users/add_to_wishlist/{id} [post]
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

// @Summary      Delete book in wishlist
// @Description  Use this API to delete book in wishlist
// @Tags         Carts
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Delete book in wishlist"
// @Success      200
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /users/add_to_wishlist/{id} [delete]
func (server *Server) deleteBookInWishlist(ctx *gin.Context) {
	var req deleteBookInWishlistRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.DeleteWishlistParams{
		ID:       req.ID,
		Username: authPayLoad.Username,
	}

	err := server.store.DeleteWishlist(ctx, arg)
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

type deleteMultiBookInWishlistRequest struct {
	List_ID []int64 `uri:"list_id" binding:"required"`
}

func (server *Server) deleteMutilBookInWishlist(ctx *gin.Context) {
	var req deleteMultiBookInWishlistRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	for _, wishlistID := range req.List_ID {
		arg := db.DeleteWishlistParams{
			ID:       wishlistID,
			Username: authPayLoad.Username,
		}

		err := server.store.DeleteWishlist(ctx, arg)
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

type ListBooksInWishlistRespone struct {
	WishlistID int64   `json:"wishlist_id"`
	Book       db.Book `json:"book"`
}

// @Summary      Get book in wishlist
// @Description  Use this API to get book in wishlis
// @Tags         Books
// @Accept       json
// @Produce      json
// @Success      200  {object}  []db.Book
// @failure	 	 400
// @failure		 404
// @failure		 500
// @Router       /users/list_book_in_wishlist [get]
func (server *Server) listBookInWishlist(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	wishlist, err := server.store.ListWishlistsByUsername(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := []ListBooksInWishlistRespone{}

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

		rsp = append(rsp, ListBooksInWishlistRespone{
			WishlistID: wish.ID,
			Book:       book,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}
