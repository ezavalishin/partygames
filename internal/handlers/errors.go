package handlers

import (
	"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorLog struct {
}

func StoreErrorLog(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {

		errorLog := models.ErrorLog{}

		user := auth.ForContext(context.Request.Context())

		err := context.BindJSON(&errorLog)

		if err != nil {
			http.Error(context.Writer, "", http.StatusUnprocessableEntity)
			context.Abort()
			return
		}

		ua := context.Request.UserAgent()
		errorLog.Ua = &ua
		errorLog.UserId = user.ID

		orm.DB.Create(&errorLog)

		context.JSON(http.StatusOK, nil)
	}
}
