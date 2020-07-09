package server

import (
	"context"
	"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/handlers"
	log "github.com/ezavalishin/partygames/internal/logger"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-gonic/gin"
)

var host, port string

func init() {
	host = utils.MustGet("SERVER_HOST")
	port = utils.MustGet("SERVER_PORT")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Run(orm *orm.ORM) {
	log.Info("GORM_CONNECTION_DSN: ", utils.MustGet("GORM_CONNECTION_DSN"))

	r := gin.Default()

	r.GET("/ping", handlers.Ping())

	authorized := r.Group("/vkma")

	authorized.Use(auth.Middleware(orm))
	{
		authorized.GET("/me", handlers.CurrentUser(orm))
	}

	log.Info("Running @ http://" + host + ":" + port)
	log.Info(r.Run(host + ":" + port))

}
