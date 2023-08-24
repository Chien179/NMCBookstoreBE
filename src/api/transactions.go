package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) listBookByTransactions(ctx *gin.Context, transactions []db.Transaction) ([]db.Book, error) {
	books := []db.Book{}
	for _, transaction := range transactions {
		book, err := server.store.GetBook(ctx, transaction.BooksID)
		if err != nil {
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return nil, err
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return nil, err
			}
		}
		books = append(books, book)
	}

	return books, nil
}
