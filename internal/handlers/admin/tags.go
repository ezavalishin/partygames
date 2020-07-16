package admin

import (
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTags(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {

		var tags []*models.Tag

		orm.DB.Find(&tags)

		context.JSON(http.StatusOK, utils.WrapJSON(tags))
	}
}