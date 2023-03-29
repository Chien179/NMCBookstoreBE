package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBookRequest struct {
	Name        string   `json:"name" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Image       []string `json:"image" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Author      string   `json:"author" binding:"required"`
	Publisher   string   `json:"publisher" binding:"required"`
	Quantity    int32    `json:"quanlity" binding:"required"`
}

// @Summary      Create book
// @Description  Use this API to create book
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Request body createBookRequest  true  "Create book"
// @Success      200  {object}  db.Book
// @failure	 	 400
// @failure		 500
// @Router       /admin/books [post]
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

// @Summary      Get book
// @Description  Use this API to get book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        id path int  true  "get book"
// @Success      200  {object}  db.Book
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /books/{id} [get]
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

type updateBookData struct {
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Image       []string `json:"image"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Publisher   string   `json:"publisher"`
	Quantity    int32    `json:"quantity"`
}

type updateBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	updateBookData
}

// @Summary      Update book
// @Description  Use this API to update book
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Update book id"
// @Param        request body updateBookData  false  "Update book request"
// @Success      200  {object}  db.Book
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /admin/books/update/{id} [put]
func (server *Server) updateBook(ctx *gin.Context) {
	var req updateBookRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.updateBookData); err != nil {
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

	arg := db.UpdateBookParams{
		ID: book.ID,
		Name: sql.NullString{
			String: req.updateBookData.Name,
			Valid:  true,
		},
		Price: sql.NullFloat64{
			Float64: req.updateBookData.Price,
			Valid:   true,
		},
		Image: req.updateBookData.Image,
		Description: sql.NullString{
			String: req.updateBookData.Description,
			Valid:  true,
		},
		Author: sql.NullString{
			String: req.updateBookData.Author,
			Valid:  true,
		},
		Publisher: sql.NullString{
			String: req.updateBookData.Publisher,
			Valid:  true,
		},
		Quantity: sql.NullInt32{
			Int32: req.updateBookData.Quantity,
			Valid: true,
		},
	}

	updatedBook, err := server.store.UpdateBook(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

type deleteBookRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete book
// @Description  Use this API to Delete book
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id path int  true  "Delete book"
// @Success      200
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /admin/books/delete/{id} [delete]
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
	PageSize int32 `form:"page_size" binding:"required,min=24,max=100"`
}

// @Summary      List book
// @Description  Use this API to List book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        query query listBookRequest  true  "List book"
// @Success      200 {object} []db.Book
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /books [get]
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
