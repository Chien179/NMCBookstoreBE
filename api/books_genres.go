package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBookGenreRequest struct {
	BookID  int64 `json:"book_id" binding:"required,min=1"`
	GenreID int64 `json:"subgenre_id" binding:"required,min=1"`
}

func (server *Server) createBookGenre(ctx *gin.Context) {
	var req createBookGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBookGenreParams{
		BooksID:  req.BookID,
		GenresID: req.GenreID,
	}

	bookGenre, err := server.store.CreateBookGenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookGenre)
}

type getBookGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBookGenre(ctx *gin.Context) {
	var req getBookGenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookGenre, err := server.store.GetBookGenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookGenre)
}

type deleteBookGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBookGenre(ctx *gin.Context) {
	var req deleteBookGenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteBookGenre(ctx, req.ID)
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

type listBooksGenresByBookIDRequest struct {
	BookID int64 `uri:"Book_id" binding:"required,min=1"`
}

func (server *Server) listBooksGenresByBookID(ctx *gin.Context) {
	var req listBooksGenresByBookIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookGenres, err := server.store.ListBooksGenresByBookID(ctx, req.BookID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookGenres)
}

type listBooksGenresByGenreIDRequest struct {
	GenreID int64 `uri:"Book_id" binding:"required,min=1"`
}

func (server *Server) listBooksGenresByGenreID(ctx *gin.Context) {
	var req listBooksGenresByGenreIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookGenres, err := server.store.ListBooksGenresByGenreID(ctx, req.GenreID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bookGenres)
}
