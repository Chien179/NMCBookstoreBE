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

// @Summary      Create subgenre
// @Description  Use this API to create subgenre
// @Tags         Subgenres
// @Accept       json
// @Produce      json
// @Param        Request body createSubgenreRequest  true  "Create subgenre"
// @Success      200 {object} db.Subgenre
// @failure	 	 400
// @failure		 500
// @Router       /admin/subgenres [post]
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

type updateSubgenreData struct {
	GenreID int64  `json:"genre_id"`
	Name    string `json:"name"`
}

type updateSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	updateSubgenreData
}

// @Summary      Update subgenre
// @Description  Use this API to update subgenre
// @Tags         Subgenres
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "Update subgenre id"
// @Param        Request body updateSubgenreData  false  "Update subgenre data"
// @Success      200 {object} db.Subgenre
// @failure	 	 400
// @failure		 500
// @Router       /admin/subgenres/update/{id} [put]
func (server *Server) updateSubgenre(ctx *gin.Context) {
	var req updateSubgenreRequest
	if err := ctx.ShouldBindJSON(&req.updateSubgenreData); err != nil {
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
			Int64: req.updateSubgenreData.GenreID,
			Valid: req.updateSubgenreData.GenreID > 0,
		},
		Name: sql.NullString{
			String: req.updateSubgenreData.Name,
			Valid:  req.updateSubgenreData.Name != "",
		},
	}

	updateSubgenre, err := server.store.UpdateSubgenre(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updateSubgenre)
}

type deleteSubgenreRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary      Delete subgenre
// @Description  Use this API to delete subgenre
// @Tags         Subgenres
// @Accept       json
// @Produce      json
// @Param        ID path int  true  "delete subgenre"
// @Success      200
// @failure	 	 400
// @failure		 500
// @Router       /admin/subgenres/delete/{id} [delete]
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

// @Summary      List subgenre
// @Description  Use this API to list subgenre
// @Tags         Subgenres
// @Accept       json
// @Produce      json
// @Param        genre_id path int  true  "list subgenre"
// @Success      200 {object} []db.Subgenre
// @failure	 	 400
// @failure	 	 404
// @failure		 500
// @Router       /subgenres/{genre_id} [get]
func (server *Server) listSubgenre(ctx *gin.Context) {
	var req listSubgenreRequest
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
