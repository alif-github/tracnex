package applicationModel

import (
	"go.uber.org/zap"
	"nexsoft.co.id/nexcommon/util"
	"os"
)

type LoggerModel struct {
	IP          string `json:"ip"`
	Category    string `json:"category"`
	PID         int    `json:"pid"`
	Thread      string `json:"thread"`
	RequestID   string `json:"request_id"`
	Source      string `json:"source"`
	AccessToken string `json:"access_token"`
	ClientID    string `json:"client_id"`
	UserID      string `json:"user_id" `
	Resource    string `json:"resource"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Class       string `json:"class"`
	Time        int64  `json:"time"`
	ByteIn      int    `json:"byte_in"`
	ByteOut     int    `json:"byte_out"`
	Status      int    `json:"status" `
	Code        string `json:"code"`
	Message     string `json:"message"`
}

func GenerateLogModel(version string, application string) (output LoggerModel) {
	output.IP = "-"
	output.Category = "-"
	output.PID = os.Getpid()
	output.Thread = "-"
	output.RequestID = "-"
	output.Source = "-"
	output.AccessToken = "-"
	output.Resource = "-"
	output.Application = application
	output.Version = version
	output.Class = "-"
	output.Code = "-"
	output.Message = "-"
	return output
}

func (object LoggerModel) String() string {
	return util.StructToJSON(object)
}

func (object LoggerModel) ToLoggerObject() (output []zap.Field) {
	output = append(output, zap.String("ip", object.IP))
	output = append(output, zap.String("category", object.Category))
	output = append(output, zap.Int("pid", object.PID))
	output = append(output, zap.String("class", object.Class))
	output = append(output, zap.String("thread", object.Thread))
	output = append(output, zap.String("request_id", object.RequestID))
	output = append(output, zap.String("source", object.Source))
	output = append(output, zap.String("access_token", object.AccessToken))
	output = append(output, zap.String("client_id", object.ClientID))
	output = append(output, zap.String("user_id", object.UserID))
	output = append(output, zap.String("resource", object.Resource))
	output = append(output, zap.String("application", object.Application))
	output = append(output, zap.String("version", object.Version))
	output = append(output, zap.Int64("time", object.Time))
	output = append(output, zap.Int("byte_in", object.ByteIn))
	output = append(output, zap.Int("byte_out", object.ByteOut))
	output = append(output, zap.Int("status", object.Status))
	output = append(output, zap.String("code", object.Code))
	output = append(output, zap.String("message", object.Message))

	return output
}

type PanicLogger struct {
	FileName     string
	FunctionName string
	Input        interface{}
	ErrorMessage string
}
