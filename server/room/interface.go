package room

/** How does it works?
 *	. ## begining of a WS handler
 *  . In WS handler we crate an instance if IUser interface
 *	. Than we call a Listen method in the instance in goroutine
 *	. The Listen() goes into a cycle and listen to the WS connection
 *	. Then in a main thread we call a AddToRoom() and send a main room into
 *	. ## and of the handler
 */

type RoomID string
type UserID string

// IInMessage - message from a user
type IInMessage interface {
	GetRoomID() RoomID          // target room id
	GetSignal() string          // message method
	ToStruct(interface{}) error // unmurshal data to struct
}

// IOutMessage - message to a client
type IOutMessage interface {
	GetRoomID() RoomID            // message room id
	GetPayload() []byte           // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// IUser - client websocket wrapper
type IUser interface {
	GetID() UserID              // User id
	AddToRoom(IRoom) error      // order to add self into the room
	RemoveFromRoom(IRoom) error // order to romve self from the room
	Listen() error              // Listen to messages
	StopListening() error       // Stop listening
	Send(IOutMessage) error     // Send message to user
}

// IRoom - abstruct room inteface
type IRoom interface {
	GetID() RoomID                // room id
	Run() error                   // run thr room
	Stop() error                  // stop the room
	AddUser(IUser) error          // add the user to the room
	RemoveUser(IUser) error       // remove  the user from the room
	SendMessage(IInMessage) error // send message to the room
}
