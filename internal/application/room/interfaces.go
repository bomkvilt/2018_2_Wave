package room

import "time"

/** How does it works?
 *	. ## begining of a WS handler
 *  . In WS handler we crate an instance if IUser interface
 *	. Than we call a Listen method in the instance in goroutine
 *	. The Listen() goes into a cycle and listen to the WS connection
 *	. Then in a main thread we call a AddToRoom() and send a main room into
 *	. ## and of the handler
 */

//go:generate easyjson .

type RoomToken string
type UserID string
type RoomType string
type RoomFactory func(id RoomToken, step time.Duration, _ IRoomManager, db interface{}) IRoom

// IInMessage - message from a user
type IInMessage interface {
	GetRoomID() RoomToken       // target room id
	GetSignal() string          // message method
	GetPayload() interface{}    // message payload
	ToStruct(interface{}) error // unmurshal data to struct
}

// IOutMessage - message to a client
type IOutMessage interface {
	GetRoomID() RoomToken         // message room id
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// IRouteResponse - response from route
type IRouteResponse interface {
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// IUser - client websocket wrapper
type IUser interface {
	GetName() string
	GetID() UserID              // User id
	AddToRoom(IRoom) error      // order to add self into the room
	RemoveFromRoom(IRoom) error // order to romve self from the room
	Listen() error              // Listen to messages
	StopListening() error       // Stop listening
	Consume(IOutMessage) error  // Send message to user

	Task(func())
}

// IRoom - abstruct room inteface
type IRoom interface {
	GetID() RoomToken                     // room id
	GetType() RoomType                    // room type
	Run() error                           // run thr room
	Stop() error                          // stop the room
	AddUser(IUser) error                  // add the user to the room
	RemoveUser(IUser) error               // remove  the user from the room
	OnDisconnected(IUser)                 // inform the room the user was disconnected
	ApplyMessage(IUser, IInMessage) error // send message to the room

	Task(func())
	IsAbleToRemove(IUser) bool
}

// IRoomManager -
type IRoomManager interface {
	CreateLobby(RoomType, RoomToken) (IRoom, error)
	RemoveLobby(RoomToken, IUser) error
}
