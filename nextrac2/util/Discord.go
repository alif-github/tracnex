package util

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"os"
	"strings"
	"time"
)

func DiscordSendThread(contextModel applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName = "Discord.go"
		funcName = "DiscordSendThread"
		channel  = config.ApplicationConfiguration.GetDiscordLogChannelId()
		sess     = serverconfig.ServerAttribute.DiscordConn
		errS     error
	)

	defer func() {
		if errS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errS)
			return
		}
	}()

	c := contextModel.LoggerModel
	timeErr := time.Now().Format(constanta.DefaultDBSQLTimeFormat)
	location := strings.Split(c.Class, ",")
	str := fmt.Sprintf("------------------------------------\nRequest ID : %s\nError Status : %d\nFile Name : %s\nFunction Name : %s\nTime Error : %s\nError Code : %s\nMessage : %s\n------------------------------------",
		c.RequestID,
		c.Status,
		strings.ReplaceAll(location[0], "[", ""),
		strings.ReplaceAll(location[1], "]", ""),
		timeErr,
		c.Code,
		c.Message)

	fmt.Println("Channel : ", channel)
	_, errS = sess.ChannelMessageSend(channel, str)
	if errS != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func DiscordInfoServerRunning(args string) {
	var (
		channel = config.ApplicationConfiguration.GetDiscordLogChannelId()
		sess    = serverconfig.ServerAttribute.DiscordConn
		owner   = "-"
	)

	if args == "local" {
		owner, _ = os.Hostname()
	}

	ip, _ := util.GenerateIPAddress(config.ApplicationConfiguration.GetServerEthernet())
	_, _ = sess.ChannelMessageSend(channel, fmt.Sprintf("------------------------------------\nSERVER NEXTRAC RE-RUNNING BLOKKKK !!! \U0001F590\nTime Started : %s\nIP : %s\nOwner : %s\n------------------------------------\n",
		time.Now().Format(constanta.DefaultDBSQLTimeFormat), ip, owner))
}
