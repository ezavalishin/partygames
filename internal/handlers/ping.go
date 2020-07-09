package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "OK")
	}
}
