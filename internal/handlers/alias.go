package handlers

import (
	//"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AliasWords(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {
		//user := auth.ForContext(context.Request.Context())

		var words []*models.Word

		count := 100

		if context.Request.URL.Query().Get("count") != "" {
			count, _ = strconv.Atoi(context.Request.URL.Query().Get("count"))
		}

		orm.DB.Order("random()").Limit(count).Find(&words)

		context.JSON(http.StatusOK, words)
	}
}