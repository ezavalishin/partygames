package server

import (
	"context"
	"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/handlers"
	"github.com/ezavalishin/partygames/internal/handlers/admin"
	"github.com/ezavalishin/partygames/internal/handlers/games"
	log "github.com/ezavalishin/partygames/internal/logger"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
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

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Vk-Params"},
		AllowCredentials: false,
		AllowAllOrigins: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", handlers.Ping())

	authorized := r.Group("/vkma")

	authorized.Use(auth.Middleware(orm))
	{
		authorized.GET("/me", handlers.CurrentUser(orm))
		authorized.GET("/alias/words", games.AliasWords(orm))
	}

	adminized := r.Group("/admin")

	adminized.Use()
	{
		adminized.GET("/tags", admin.GetTags(orm))
		adminized.POST("/words", admin.CreateWords(orm))
	}

	log.Info("Running @ http://" + host + ":" + port)
	log.Info(r.Run(host + ":" + port))

}
