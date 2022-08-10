package modules

import (
	"os"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/tg"
)

func (m *module) LoadProg(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("PROG")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("prog", `
	This module provides help for the commads like stopping userbot, etc.
	
	<b>Commands</b>:
	 • <code>.killub</code>: Use this command to turn off the userbot.   
`)
	dispatcher.AddHandler(handlers.NewCommand("killub", authorised(killub)))
	dispatcher.AddHandler(handlers.NewCommand("restart", authorised(restart)))
}

func killub(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		ID:      u.EffectiveMessage.ID,
		Message: "Exiting...",
	})
	os.Exit(1)
	return dispatcher.EndGroups
}

func restart(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		ID:      u.EffectiveMessage.ID,
		Message: "Restarting",
	})
	utils.Restart(5, chat.GetID(), u.EffectiveMessage.ID, "Restarted Successfully!")
	return dispatcher.EndGroups
}
