package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type cartData struct {
	Amount int32 `json:"amount" binding:"required,min=1"`
}

type addToCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	cartData
}

// @Summary      Add to cart
// @Description  Use this API to add to cart
// @Tags         Carts
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Add to cart"
// @Success      200 {object} db.Cart
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /users/add_to_cart/{id} [post]
func (server *Server) addToCart(ctx *gin.Context) {
	var req addToCartRequest
	if err := ctx.ShouldBindJSON(&req.cartData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	carts, err := server.store.ListCartsByUsername(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	for _, cart := range carts {
		if cart.BooksID == book.ID {
			newAmount := cart.Amount + req.Amount
			arg := db.UpdateAmountParams{
				ID:     cart.ID,
				Amount: newAmount,
				Total:  book.Price * float64(newAmount),
			}
			cart, err = server.store.UpdateAmount(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusOK, cart)
			return
		}
	}

	arg := db.CreateCartParams{
		BooksID:  book.ID,
		Username: authPayLoad.Username,
		Amount:   req.Amount,
		Total:    book.Price * float64(req.Amount),
	}

	bookCart, err := server.store.CreateCart(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookCart)
}

type updateAmountCartRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	cartData
}

func (server *Server) upatdeAmountCart(ctx *gin.Context) {
	var req updateAmountCartRequest
	if err := ctx.ShouldBindJSON(&req.cartData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cart, err := server.store.GetCart(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	book, err := server.store.GetBook(ctx, cart.BooksID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if cart.Username != authPayLoad.Username {
		err := errors.New("cart doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateAmountParams{
		ID:     cart.ID,
		Amount: req.Amount,
		Total:  book.Price * float64(req.Amount),
	}

	cart, err = server.store.UpdateAmount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// @Summary      Delete book in cart
// @Description  Use this API to delete book in cart
// @Tags         Carts
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Delete book in cart"
// @Success      200
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /users/delete_book_in_cart/{id} [delete]

type deleteMultiBookInCartRequest struct {
	IDs []int64 `form:"ids" binding:"required"`
}

func (server *Server) deleteBookInCart(ctx *gin.Context) {
	var req deleteMultiBookInCartRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	for _, cartID := range req.IDs {
		cart, err := server.store.GetCart(ctx, cartID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if cart.Username != authPayLoad.Username {
			err := errors.New("cart doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		arg := db.DeleteCartParams{
			ID:       cartID,
			Username: authPayLoad.Username,
		}

		err = server.store.DeleteCart(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "Books in cart deleted successfully")
}

type ListBooksInCartRespone struct {
	CartID   int64   `json:"cart_id"`
	BookID   int64   `json:"book_id"`
	BookName string  `json:"book_name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
	Amount   int32   `json:"amount"`
}

// @Summary      Get book in cart
// @Description  Use this API to get book in cart
// @Tags         Books
// @Accept       json
// @Produce      json
// @Success      200  {object}  []db.Book
// @failure	 	 400
// @failure		 404
// @failure		 500
// @Router       /users/list_book_in_cart [get]
func (server *Server) listBookInCart(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	carts, err := server.store.ListCartsByUsername(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []ListBooksInCartRespone{}
	for _, cart := range carts {
		book, err := server.store.GetBook(ctx, cart.BooksID)
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
		rsp = append(rsp, ListBooksInCartRespone{
			CartID:   cart.ID,
			BookID:   cart.BooksID,
			BookName: book.Name,
			Image:    book.Image[0],
			Price:    book.Price,
			Amount:   cart.Amount,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}
