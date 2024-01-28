package grochat_request

import (
	"reflect"
)

type SendMessageGroChatRequest struct {
	Status  string     `json:"status" default:"success"`
	Message string     `json:"message" default:"Message sent"`
	Data    DetailData `json:"data"`
}

type DetailData struct {
	ClientID       string         `json:"clientID"`
	Username       string         `json:"userName"`
	PhoneNumber    string         `json:"phoneNumber"`
	MessageType    string         `json:"messageType" default:"1"`
	MessageContent MessageContent `json:"messageContent"`
}

type MessageContent struct {
	Message          string `json:"message"`
	MessageEncrypted string `json:"messageEncrypted" default:"N"`
	Status           string `json:"status" default:"1"`
	StatusInfo       string `json:"statusInfo" default:"sent"`
	AdditionalKey    string `json:"additionalKey"`
	LocalID          string `json:"localID"`
	TypeRoom         string `json:"typeRoom" default:"U"`
}

func (input SendMessageGroChatRequest) GetDefault(p *SendMessageGroChatRequest) {
	typ := reflect.TypeOf(SendMessageGroChatRequest{})

	f, _ := typ.FieldByName("Status")
	p.Status = f.Tag.Get("default")

	f, _ = typ.FieldByName("Message")
	p.Message = f.Tag.Get("default")

	typ = reflect.TypeOf(DetailData{})
	f, _ = typ.FieldByName("MessageType")
	p.Data.MessageType = f.Tag.Get("default")

	typ = reflect.TypeOf(MessageContent{})
	f, _ = typ.FieldByName("MessageEncrypted")
	p.Data.MessageContent.MessageEncrypted = f.Tag.Get("default")

	f, _ = typ.FieldByName("Status")
	p.Data.MessageContent.Status = f.Tag.Get("default")

	f, _ = typ.FieldByName("StatusInfo")
	p.Data.MessageContent.StatusInfo = f.Tag.Get("default")

	f, _ = typ.FieldByName("TypeRoom")
	p.Data.MessageContent.TypeRoom = f.Tag.Get("default")

	return
}
