package room

import (
	"encoding/json"
)

// ----------------| InMessage

// InMessage - default IInMessage
type InMessage struct {
	RoomID  RoomID
	Signal  string
	Payload []byte
}

func (im *InMessage) GetRoomID() RoomID { return im.RoomID }
func (im *InMessage) GetSignal() string { return im.Signal }
func (im *InMessage) ToStruct(s interface{}) error {
	return json.Unmarshal(im.Payload, s)
}

// ----------------| OutMessage

// OutMessage - default IOutMessage
type OutMessage struct {
	RoomID  RoomID
	Status  string
	Payload []byte
}

func (om *OutMessage) GetRoomID() RoomID  { return om.RoomID }
func (om *OutMessage) GetStatus() string  { return om.Status }
func (om *OutMessage) GetPayload() []byte { return om.Payload }
func (om *OutMessage) FromStruct(s interface{}) (err error) {
	om.Payload, err = json.Marshal(s)
	return err
}
