package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBookSubgenreRequest struct {
	BookID     int64 `json:"book_id" binding:"required,min=1"`
	SubgenreID int64 `json:"subgenre_id" binding:"required,min=1"`
}

func (server *Server) createBookSubgenre(ctx *gin.Context) {
	var req createBookSubgenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBookSubgenreParams{
		BooksID:     req.BookID,
		SubgenresID: req.SubgenreID,
	}

	bookSubgenre, err := server.store.CreateBookSubgenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookSubgenre)
}

type getBookSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBookSubgenre(ctx *gin.Context) {
	var req getBookSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookSubgenre, err := server.store.GetBookSubgenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookSubgenre)
}

type deleteBookSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBookSubgenre(ctx *gin.Context) {
	var req deleteBookSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteBookSubgenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Book deleted successfully")
}

type listBooksSubgenresByBookIDRequest struct {
	BookID int64 `uri:"Book_id" binding:"required,min=1"`
}

func (server *Server) listBooksSubgenresByBookID(ctx *gin.Context) {
	var req listBooksSubgenresByBookIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookSubgenres, err := server.store.ListBooksSubgenresByBookID(ctx, req.BookID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookSubgenres)
}

type listBooksSubgenresBySubgenreIDRequest struct {
	SubgenreID int64 `uri:"Book_id" binding:"required,min=1"`
}

func (server *Server) listBooksSubgenresBySubgenreID(ctx *gin.Context) {
	var req listBooksSubgenresBySubgenreIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookSubgenres, err := server.store.ListBooksSubgenresBySubgenreID(ctx, req.SubgenreID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookSubgenres)
}
