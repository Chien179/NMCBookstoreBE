package api

import (
	"net/http"

	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) getRank(ctx *gin.Context) {
	var req models.RankRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rank, err := server.store.GetRank(ctx, user.Rank)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	vote, err := server.store.GetCountLikeByUser(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	review, err := server.store.GetCountReviewByUser(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := models.RankReponse{
		Rank:    rank.Name,
		Vote:    int(vote),
		Reviews: int(review),
	}

	ctx.JSON(http.StatusOK, rsp)

}
