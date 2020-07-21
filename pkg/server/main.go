package server

import (
	"encoding/json"
	"fmt"
	"github.com/ezavalishin/partygames/internal/auth"
	"github.com/ezavalishin/partygames/internal/games"
	"github.com/ezavalishin/partygames/internal/handlers"
	"github.com/ezavalishin/partygames/internal/handlers/admin"
	log "github.com/ezavalishin/partygames/internal/logger"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"time"
)

var host, port string

func init() {
	host = utils.MustGet("SERVER_HOST")
	port = utils.MustGet("SERVER_PORT")
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func Run(orm *orm.ORM) {
	log.Info("GORM_CONNECTION_DSN: ", utils.MustGet("GORM_CONNECTION_DSN"))

	r := gin.Default()

	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		vkParams := url.Query().Get("vk-params")
		fmt.Println("connected:", s.ID())

		err := auth.WsValidateAndSetUser(orm, s, vkParams)

		user := auth.ForWssContext(s)

		if user != nil {
			go games.SetUserIsOnline(user, server)

			fmt.Println("CONNECTED " + *user.FirstName + ", " + s.ID())
		}

		if err != nil {
			return err
		}

		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		user := auth.ForWssContext(s)

		if user != nil {
			go games.SetUserIsOffline(user, server)

			fmt.Println("DISCONNECTED " + *user.FirstName + ", " + s.ID())
		}
	})

	server.OnEvent("/", "kek", func(s socketio.Conn, msg string) error {

		user := auth.ForWssContext(s)

		s.Emit("reply", utils.WrapJSON(user))
		return nil
	})

	server.OnEvent("/", "create-game", func(s socketio.Conn, msg string) error {

		fmt.Println("CREAT EVENT")

		user := auth.ForWssContext(s)

		activeGame := games.CreateStickerGame(user)

		fmt.Println("ACTIVE GAME")
		fmt.Printf("%+v", activeGame)

		s.Join(activeGame.Id.String())

		s.Emit("game-created", utils.WrapJSON(activeGame))
		return nil
	})

	server.OnEvent("/", "join-game", func(s socketio.Conn, msg string) error {

		fmt.Println("JOIN")

		user := auth.ForWssContext(s)

		activeGame := games.JoinStickerGame(user, msg)

		fmt.Printf("%+v", activeGame)

		s.Join(activeGame.Id.String())

		fmt.Println("JOINED")

		server.BroadcastToRoom("/", activeGame.Id.String(), "game-updated", utils.WrapJSON(activeGame))

		fmt.Println("SENT")

		return nil
	})

	server.OnEvent("/", "start-game-prepare", func(s socketio.Conn, msg string) error {

		fmt.Println("JOIN")

		user := auth.ForWssContext(s)

		activeGame := games.StartStickerGame(user, msg)

		server.BroadcastToRoom("/", activeGame.Id.String(), "game-prepared", utils.WrapJSON(activeGame))

		fmt.Println("SENT")

		return nil
	})

	server.OnEvent("/", "restart-game", func(s socketio.Conn, msg string) error {

		fmt.Println("RESTART")

		user := auth.ForWssContext(s)

		activeGame := games.RestartStickerGame(user, msg)

		server.BroadcastToRoom("/", activeGame.Id.String(), "game-restarted", utils.WrapJSON(activeGame))

		fmt.Println("SENT")

		return nil
	})

	server.OnEvent("/", "set-word", func(s socketio.Conn, msg string) error {

		fmt.Println("SET WORD")

		user := auth.ForWssContext(s)

		setWord := games.SetWord{}

		json.Unmarshal([]byte(msg), &setWord)

		fmt.Printf("%+v", setWord)

		activeGame := games.SetWordInGame(user, setWord.GameId, setWord.Word)

		server.BroadcastToRoom("/", activeGame.Id.String(), "word-set", utils.WrapJSON(activeGame))

		fmt.Println("SENT")

		return nil
	})

	server.OnEvent("/", "got-word", func(s socketio.Conn, msg string) error {

		fmt.Println("GOT WORD")

		user := auth.ForWssContext(s)

		activeGame := games.GotWordInGame(user, msg, server)

		server.BroadcastToRoom("/", activeGame.Id.String(), "game-updated", utils.WrapJSON(activeGame))

		fmt.Println("SENT")

		return nil
	})

	server.OnEvent("/", "leave", func(s socketio.Conn, msg string) error {

		fmt.Println("LEAVE")

		user := auth.ForWssContext(s)

		s.Leave(msg)

		games.LeaveFromStickerGame(user, msg, server)

		return nil
	})

	go server.Serve()
	defer server.Close()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "WS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Vk-Params", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", handlers.Ping())

	authorized := r.Group("/vkma")

	authorized.Use(auth.Middleware(orm))
	{
		authorized.GET("/me", handlers.CurrentUser(orm))
		authorized.GET("/alias/words", handlers.AliasWords(orm))
		authorized.GET("/stickers/word", handlers.GetRandomStickerWord(orm))
	}

	adminized := r.Group("/admin")

	adminized.Use()
	{
		adminized.GET("/tags", admin.GetTags(orm))
		adminized.POST("/words", admin.CreateWords(orm))
	}

	ws := r.Group("/socket.io")

	ws.Use(auth.Middleware(orm))
	ws.Use(GinMiddleware("*"))
	{
		ws.Handle("POST", "/*any", gin.WrapH(server))
		ws.Handle("GET", "/*any", gin.WrapH(server))
	}

	//r.Handle ( "WS", "/socket.io/*any", gin.WrapH(server) )
	//r.Handle ( "WSS", "/socket.io/*any", gin.WrapH(server) )

	log.Info("Running @ http://" + host + ":" + port)
	log.Info(r.Run(host + ":" + port))

}
