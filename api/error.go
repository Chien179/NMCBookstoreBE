package api

import "github.com/gin-gonic/gin"

type errorCustom struct {
	Field   string
	Message gin.H
}
