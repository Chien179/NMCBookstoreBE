package api

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) uploadFile(ctx *gin.Context, fileHeader *multipart.FileHeader, filePath string, fileName string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return "", err
	}

	fileUrl, err := server.uploader.FileUpload(file, filePath, fileName)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return "", err
	}

	return fileUrl, nil
}
