package handlers

import (
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRandomStickerWord(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {

		var tag models.Tag
		var word models.Word

		orm.DB.Where("value = ?", "stickers").Find(&tag)

		orm.DB.Order("random()").Model(&tag).Related(&word, "Words")

		context.JSON(http.StatusOK, utils.WrapJSON(word))
	}
}
