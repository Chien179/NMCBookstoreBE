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

func (server *Server) createOrder(ctx *gin.Context) {
	var req models.CreateOrderRequest
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
	sale := float64(0)
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

		sale += book.Sale

		argBook := db.UpdateBookParams{
			ID: book.ID,
			Quantity: sql.NullInt32{
				Int32: book.Quantity - cart.Amount,
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
			Float64: subTotal + req.TotalShipping,
			Valid:   subTotal > 0,
		},
		Sale: sql.NullFloat64{
			Float64: sale,
			Valid:   true,
		},
		Status: sql.NullString{
			String: req.Status,
			Valid:  true,
		},
		Note: sql.NullString{
			String: req.Note,
			Valid:  req.Note != "",
		},
	}

	order, err = server.store.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.createPayment(ctx, req.PaymentID, order.ID, req.ToAddress, req.TotalShipping, order.SubTotal, req.Status, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := models.OrderReponse{
		ID:        order.ID,
		Username:  order.Username,
		ToAddress: req.ToAddress,
		Note:      order.Note.String,
		SubAmount: order.SubAmount,
		SubTotal:  order.SubTotal,
		Sale:      order.Sale,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteOrder(ctx *gin.Context) {
	var req models.DeleteOrderRequest
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
	var req models.DeleteOrderRequest
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

func (server *Server) listOrder(ctx *gin.Context) {
	var req models.ListOrderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListOdersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orders, err := server.store.ListOders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []models.ListOrderResponse{}

	listOrder := []db.Order{}
	json.Unmarshal([]byte(orders.Orders), &listOrder)

	for _, order := range listOrder {
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

		rsp = append(rsp, models.ListOrderResponse{
			ID:           order.ID,
			Username:     order.Username,
			Books:        books,
			Transactions: transactions,
			Status:       order.Status,
			SubTotal:     order.SubTotal,
			Sale:         float64(order.Sale),
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

	rsp := []models.ListOrderResponse{}

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

			rsp = append(rsp, models.ListOrderResponse{
				ID:           order.ID,
				Username:     order.Username,
				Books:        books,
				Transactions: transactions,
				Status:       order.Status,
				SubTotal:     order.SubTotal,
				Sale:         float64(order.Sale),
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

	rsp := []models.ListOrderResponse{}

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

			rsp = append(rsp, models.ListOrderResponse{
				ID:           order.ID,
				Username:     order.Username,
				Books:        books,
				Transactions: transactions,
				Status:       order.Status,
				SubTotal:     order.SubTotal,
				Sale:         float64(order.Sale),
				SubAmount:    order.SubAmount,
			})
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listAllOrder(ctx *gin.Context) {
	orders, err := server.store.ListAllOders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []models.ListOrderResponse{}

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

		rsp = append(rsp, models.ListOrderResponse{
			ID:           order.ID,
			Username:     order.Username,
			Books:        books,
			Transactions: transactions,
			Status:       order.Status,
			SubTotal:     order.SubTotal,
			Sale:         float64(order.Sale),
			SubAmount:    order.SubAmount,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}
