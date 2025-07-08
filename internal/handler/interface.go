package handler

import "github.com/gin-gonic/gin"

type URLHandler interface {
	RegisterRoutes(r *gin.Engine)
}
