package helper

import "github.com/gin-gonic/gin"

type ErrorCustom struct {
	Field   string
	Message gin.H
}
