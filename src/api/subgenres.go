package api

import (
	"database/sql"
	"net/http"

	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (server *Server) createSubgenre(ctx *gin.Context) {
	var req models.CreateSubgenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateSubgenreParams{
		GenresID: req.GenresID,
		Name:     req.Name,
	}

	subgenre, err := server.store.CreateSubgenre(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenre)
}

func (server *Server) getSubgenre(ctx *gin.Context) {
	var req models.GetSubgenreRequest
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

func (server *Server) updateSubgenre(ctx *gin.Context) {
	var req models.UpdateSubgenreRequest
	if err := ctx.ShouldBindJSON(&req.UpdateSubgenreData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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

	arg := db.UpdateSubgenreParams{
		ID: subgenre.ID,
		GenresID: sql.NullInt64{
			Int64: req.GenreID,
			Valid: req.GenreID > 0,
		},
		Name: sql.NullString{
			String: req.Name,
			Valid:  req.Name != "",
		},
	}

	updateSubgenre, err := server.store.UpdateSubgenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updateSubgenre)
}

func (server *Server) deleteSubgenre(ctx *gin.Context) {
	var req models.DeleteSubgenreRequest
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

func (server *Server) softDeleteSubgenre(ctx *gin.Context) {
	var req models.DeleteSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.SoftDeleteSubgenre(ctx, req.ID)
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

func (server *Server) listSubgenre(ctx *gin.Context) {
	var req models.ListSubgenreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subgenres, err := server.store.ListSubgenres(ctx, req.GenreID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenres)
}

func (server *Server) listAllSubgenre(ctx *gin.Context) {
	subgenres, err := server.store.ListAllSubgenres(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subgenres)
}

func (server *Server) listSubgenresNoticeable(ctx *gin.Context) {
	fiveSubgenres, err := server.store.ListSubgenresNoticeable(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, fiveSubgenres)
}
