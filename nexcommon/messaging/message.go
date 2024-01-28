package messaging

import (
	"errors"
	"nexsoft.co.id/nexcommon/util"
)

type Message struct {
	NexsoftMessage struct {
		Header struct {
			MessageID   string `json:"message_id"`
			UserID      string `json:"user_id"`
			Password    string `json:"password"`
			Version     string `json:"version"`
			PrincipalID string `json:"principal_id"`
			Timestamp   string `json:"timestamp"`
			Action      struct {
				Class  string   `json:"class_name"`
				Type   string   `json:"type_name"`
				Custom []Custom `json:"custom"`
			} `json:"action"`
			Custom []Custom `json:"custom"`
		} `json:"header"`
		Payload interface{} `json:"payload"`
		Custom  []Custom    `json:"custom"`
	} `json:"nexsoft_message"`
}

type Custom struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (input Message) CheckMandatoryOrigin() error {
	if input.NexsoftMessage.Header.MessageID == "" {
		return errors.New("MESSAGE_ID_MISSING")
	}
	if input.NexsoftMessage.Header.Version == "" {
		return errors.New("VERSION_MISSING")
	}
	if input.NexsoftMessage.Header.UserID == "" {
		return errors.New("USER_ID_MISSING")
	}
	if input.NexsoftMessage.Header.Password == "" {
		return errors.New("PASSWORD_MISSING")
	}
	if input.NexsoftMessage.Header.PrincipalID == "" {
		return errors.New("PRINCIPAL_ID_MISSING")
	}
	if input.NexsoftMessage.Header.Action.Class == "" {
		return errors.New("ACTION_CLASS_MISSING")
	}
	if input.NexsoftMessage.Header.Action.Type == "" {
		return errors.New("ACTION_TYPE_MISSING")
	}
	timestampValid, timestamp := util.IsTimestampValid(input.NexsoftMessage.Header.Timestamp)

	if timestampValid {
		input.NexsoftMessage.Header.Timestamp = timestamp
		return nil
	} else {
		return errors.New("INVALID_TIMESTAMP")
	}
}

func (input Message) String() (output string) {
	return util.StructToJSON(input)
}
