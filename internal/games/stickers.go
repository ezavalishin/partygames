package games

import (
	"fmt"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
	"math/rand"
	"time"
)

var stickerGames map[string]*ActiveStickerGame

type AttachedGameUser struct {
	User  *models.User `json:"user"`
	Index *int         `json:"index"`
}

type StickerGameUser struct {
	User             *models.User      `json:"user"`
	AttachedGameUser *AttachedGameUser `json:"attachedGameUser"`
	Word             *string           `json:"word"`
	IsFinished       bool              `json:"isFinished"`
}

type ActiveStickerGame struct {
	Id         uuid.UUID          `json:"id"`
	Creator    *models.User       `json:"creator"`
	StartedAt  *time.Time         `json:"startedAt"`
	FinishedAt *time.Time         `json:"finishedAt"`
	GameUsers  []*StickerGameUser `json:"gameUsers"`
}

type SetWord struct {
	GameId string `json:"gameId"`
	Word   string `json:"word"`
}

func CreateStickerGame(authedUser *models.User) ActiveStickerGame {

	if stickerGames == nil {
		stickerGames = make(map[string]*ActiveStickerGame)
	}

	fmt.Println("CREAATE")
	id, _ := uuid.NewRandom()

	var users []*StickerGameUser

	activeGame := ActiveStickerGame{
		Id:        id,
		Creator:   authedUser,
		GameUsers: users,
	}
	fmt.Println("ACTIVE GAME")

	activeGame.AddPlayer(authedUser)

	stickerGames[id.String()] = &activeGame

	fmt.Println("STORED")

	return activeGame
}

func JoinStickerGame(authedUser *models.User, id string) ActiveStickerGame {

	activeGame := stickerGames[id]

	activeGame.AddPlayer(authedUser)

	return *activeGame
}

func StartStickerGame(authedUser *models.User, id string) ActiveStickerGame {

	activeGame := stickerGames[id]

	activeGame.StartTyping()

	return *activeGame
}

func RestartStickerGame(authedUser *models.User, id string) ActiveStickerGame {

	activeGame := stickerGames[id]

	activeGame.Clear()
	activeGame.StartTyping()

	return *activeGame
}

func SetWordInGame(authedUser *models.User, id string, word string) ActiveStickerGame {

	activeGame := stickerGames[id]

	fmt.Printf("%+v", stickerGames)

	fmt.Printf("%+v", activeGame)

	activeGame.SetWordForUser(authedUser, word)

	return *activeGame
}

func GotWordInGame(authedUser *models.User, id string, server *socketio.Server) ActiveStickerGame {

	activeGame := stickerGames[id]

	activeGame.GotWordForUser(authedUser, server)

	return *activeGame
}

func (g *ActiveStickerGame) AddPlayer(p *models.User) {

	for _, gameUser := range g.GameUsers {
		if gameUser.User.ID == p.ID {
			return
		}
	}

	gameUser := StickerGameUser{
		User:       p,
		Word:       nil,
		IsFinished: false,
	}

	g.GameUsers = append(g.GameUsers, &gameUser)
}

func (g *ActiveStickerGame) Clear() {

	g.StartedAt = nil
	g.FinishedAt = nil

	for _, gameUser := range g.GameUsers {
		gameUser.Word = nil
		gameUser.AttachedGameUser = nil
		gameUser.IsFinished = false
	}
}

func (g *ActiveStickerGame) StartTyping() {
	now := time.Now()
	g.StartedAt = &now

	len := len(g.GameUsers)

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len, func(i, j int) { g.GameUsers[i], g.GameUsers[j] = g.GameUsers[j], g.GameUsers[i] })

	for i, gameUser := range g.GameUsers {
		index := (i + 1) % len
		gs := g.GameUsers[index]
		attachedGameUser := AttachedGameUser{
			User:  gs.User,
			Index: &index,
		}
		gameUser.AttachedGameUser = &attachedGameUser
	}
}

func (g *ActiveStickerGame) SetWordForUser(u *models.User, word string) {

	var currentGameUser *StickerGameUser

	fmt.Println("BEFORE FOUND")
	fmt.Println("BEFORE FOUND2")

	fmt.Printf("%+v", g.GameUsers)

	for _, gameUser := range g.GameUsers {
		if gameUser.User.ID == u.ID {
			currentGameUser = gameUser
		}
	}

	fmt.Println("AFTER FOUND")

	if currentGameUser == nil {
		return
	}

	index := currentGameUser.AttachedGameUser.Index

	attachedUser := g.GameUsers[*index]

	attachedUser.Word = &word
}

func (g *ActiveStickerGame) GotWordForUser(u *models.User, server *socketio.Server) {

	var currentGameUser *StickerGameUser

	for _, gameUser := range g.GameUsers {
		if gameUser.User.ID == u.ID {
			currentGameUser = gameUser
		}
	}

	if currentGameUser == nil {
		return
	}

	currentGameUser.IsFinished = true

	g.CheckGameIsFinished(server)
}

func (g *ActiveStickerGame) CheckGameIsFinished(server *socketio.Server) {

	finished := true

	for _, gameUser := range g.GameUsers {
		finished = finished && gameUser.IsFinished
	}

	if finished {
		now := time.Now()
		g.FinishedAt = &now

		server.BroadcastToRoom("/", g.Id.String(), "game-updated", utils.WrapJSON(g))
	}
}
