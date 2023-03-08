package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBookRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Image       string  `json:"image" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Publisher   string  `json:"publisher" binding:"required"`
	Quantity    int32   `json:"quanlity" binding:"required"`
}

func (server *Server) createBook(ctx *gin.Context) {
	var req createBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBookParams{
		Name:        req.Name,
		Price:       req.Price,
		Image:       req.Image,
		Description: req.Description,
		Author:      req.Author,
		Publisher:   req.Publisher,
		Quantity:    req.Quantity,
	}

	book, err := server.store.CreateBook(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, book)
}

type getBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBook(ctx *gin.Context) {
	var req getBookRequest
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

	ctx.JSON(http.StatusOK, book)
}

type updateBookRequest struct {
	ID   int64 `uri:"id" binding:"required,min=1"`
	data struct {
		Name        string  `json:"name" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Image       string  `json:"image" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		Publisher   string  `json:"publisher" binding:"required"`
		Quantity    int32   `json:"quantity" binding:"required"`
	}
}

func (server *Server) updateBook(ctx *gin.Context) {
	var req updateBookRequest
	if err := ctx.ShouldBindUri(&req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.data); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateBookParams{
		ID:          req.ID,
		Name:        req.data.Name,
		Price:       req.data.Price,
		Image:       req.data.Image,
		Description: req.data.Description,
		Author:      req.data.Author,
		Publisher:   req.data.Publisher,
		Quantity:    req.data.Quantity,
	}

	book, err := server.store.UpdateBook(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, book)
}

type deleteBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBook(ctx *gin.Context) {
	var req deleteBookRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteBook(ctx, req.ID)
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

type listBookRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listBook(ctx *gin.Context) {
	var req listBookRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListBooksParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	books, err := server.store.ListBooks(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (server *Server) listBookInCart(ctx *gin.Context) {
	bookIDs, err := server.listBookInCartByUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	books := []db.Book{}
	for _, book := range bookIDs {
		book, err := server.store.GetBook(ctx, book.BooksID)
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
		books = append(books, book)
	}

	ctx.JSON(http.StatusOK, books)
}

func (server *Server) listBookInWishlist(ctx *gin.Context) {
	bookIDs, err := server.listBookInWishlistByUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	books := []db.Book{}
	for _, book := range bookIDs {
		book, err := server.store.GetBook(ctx, book.BooksID)
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
		books = append(books, book)
	}

	ctx.JSON(http.StatusOK, books)
}
