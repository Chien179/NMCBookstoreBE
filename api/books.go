package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type bookResponse struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	Image       []string      `json:"image"`
	Description string        `json:"description"`
	Author      string        `json:"author"`
	Publisher   string        `json:"publisher"`
	Quantity    int32         `json:"quantity"`
	Genres      []db.Genre    `json:"genres"`
	Subgenres   []db.Subgenre `json:"subgenres"`
}

func (server *Server) createBook(ctx *gin.Context) {
	req, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var imgUrls []string

	for i, img := range req.File["image"] {
		imgUrl, err := server.uploadFile(ctx, img, "NMCBookstore/Image/Books/"+req.Value["name"][0], req.Value["name"][0]+" "+strconv.Itoa(int(i)))
		if err != nil {
			return
		} else {
			imgUrls = append(imgUrls, imgUrl)
		}
	}

	price, err := strconv.Atoi(req.Value["price"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quantity, err := strconv.Atoi(req.Value["quantity"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateBookParams{
		Name:        req.Value["name"][0],
		Price:       float64(price),
		Image:       imgUrls,
		Description: req.Value["description"][0],
		Author:      req.Value["author"][0],
		Publisher:   req.Value["publisher"][0],
		Quantity:    int32(quantity),
	}

	book, err := server.store.CreateBook(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	genres := []db.Genre{}
	for _, ID := range req.Value["genres_id"] {
		id, err := strconv.Atoi(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argBG := db.CreateBookGenreParams{
			BooksID:  book.ID,
			GenresID: int64(id),
		}

		_, err = server.store.CreateBookGenre(ctx, argBG)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		genre, err := server.store.GetGenre(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		genres = append(genres, genre)
	}

	subgenres := []db.Subgenre{}
	for _, ID := range req.Value["subgenres_id"] {
		id, err := strconv.Atoi(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argBS := db.CreateBookSubgenreParams{
			BooksID:     book.ID,
			SubgenresID: int64(id),
		}

		_, err = server.store.CreateBookSubgenre(ctx, argBS)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		subgenre, err := server.store.GetSubgenre(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		subgenres = append(subgenres, subgenre)
	}

	rsp := bookResponse{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Image:       book.Image,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		Quantity:    book.Quantity,
		Genres:      genres,
		Subgenres:   subgenres,
	}

	ctx.JSON(http.StatusOK, rsp)
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

	genresID, err := server.store.ListBooksGenresIDByBookID(ctx, book.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	subgenresID, err := server.store.ListBooksSubgenresIDByBookID(ctx, book.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	genres := []db.Genre{}
	for _, ID := range genresID {
		genre, err := server.store.GetGenre(ctx, ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		genres = append(genres, genre)
	}

	subgenres := []db.Subgenre{}
	for _, ID := range subgenresID {
		subgenre, err := server.store.GetSubgenre(ctx, ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		subgenres = append(subgenres, subgenre)
	}

	rsp := bookResponse{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Image:       book.Image,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		Quantity:    book.Quantity,
		Genres:      genres,
		Subgenres:   subgenres,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) updateBook(ctx *gin.Context) {
	req, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(req.Value["id"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	book, err := server.store.GetBook(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	imgUrls := req.Value["image"]

	if req.File["files"] != nil {
		for i, img := range req.File["files"] {
			imgUrl, err := server.uploadFile(ctx, img, "NMCBookstore/Image/Books/"+req.Value["name"][0], req.Value["name"][0]+" "+strconv.Itoa(int(i)))
			if err != nil {
				return
			} else {
				imgUrls = append(imgUrls, imgUrl)
			}
		}
	}

	price, err := strconv.Atoi(req.Value["price"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	quantity, err := strconv.Atoi(req.Value["quantity"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateBookParams{
		ID: book.ID,
		Name: sql.NullString{
			String: req.Value["name"][0],
			Valid:  req.Value["name"][0] != "",
		},
		Price: sql.NullFloat64{
			Float64: float64(price),
			Valid:   price > 0,
		},
		Image: imgUrls,
		Description: sql.NullString{
			String: req.Value["description"][0],
			Valid:  req.Value["description"][0] != "",
		},
		Author: sql.NullString{
			String: req.Value["author"][0],
			Valid:  req.Value["author"][0] != "",
		},
		Publisher: sql.NullString{
			String: req.Value["publisher"][0],
			Valid:  req.Value["publisher"][0] != "",
		},
		Quantity: sql.NullInt32{
			Int32: int32(quantity),
			Valid: quantity > 0,
		},
	}

	updatedBook, err := server.store.UpdateBook(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBookGenreByBooksID(ctx, book.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBookSubgenreByBooksID(ctx, book.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	genres := []db.Genre{}
	for _, ID := range req.Value["genres_id"] {
		id, err := strconv.Atoi(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argBG := db.CreateBookGenreParams{
			BooksID:  book.ID,
			GenresID: int64(id),
		}

		_, err = server.store.CreateBookGenre(ctx, argBG)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		genre, err := server.store.GetGenre(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		genres = append(genres, genre)
	}

	subgenres := []db.Subgenre{}
	for _, ID := range req.Value["subgenres_id"] {
		id, err := strconv.Atoi(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argBS := db.CreateBookSubgenreParams{
			BooksID:     book.ID,
			SubgenresID: int64(id),
		}

		_, err = server.store.CreateBookSubgenre(ctx, argBS)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		subgenre, err := server.store.GetSubgenre(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		subgenres = append(subgenres, subgenre)
	}

	rsp := bookResponse{
		ID:          updatedBook.ID,
		Name:        updatedBook.Name,
		Price:       updatedBook.Price,
		Image:       updatedBook.Image,
		Description: updatedBook.Description,
		Author:      updatedBook.Author,
		Publisher:   updatedBook.Publisher,
		Quantity:    updatedBook.Quantity,
		Genres:      genres,
		Subgenres:   subgenres,
	}

	ctx.JSON(http.StatusOK, rsp)
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
