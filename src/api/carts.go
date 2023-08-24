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

func (server *Server) addToCart(ctx *gin.Context) {
	var req models.AddToCartRequest
	if err := ctx.ShouldBindJSON(&req.Amount); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req.ID); err != nil {
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

func (server *Server) upatdeAmountCart(ctx *gin.Context) {
	var req models.UpdateAmountCartRequest
	if err := ctx.ShouldBindJSON(&req.Amount); err != nil {
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

func (server *Server) deleteBookInCart(ctx *gin.Context) {
	var req models.DeleteMultiBookInCartRequest
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

	rsp := []models.ListBooksInCartRespone{}
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
		rsp = append(rsp, models.ListBooksInCartRespone{
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
