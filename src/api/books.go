package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) createBook(ctx *gin.Context) {
	req, err := ctx.MultipartForm()
	price, err := strconv.Atoi(req.Value["price"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	sale, err := strconv.Atoi(req.Value["sale"][0])
	if err != nil {
		sale = 0
	}
	quantity, err := strconv.Atoi(req.Value["quantity"][0])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if price < 0 {
		ctx.JSON(http.StatusBadRequest, "Book price must be positive")
		return
	}

	if quantity < 1 {
		ctx.JSON(http.StatusBadRequest, "Book quantity price must be greater than 0")
		return
	}

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

	arg := db.CreateBookParams{
		Name:        req.Value["name"][0],
		Price:       float64(price),
		Sale:        float64(sale),
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

	rsp := models.BookResponse{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Sale:        float64(book.Sale),
		Image:       book.Image,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		Quantity:    book.Quantity,
		Genres:      genres,
		Subgenres:   subgenres,
		Rating:      0,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getBook(ctx *gin.Context) {
	var req models.GetBookRequest
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

	rsp := models.BookResponse{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Sale:        float64(book.Sale),
		Image:       book.Image,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		Quantity:    book.Quantity,
		Genres:      genres,
		Subgenres:   subgenres,
		Rating:      book.Rating,
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
	sale, err := strconv.Atoi(req.Value["sale"][0])
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
		Sale: sql.NullFloat64{
			Float64: book.Sale,
			Valid:   sale >= 0,
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

	rsp := models.BookResponse{
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
		Rating:      updatedBook.Rating,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteBook(ctx *gin.Context) {
	var req models.DeleteBookRequest
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

func (server *Server) listBook(ctx *gin.Context) {
	var req models.ListBookRequest
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

func (server *Server) softDeleteBook(ctx *gin.Context) {
	var req models.DeleteBookRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.SoftDeleteBook(ctx, req.ID)
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
