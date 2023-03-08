package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createSubgenreRequest struct {
	GenreID int64  `json:"genre_id" binding:"required,min=1"`
	Name    string `json:"name" binding:"required"`
}

func (server *Server) createSubgenre(ctx *gin.Context) {
	var req createSubgenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateSubgenreParams{
		GenresID: req.GenreID,
		Name:     req.Name,
	}

	subgenre, err := server.store.CreateSubgenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenre)
}

type getSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSubgenre(ctx *gin.Context) {
	var req getSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subgenre, err := server.store.GetSubgenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenre)
}

type updateSubgenreRequest struct {
	ID      int64  `uri:"id" binding:"required,min=1"`
	GenreID int64  `json:"genre_id" binding:"required,min=1"`
	Name    string `json:"name" binding:"required"`
}

func (server *Server) updateSubgenre(ctx *gin.Context) {
	var req updateSubgenreRequest
	if err := ctx.ShouldBindJSON(&req.Name); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateSubgenreParams{
		ID:       req.ID,
		GenresID: req.GenreID,
		Name:     req.Name,
	}

	subgenre, err := server.store.UpdateSubgenre(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenre)
}

type deleteSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteSubgenre(ctx *gin.Context) {
	var req deleteSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteSubgenre(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Subgenre deleted successfully")
}

type listSubgenreRequest struct {
	GenreID int64 `uri:"genre_id" binding:"required,min=1"`
}

func (server *Server) listSubgenre(ctx *gin.Context) {
	var req listSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subgenres, err := server.store.ListSubgenres(ctx, req.GenreID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenres)
}
