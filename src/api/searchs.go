package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) fullSearch(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) recommend(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) justForYou(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
