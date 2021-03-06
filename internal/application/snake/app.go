package snake

import (
	"Wave/internal/application/proto"
	"Wave/internal/application/snake/core"
	"time"
)

//go:generate easyjson .

type App struct {
	*proto.Room // base class
	game        *game
}

const RoomType proto.RoomType = "snake"

// ----------------|

// New snake app
func New(id proto.RoomToken, m proto.IManager, db interface{}, step time.Duration) proto.IRoom {
	s := &App{
		Room: proto.NewRoom(id, RoomType, m, step),
		game: newGame(core.Vec2i{
			X: 60,
			Y: 40,
		}),
	}
	s.SetCounterType(proto.FillGaps)
	s.OnTick = s.onTick
	s.OnUserRemove = s.onUserRemoved
	s.game.OnSnakeDead = s.onSnakeDead
	s.Routes["game_action"] = s.onGameAction
	s.Routes["game_play"] = s.onGamePlay
	s.Routes["game_exit"] = s.onGameExit
	return s
}

// ----------------| handlers

// <- STATUS_TICK
func (a *App) onTick(dt time.Duration) {
	a.game.Tick(dt)
	info := a.game.GetGameInfo()
	for i, s := range info.Snakes {
		serial, _ := a.GetTokenSerial(s.UserToken)
		info.Snakes[i].Serial = serial
	}
	a.Broadcast(messageTick.WithStruct(info))
}

func (a *App) onUserRemoved(u proto.IUser) {
	a.game.DeleteSnake(u)
}

// -> game_action
func (a *App) onGameAction(u proto.IUser, im proto.IInMessage) {
	ac := &gameAction{}
	if im.ToStruct(ac) != nil {
		return
	}

	switch ac.ActionName {
	case "move_left":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Left) })
	case "move_right":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Right) })
	case "move_up":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Up) })
	case "move_down":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Down) })
	}
}

// -> game_play
func (a *App) onGamePlay(u proto.IUser, im proto.IInMessage) {
	a.game.CreateSnake(u, 3)
}

// -> game_exit
func (a *App) onGameExit(u proto.IUser, im proto.IInMessage) {
	a.game.DeleteSnake(u)

	if len(a.game.user2snake) == 0 {
		a.exit()
	}
}

// <- STATUS_DEAD | win
func (a *App) onSnakeDead(u proto.IUser) {
	serial, _ := a.GetUserSerial(u)
	a.Broadcast(messageDead.WithStruct(&playerPayload{
		UserToken:  u.GetToken(),
		UserName:   u.GetName(),
		UserSerial: serial,
	}))

	if len(a.game.user2snake) <= 1 {
		// send a victory message
		if len(a.game.user2snake) == 1 {
			var w proto.IUser
			for w = range a.game.user2snake {
				// empty
			}
			serial, _ := a.GetUserSerial(w)
			a.Broadcast(messageWin.WithStruct(&playerPayload{
				UserToken:  w.GetToken(),
				UserName:   w.GetName(),
				UserSerial: serial,
			}))
		}
		// stop the room
		a.exit()
	}
}

func (a *App) exit() {
	m := a.GetManager()
	m.Task(a, func() {
		m.RemoveLobby(a.GetToken(), nil)
	})
}

// ----------------| helpers

// easyjson:json
type gameAction struct {
	ActionName string `json:"action"`
}

type playerPayload struct {
	UserName   string          `json:"user_name"`
	UserToken  proto.UserToken `json:"user_token"`
	UserSerial int64           `json:"user_serial"`
}

var (
	messageWin            = proto.Response{Status: "win"}.WithStruct("")
	messageDead           = proto.Response{Status: "STATUS_DEAD"}.WithStruct("")
	messageNoSnake        = proto.Response{Status: "STATUS_ERROR"}.WithReason("No snake")
	messageAlreadyPlays   = proto.Response{Status: "STATUS_ERROR"}.WithReason("already plays")
	messageUnknownCommand = proto.Response{Status: "STATUS_ERROR"}.WithReason("unknown command")
	messageTick           = proto.Response{Status: "STATUS_TICK"}.WithReason("")
)

func (a *App) withSnake(u proto.IUser, next func(s *snake)) {
	if s, err := a.game.GetSnake(u); err == nil {
		next(s)
	}
}
