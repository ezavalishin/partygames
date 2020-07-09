package handlers

import (
	"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CurrentUser(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {
		user := auth.ForContext(context.Request.Context())

		context.JSON(http.StatusOK, user)
	}
}
