package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (server *Server) getGenre(ctx *gin.Context) {
	var req models.GetGenreRequest
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

	ctx.JSON(http.StatusOK, genre)
}

func (server *Server) createGenre(ctx *gin.Context) {
	var req models.CreateGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	genre, err := server.store.CreateGenre(ctx, req.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genre)
}

func (server *Server) updateGenre(ctx *gin.Context) {
	var req models.UpdateGenreRequest
	if err := ctx.ShouldBindJSON(&req.Name); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req.ID); err != nil {
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
		Name: req.Name,
	}

	upadteGenre, err := server.store.UpdateGenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, upadteGenre)
}

func (server *Server) deleteGenre(ctx *gin.Context) {
	var req models.DeleteGenreRequest
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

func (server *Server) softDeleteGenre(ctx *gin.Context) {
	var req models.DeleteGenreRequest
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

	subgenres, err := server.store.ListSubgenres(ctx, genre.ID)
	for _, subgenre := range subgenres {
		_, err := server.store.SoftDeleteSubgenre(ctx, subgenre.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, err = server.store.SoftDeleteGenre(ctx, genre.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Genre deleted successfully")
}

func (server *Server) listGenre(ctx *gin.Context) {
	genres, err := server.store.ListGenres(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, genres)
}
