package api

import (
	"database/sql"
	"mime/multipart"
	"net/http"
	"strconv"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBookRequest struct {
	Name        string                 `form:"name" binding:"required"`
	Price       float64                `form:"price" binding:"required"`
	Image       []multipart.FileHeader `form:"image" binding:"required"`
	Description string                 `form:"description" binding:"required"`
	Author      string                 `form:"author" binding:"required"`
	Publisher   string                 `form:"publisher" binding:"required"`
	Quantity    int32                  `form:"quantity" binding:"required"`
	GenresID    int64                  `form:"genres_id" binding:"required"`
	SubgenresID int64                  `form:"subgenres_id" binding:"required"`
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
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var imgUrls []string

	for i, img := range req.Image {
		imgUrl, err := server.uploadFile(ctx, &img, "NMCBookstore/Image/Books/"+req.Name, req.Name+" "+strconv.Itoa(int(i)))
		if err != nil {
			return
		} else {
			imgUrls = append(imgUrls, imgUrl)
		}
	}

	arg := db.CreateBookParams{
		Name:        req.Name,
		Price:       req.Price,
		Image:       imgUrls,
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

	argBG := db.CreateBookGenreParams{
		BooksID:  book.ID,
		GenresID: req.GenresID,
	}

	_, err = server.store.CreateBookGenre(ctx, argBG)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argBS := db.CreateBookSubgenreParams{
		BooksID:     book.ID,
		SubgenresID: req.SubgenresID,
	}

	_, err = server.store.CreateBookSubgenre(ctx, argBS)
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
	Name        string                 `form:"name"`
	Price       float64                `form:"price"`
	Image       []multipart.FileHeader `form:"image" binding:"required"`
	Description string                 `form:"description"`
	Author      string                 `form:"author"`
	Publisher   string                 `form:"publisher"`
	Quantity    int32                  `form:"quantity"`
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
// @Param        Request body updateBookData  false  "Update book request"
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

	if err := ctx.ShouldBind(&req.updateBookData); err != nil {
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

	var imgUrls []string

	for i, img := range req.Image {
		imgUrl, err := server.uploadFile(ctx, &img, "NMCBookstore/Image/Books/"+req.Name, req.Name+" "+strconv.Itoa(int(i)))
		if err != nil {
			return
		} else {
			imgUrls = append(imgUrls, imgUrl)
		}
	}

	arg := db.UpdateBookParams{
		ID: book.ID,
		Name: sql.NullString{
			String: req.updateBookData.Name,
			Valid:  req.updateBookData.Name != "",
		},
		Price: sql.NullFloat64{
			Float64: req.updateBookData.Price,
			Valid:   req.updateBookData.Price > -1,
		},
		Image: imgUrls,
		Description: sql.NullString{
			String: req.updateBookData.Description,
			Valid:  req.updateBookData.Description != "",
		},
		Author: sql.NullString{
			String: req.updateBookData.Author,
			Valid:  req.updateBookData.Author != "",
		},
		Publisher: sql.NullString{
			String: req.updateBookData.Publisher,
			Valid:  req.updateBookData.Publisher != "",
		},
		Quantity: sql.NullInt32{
			Int32: req.updateBookData.Quantity,
			Valid: req.updateBookData.Quantity > 0,
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

func (server *Server) listAllBook(ctx *gin.Context) {
	books, err := server.store.ListAllBooks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
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
		if books.Books == nil {
			ctx.JSON(http.StatusOK, books)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (server *Server) listTheBestBook(ctx *gin.Context) {
	books, err := server.store.ListTheBestBooks(ctx)
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

func (server *Server) listNewestBook(ctx *gin.Context) {
	books, err := server.store.ListNewestBooks(ctx)
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
