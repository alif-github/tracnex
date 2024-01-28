package grochat_ws_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrr/fastws"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model"
	"strings"
	"time"
)

type GroChatWSHandler struct {
	conn           *fastws.Conn
	url 		   string
	token          string
	clientId       string
	signId		   string
	reconnectSignal chan bool
}

type OnMessage func(message []byte)

func NewGroChatWSHandler() *GroChatWSHandler {
	groChatWS := config.ApplicationConfiguration.GetGroChatWS()

	return &GroChatWSHandler{
		url:             groChatWS.Host + groChatWS.PathRedirect.WS,
		reconnectSignal: make(chan bool),
	}
}

func (g *GroChatWSHandler) Dial() error {
	//groChatWS := config.ApplicationConfiguration.GetGroChatWS()

	conn, err := fastws.Dial(g.url)
	if err != nil {
		return err
	}

	//conn, err := fastws.DialTLS(g.url, &tls.Config{ServerName: groChatWS.Host + groChatWS.PathRedirect.WS, InsecureSkipVerify: true})
	//if err != nil {
	//	return err
	//}

	conn.WriteTimeout = 0
	conn.ReadTimeout = 0

	g.conn = conn
	return nil
}

func (g *GroChatWSHandler) Authenticate(username, password string) error {
	token, clientId := g.requestLogin(username, password)

	signId, err := g.fetchSignID(token, clientId)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to fetch sign id : %s", err.Error()))
	}

	g.token = token
	g.clientId = clientId
	g.signId = signId

	return nil
}

func (g *GroChatWSHandler) Start() {
	go func() {
		if err := g.readMessage(g.handleOnMessage); err != nil {
			return
		}
	}()

	g.handshake()
}

func (g *GroChatWSHandler) ReconnectSignal() chan bool {
	return g.reconnectSignal
}

func (g *GroChatWSHandler) Write(message []byte) error {
	if _, err := g.conn.Write(message); err != nil {
		if err == fastws.EOF {
			g.reconnect()

			return err
		}

		return err
	}

	return nil
}

func (g *GroChatWSHandler) SendNotification(destinationId string, messages []string) error {
	roomId := fmt.Sprintf("%s_room_%s", g.clientId, destinationId)

	for _, message := range messages {
		messageId := uuid.New().String()

		payload := &ChatMessage{
			Type:    constanta.MessageTypeChat,
			Message: Message{
				SourceId:      g.clientId,
				DestinationId: destinationId,
				Identity: Identity{
					ClientId: g.clientId,
					Sign:     g.signId,
				},
				MessageDetail: MessageDetail{
					Guarantee:    true,
					Type:         constanta.ChatTypeSend,
					MessageId:    messageId,
					MessageModel: MessageModel{
						Content:  ContentText{
							Text: message,
						},
						CreatedAt:     g.getCurrentMilliseconds(),
						IsEncrypted:   "N",
						RoomId:        roomId,
						TypeMessageId: 1,
						VectorX:       1,
						VectorY:       0,
					},
				},
			},
		}

		if _, err := g.conn.Write(payload.Bytes()); err != nil {
			g.logError(fmt.Sprintf("Failed to send notification : %s", err.Error()), 500)
			return err
		}

		g.logInfo(fmt.Sprintf("Notification has been sent to %s", roomId), 200)
	}

	return nil
}

func (g *GroChatWSHandler) getCurrentMilliseconds() int64 {
	now := time.Now()
	dur := now.Sub(time.Unix(0, 0))

	return dur.Milliseconds()
}

func (g *GroChatWSHandler) handleOnMessage(message []byte) {
	messageType := &MessageType{}

	if err := json.Unmarshal(message, messageType); err != nil {
		g.logError(fmt.Sprintf("(GroChat WS) Unmarshal MessageType error : %s", err.Error()), 500)
		return
	}

	if messageType.Type != constanta.MessageTypeChat {
		return
	}

	g.logInfo(fmt.Sprintf("(GroChat WS) Received message : %s", string(message)), 200)
	g.sendAckMessage(message)
}

func (g *GroChatWSHandler) sendAckMessage(message []byte) {
	chat := &ChatMessage{}

	if err := json.Unmarshal(message, chat); err != nil {
		g.logError(fmt.Sprintf("(GroChat WS) Unmarshal ChatMessage error : %s", err.Error()), 500)
		return
	}

	if chat.Message.MessageDetail.Type != constanta.ChatTypeSend {
		return
	}

	if !g.isMessageFromSystem(chat) {
		return
	}

	g.updateMessageStatus(chat, chat.Message.DestinationId, constanta.ChatTypeAck)
}

func (g *GroChatWSHandler) isMessageFromSystem(chat *ChatMessage) bool {
	return chat.Message.SourceId == g.clientId
}

func (g *GroChatWSHandler) updateMessageStatus(chat *ChatMessage, destinationId string, status string) {
	payload := &ChatMessage{
		Type:    constanta.MessageTypeChat,
		Message: Message{
			SourceId:      g.clientId,
			DestinationId: destinationId,
			Identity: Identity{
				ClientId: g.clientId,
				Sign:     g.signId,
			},
			MessageDetail: MessageDetail{
				Guarantee:    true,
				Type:         status,
				MessageId:    chat.Message.MessageDetail.MessageId,
				MessageModel: MessageModel{
					RoomId:   chat.Message.MessageDetail.MessageModel.RoomId,
				},
			},
		},
	}

	if err := g.Write(payload.Bytes()); err != nil {
		g.logError(fmt.Sprintf("(GroChat WS) Update message status error : MessageId = %s, SourceId = %s, DestinationId = %s, Type = %s -> %s, Error = %s",
			payload.Message.MessageDetail.MessageId,
			payload.Message.SourceId,
			payload.Message.DestinationId,
			chat.Message.MessageDetail.Type,
			status,
			err.Error()), 500)
		return
	}

	g.logInfo(fmt.Sprintf("(GroChat WS) Update message has been sent  : MessageId = %s, SourceId = %s, DestinationId = %s, Type = %s -> %s",
		payload.Message.MessageDetail.MessageId,
		payload.Message.SourceId,
		payload.Message.DestinationId,
		chat.Message.MessageDetail.Type,
		status), 200)
	return
}

func (g *GroChatWSHandler) handshake() {
	handshakeMessage := GroChatWSHandshake{
		Token:    g.token,
		ClientId: g.clientId,
		Sign:     g.signId,
	}

	message, _ := json.Marshal(handshakeMessage)

	for {
		if err := g.Write(message); err != nil {
			continue
		}

		return
	}
}

func (g *GroChatWSHandler) KeepAlive() {
	for {
		time.Sleep(2 * time.Minute)

		body := GroChatWSKeepAlive{Type: "keepalive"}
		bytes, _ := json.Marshal(body)

		if err := g.Write(bytes); err != nil {
			g.logError("Failed to send keep alive message", 500)
			continue
		}

		g.logInfo("Keep alive message has been sent", 200)
	}
}

func (g *GroChatWSHandler) reconnect() {
	g.reconnectSignal <- true
}

func (g *GroChatWSHandler) readMessage(onMessage OnMessage) error {
	for {
		_, message, err := g.conn.ReadMessage(nil)
		if err != nil {
			/*
				Disconnected
			*/
			if err == fastws.EOF {
				g.reconnect()
				return err
			}

			continue
		}

		onMessage(message)
	}
}

func (g *GroChatWSHandler) requestLogin(username, password string) (token string, clientId string) {
	groChatServer := config.ApplicationConfiguration.GetGroChat()

	path := fmt.Sprintf("%s%s", groChatServer.Host, groChatServer.PathRedirect.Login)

	headerRequest := make(map[string][]string)
	headerRequest["X-DEVICE"] = []string{"backend"}
	headerRequest["Content-Type"] = []string{"application/json"}

	req := util2.StructToJSON(LoginWS{
		Username:  username,
		Password:  password,
	})

	var (
		response LoginResponse
	)

	isLoginSuccess := false

	for !isLoginSuccess {
		g.logInfo(fmt.Sprintf("logging in as %s", username), 200)

		statusCode, _, bodyResult, err := g.request(path, headerRequest, req, "POST")
		if err != nil {
			g.logError(fmt.Sprintf("Failed to login as %s : %s", username, err.Error()), 500)
			continue
		}

		if err = json.Unmarshal([]byte(bodyResult), &response); err != nil {
			g.logError(fmt.Sprintf("Failed to login as %s : %s", username, err.Error()), 500)
			continue
		}

		if statusCode == 200 && response.Status == 1 {
			isLoginSuccess = true
		}
	}

	g.logInfo(fmt.Sprintf("logged in as %s successfully", username), 200)

	return response.Data.UserToken, response.Data.UserModel.Auth.ClientId
}

func (g *GroChatWSHandler) fetchSignID(token, clientId string) (string, error) {
	groChatServer := config.ApplicationConfiguration.GetGroChat()

	path := fmt.Sprintf("%s%s/%s", groChatServer.Host, groChatServer.PathRedirect.SignId, clientId)
	headerRequest := make(map[string][]string)
	headerRequest["Authorization"] = []string{token}

	var (
		response GroWSResponseSuccess
	)

	statusCode, _, bodyResult, err := g.request(path, headerRequest, "", "GET")
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(bodyResult), &response); err != nil {
		return "", err
	}

	if statusCode == 200 && response.Status == 1 {
		return response.Data.PayloadGro.SignUuid, nil
	}

	return "", nil
}

func (*GroChatWSHandler) request(urlAddress string, header map[string][]string, body string, method string) (statusCode int, headerResult map[string][]string, bodyResult string, err error) {
	var (
		reqURL   *url.URL
		request  *http.Request
		response *http.Response
	)

	reqURL, err = url.Parse(urlAddress)
	if err != nil {
		return
	}

	if header == nil {
		header = make(map[string][]string)
	}

	request = &http.Request{
		Method: strings.ToUpper(method),
		URL:    reqURL,
		Header: header,
		Body:   ioutil.NopCloser(strings.NewReader(body)),
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer func() {
		err = response.Body.Close()
	}()

	bodyResultByte, _ := ioutil.ReadAll(response.Body)
	bodyResult = string(bodyResultByte)
	statusCode = response.StatusCode
	headerResult = response.Header

	return
}

func (*GroChatWSHandler) logInfo(message string, status int) {
	logModel := model.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion())
	logModel.Status = status
	logModel.Message = message

	util2.LogInfo(logModel.ToLoggerObject())
}

func (*GroChatWSHandler) logError(message string, status int) {
	logModel := model.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion())
	logModel.Status = status
	logModel.Message = message

	util2.LogError(logModel.ToLoggerObject())
}