package admin

import (
	"fmt"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type blockCollection struct {
	Blocks []block
}

type block struct {
	Words []string
	Tags []string
}

func CreateWords(orm *orm.ORM) gin.HandlerFunc {
	return func(context *gin.Context) {

		input := blockCollection{}

		err := context.BindJSON(&input)

		if err != nil {
			http.Error(context.Writer, "", http.StatusUnprocessableEntity)
			context.Abort()
			return
		}

		fmt.Printf("%+v", input)

		var tags []*models.Tag

		orm.DB.Find(&tags)

		context.JSON(http.StatusOK, utils.WrapJSON(tags))
	}
}
