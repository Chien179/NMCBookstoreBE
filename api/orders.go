package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/token"
	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	PaymentID     string  `json:"payment_id" binding:"required"`
	CartIDs       []int64 `json:"cart_ids" binding:"required"`
	ToAddress     string  `json:"to_address" binding:"required"`
	TotalShipping float64 `json:"total_shipping" binding:"required,min=1000,max=100000000"`
	Status        string  `json:"status" binding:"required"`
}

// @Summary      Create order
// @Description  Use this API to create order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        Request body createOrderRequest  true  "Create order"
// @Success      200 {object} db.Order
// @failure	 	 400
// @failure		 500
// @Router       /users/orders [post]
func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Status == "failed" {
		ctx.JSON(http.StatusInternalServerError, "Payment failed")
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	order, err := server.store.CreateOrder(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	subTotal := 0.0
	for _, cartID := range req.CartIDs {
		cart, err := server.store.GetCart(ctx, cartID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg := db.CreateTransactionParams{
			OrdersID: order.ID,
			BooksID:  cart.BooksID,
			Amount:   cart.Amount,
			Total:    cart.Total,
		}

		_, err = server.store.CreateTransaction(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		subTotal += cart.Total

		book, err := server.store.GetBook(ctx, cart.BooksID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argBook := db.UpdateBookParams{
			ID: book.ID,
			Quantity: sql.NullInt32{
				Int32: book.Quantity - 1,
				Valid: true,
			},
			Image: book.Image,
		}
		_, err = server.store.UpdateBook(ctx, argBook)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argCart := db.DeleteCartParams{
			ID:       cart.ID,
			Username: authPayLoad.Username,
		}

		err = server.store.DeleteCart(ctx, argCart)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := db.UpdateOrderParams{
		ID: order.ID,
		SubAmount: sql.NullInt32{
			Int32: int32(len(req.CartIDs)),
			Valid: int32(len(req.CartIDs)) > 0,
		},
		SubTotal: sql.NullFloat64{
			Float64: subTotal,
			Valid:   subTotal > 0,
		},
	}

	order, err = server.store.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.createPayment(ctx, req.PaymentID, order.ID, req.ToAddress, req.TotalShipping, subTotal, req.Status)

	ctx.JSON(http.StatusOK, order)
}

type deleteOrderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete order
// @Description  Use this API to delete order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "Delete order"
// @Success      200
// @failure	 	 400
// @failure	 	 401
// @failure	 	 404
// @failure		 500
// @Router       /users/orders/delete/{id} [delete]
func (server *Server) deleteOrder(ctx *gin.Context) {
	var req deleteOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if order.Username != authPayLoad.Username {
		err := errors.New("order doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteOrder(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "deleted successfully")
}

func (server *Server) cancelOrder(ctx *gin.Context) {
	var req deleteOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if order.Username != authPayLoad.Username {
		err := errors.New("order doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateOrderParams{
		ID: req.ID,
		Status: sql.NullString{
			String: "cancelled",
			Valid:  true,
		},
	}

	order, err = server.store.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

type listOrderResponse struct {
	ID           int64            `uri:"id"`
	Books        []db.Book        `json:"books"`
	Transactions []db.Transaction `json:"transactions"`
	Status       string           `json:"status"`
	SubTotal     float64          `json:"sub_total"`
	SubAmount    int32            `json:"sub_amount"`
}

func (server *Server) listOrder(ctx *gin.Context) {
	orders, err := server.store.ListOders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []listOrderResponse{}

	for _, order := range orders {
		transactions, err := server.store.ListTransactionsByOrderID(ctx, order.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		books, err := server.listBookByTransactions(ctx, transactions)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp = append(rsp, listOrderResponse{
			ID:           order.ID,
			Books:        books,
			Transactions: transactions,
			Status:       order.Status,
			SubTotal:     order.SubTotal,
			SubAmount:    order.SubAmount,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listOrderPaid(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	orders, err := server.store.ListOdersByUserName(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []listOrderResponse{}

	for _, order := range orders {
		if order.Status == "paid" {
			transactions, err := server.store.ListTransactionsByOrderID(ctx, order.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			books, err := server.listBookByTransactions(ctx, transactions)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			rsp = append(rsp, listOrderResponse{
				ID:           order.ID,
				Books:        books,
				Transactions: transactions,
				Status:       order.Status,
				SubTotal:     order.SubTotal,
				SubAmount:    order.SubAmount,
			})
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listOrderCancelled(ctx *gin.Context) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	orders, err := server.store.ListOdersByUserName(ctx, authPayLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []listOrderResponse{}

	for _, order := range orders {
		if order.Status == "cancelled" {
			transactions, err := server.store.ListTransactionsByOrderID(ctx, order.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			books, err := server.listBookByTransactions(ctx, transactions)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			rsp = append(rsp, listOrderResponse{
				ID:           order.ID,
				Books:        books,
				Transactions: transactions,
				Status:       order.Status,
				SubTotal:     order.SubTotal,
				SubAmount:    order.SubAmount,
			})
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}
