package room

import (
	lg "Wave/internal/logger"
	json "encoding/json"

	"github.com/gorilla/websocket"
)

type User struct {
	ID    UserID
	Name  string
	Rooms map[RoomToken]IRoom
	Conn  *websocket.Conn
	LG    *lg.Logger

	input   chan IInMessage
	output  chan IOutMessage
	cancel  chan interface{}
	task    chan func()
	bClosed bool
}

func NewUser(ID UserID, Conn *websocket.Conn) *User {
	return &User{
		ID:     ID,
		Conn:   Conn,
		cancel: make(chan interface{}, 1),
		output: make(chan IOutMessage, 1000),
		input:  make(chan IInMessage, 1000),
		task:   make(chan func(), 1000),
		Rooms:  map[RoomToken]IRoom{},
	}
}

// ----------------| IUser interface

func (u *User) GetName() string { return u.Name }
func (u *User) GetID() UserID   { return u.ID }

func (u *User) AddToRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	r.Task(func() { r.AddUser(u) })
	u.Rooms[r.GetID()] = r
	return nil
}

func (u *User) RemoveFromRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	r.Task(func() { r.RemoveUser(u) })
	delete(u.Rooms, r.GetID())
	return nil
}

func (u *User) Listen() error {
	u.LG.Sugar.Infof("User started: id=%s", u.GetID())
	defer func() {
		u.LG.Sugar.Infof("User stopped: id=%s", u.GetID())
	}()

	go u.sendWorker()
	go u.receiveWorker()

	// send current user_id
	u.Consume(&OutMessage{
		Status: "STATUS_TOKEN",
		Payload: &userTokenPayload{
			UserToken: u.GetID(),
		},
	})

	for { // stops when connection closes
		select {
		case m := <-u.input:
			// log input
			if u.LG != nil {
				data, _ := json.Marshal(m)
				u.LG.Sugar.Infof("in_message: %v %v", u.GetID(), string(data))
			}

			// apply the message to a room
			if r, ok := u.getRoom(m.GetRoomID()); ok {
				r.Task(func() { r.ApplyMessage(u, m) })
			} else {
				u.Consume(&OutMessage{
					RoomToken: m.GetRoomID(),
					Status:    StatusError,
					Payload:   "Unknown room:" + m.GetRoomID(),
				})
				continue
			}
			u.LG.Sugar.Infof("loop: %v", u.GetID())
		case t := <-u.task:
			t()
		case <-u.cancel:
			return nil
		}
	}
}

func (u *User) StopListening() error {
	if !u.bClosed {
		u.removeFromAllRooms() // do we need this?
		u.stop()
	}
	return nil
}

func (u *User) Consume(m IOutMessage) error {
	if m == nil {
		return ErrorNil
	}
	u.output <- m
	return nil
}

func (u *User) Task(t func()) {
	if t != nil {
		u.task <- t
	}
}

// ----------------| internal function

// easyjson:json
type userTokenPayload struct {
	UserToken UserID `json:"user_token"`
}

func (u *User) getRoom(name RoomToken) (IRoom, bool) {
	r, ok := u.Rooms[name]
	return r, ok
}

func (u *User) sendWorker() {
	defer func() {
		if err := recover(); err != nil {
			u.StopListening()
		}
	}()

	for {
		select {
		case m := <-u.output:
			if err := u.Conn.WriteJSON(m); err != nil {
				u.StopListening()
				u.stop()
			}
		case <-u.cancel:
			return
		}
	}
}

func (u *User) receiveWorker() {
	for {
		m := &InMessage{}

		// read a message
		if err := u.Conn.ReadJSON(m); err != nil {
			if websocket.IsCloseError(err, wsCloseErrors...) {
				u.onDisconnected()
				u.stop()
				return
			}
			u.LG.Sugar.Infof("wrong_message: %v", u.GetID())
			u.Consume(&OutMessage{
				RoomToken: m.GetRoomID(),
				Status:    StatusError,
				Payload:   "Wrong message",
			})
			continue
		}

		u.input <- m
	}
}

func (u *User) removeFromAllRooms() error {
	for _, r := range u.Rooms {
		err := u.RemoveFromRoom(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) onDisconnected() {
	for _, r := range u.Rooms {
		r.OnDisconnected(u)
	}
}

func (u *User) stop() {
	u.LG.Sugar.Infof("ws closed, uid: %s", u.GetID())
	u.bClosed = true
	u.Conn.Close()
	u.cancel <- ""
	u.cancel <- ""
}