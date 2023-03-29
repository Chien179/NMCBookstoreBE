package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

// @Summary      Create genre
// @Description  Use this API to create genre
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Param        Request body createGenreRequest  true  "Create genre"
// @Success      200 {object} db.Genre
// @failure	 	 400
// @failure		 500
// @Router       /admin/genres [post]
func (server *Server) createGenre(ctx *gin.Context) {
	var req createGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	genre, err := server.store.CreateGenre(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genre)
}

type updateGenreData struct {
	Name string `json:"name" binding:"required"`
}

type updateGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	updateGenreData
}

// @Summary      Update genre
// @Description  Use this API to update genre
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Param        ID path int true  "Update genre id"
// @Param        Request body updateGenreData true  "Update genre data"
// @Success      200 {object} db.Genre
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /admin/genres/update/{id} [put]
func (server *Server) updateGenre(ctx *gin.Context) {
	var req updateGenreRequest
	if err := ctx.ShouldBindJSON(&req.updateGenreData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	genre, err := server.store.GetGenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateGenreParams{
		ID:   genre.ID,
		Name: req.updateGenreData.Name,
	}

	upadteGenre, err := server.store.UpdateGenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, upadteGenre)
}

type deleteGenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete genre
// @Description  Use this API to delete genre
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Param        ID path int true  "Delete genre"
// @Success      200
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /admin/genres/delete/{id} [delete]
func (server *Server) deleteGenre(ctx *gin.Context) {
	var req deleteGenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	genre, err := server.store.GetGenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteGenre(ctx, genre.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Genre deleted successfully")
}

// @Summary      List genre
// @Description  Use this API to list genre
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Success      200 {object} []db.Genre
// @failure		 500
// @Router       /genres [get]
func (server *Server) listGenre(ctx *gin.Context) {
	genres, err := server.store.ListGenres(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genres)
}
