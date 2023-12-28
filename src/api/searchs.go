package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/helper"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/Chien179/NMCBookstoreBE/src/pb"
	"github.com/gin-gonic/gin"
)

func (server *Server) elasticSearch(ctx *gin.Context) {
	var req models.SearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	query, err := helper.QueryElastic(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err := server.elastic.Search(
		server.elastic.Search.WithIndex("books"),
		server.elastic.Search.WithBody(strings.NewReader(query)),
		server.elastic.Search.WithFilterPath("hits"),
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var result models.SearchResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err = server.elastic.Search(
		server.elastic.Search.WithIndex("books"),
		server.elastic.Search.WithBody(strings.NewReader(query)),
		server.elastic.Search.WithFilterPath("aggregations"),
	)
	if err != nil { // Parse []byte to go struct pointer
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	body, err = io.ReadAll(res.Body)
	var aggs models.Aggs
	if err := json.Unmarshal(body, &aggs); err != nil { // Parse []byte to go struct pointer
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var books []models.BookResponse
	for _, inf := range result.Hits.Hits {
		source := inf.Source
		book := models.BookResponse{
			ID:          source.ID,
			Name:        source.Name,
			Price:       source.Price,
			Sale:        source.Sale,
			Image:       source.Image,
			Description: source.Description,
			Author:      source.Author,
			Publisher:   source.Publisher,
			Quantity:    source.Quantity,
			Rating:      source.Rating,
			IsDeleted:   source.IsDeleted,
		}

		books = append(books, book)
	}

	booksByte, err := json.Marshal(books)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := db.ListBooksRow{TotalPage: aggs.Aggregations.UniqueBooks.Value, Books: json.RawMessage(booksByte)}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) recommend(ctx *gin.Context) {
	var req models.RecommedRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bookRequest := pb.BookRequest{
		Name: req.Name,
		Size: req.Size,
	}

	results, err := getBookRCM(ctx, &bookRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (server *Server) justForYou(ctx *gin.Context) {
	var req models.JustForYouRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	book, err := server.store.GetBestBookByUser(ctx, req.UserName)
	if err != nil {
		b, err := server.store.GetBook(ctx, book.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		book = db.GetBestBookByUserRow{ID: b.ID, Name: b.Name, Rating: int32(b.Rating)}
	}

	bookRequest := pb.BookRequest{
		Name: book.Name,
		Size: 17,
	}

	results, err := getBookRCM(ctx, &bookRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func getBookRCM(ctx context.Context, bookRequest *pb.BookRequest) ([]models.BookResponse, error) {
	books, err := pb.GRPCGetRCM(ctx, bookRequest)
	if err != nil {
		return []models.BookResponse{}, err
	}

	results := []models.BookResponse{}
	for _, book := range books.GetBooks() {
		b := models.BookResponse{
			ID:          book.Id,
			Name:        book.Name,
			Author:      book.Author,
			Publisher:   book.Publisher,
			Price:       float64(book.Price),
			Description: book.Description,
			Rating:      float64(book.Rating),
			Image:       book.Image,
			Quantity:    book.Quantity,
		}

		results = append(results, b)
	}

	return results, nil
}
