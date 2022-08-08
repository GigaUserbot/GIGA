package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
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
	command := fmt.Sprintf("run main.go -delay=5 -chat=%d -msg=%d", chat.GetID(), u.EffectiveMessage.ID)
	cmd := exec.Command("go", strings.Fields(command)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Start()
	cmd.Process.Release()
	os.Exit(1)
	return dispatcher.EndGroups
}
