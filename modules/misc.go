package modules

import (
	"encoding/json"
	"github.com/anonyindian/gotgproto"
	"strconv"
	"time"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/entityhelper"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/tg"
)

func (m *module) LoadMisc(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("MISC")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("misc", `
	This module provides help for the miscellaneous features like ping, json, etc.
	
	<b>Commands</b>:
	 • <code>.ping</code>: Use this command to check ping between telegram and userbot client.   
	 • <code>.json</code>: Get JSON output of a message.   
	 • <code>.alive</code>: Use this command to check whether the userbot is alive or not.   
`)
	dispatcher.AddHandler(handlers.NewCommand("ping", authorised(ping)))
	dispatcher.AddHandler(handlers.NewCommand("save", authorised(dotsave)))
	dispatcher.AddHandler(handlers.NewCommand("alive", authorised(alive)))
	dispatcher.AddHandler(handlers.NewCommand("json", authorised(jsonify)))
}

func jsonify(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	msg := u.EffectiveMessage
	if id := msg.ReplyTo.ReplyToMsgID; id != 0 {
		m, err := ctx.GetMessages(chat.GetID(), []tg.InputMessageClass{&tg.InputMessageID{
			ID: id,
		}})
		if err != nil {
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      msg.ID,
				Message: "failed to jsonify: " + err.Error(),
			})
			return dispatcher.EndGroups
		}
		md, ok := m[0].(*tg.Message)
		if ok {
			msg = md
		}
	}
	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: "failed to jsonify: " + err.Error(),
		})
		return dispatcher.EndGroups
	}
	text := entityhelper.Code(string(b))
	ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		ID:       u.EffectiveMessage.ID,
		Message:  text.GetString(),
		Entities: text.GetEntities(),
	})
	return dispatcher.EndGroups
}

func alive(ctx *ext.Context, u *ext.Update) error {
	text := entityhelper.Bold(`
The GIGA Userbot is currently up and working fine,
Written using @gotgproto by @GIGADevs.
`)
	ctx.EditMessage(u.EffectiveChat().GetID(), &tg.MessagesEditMessageRequest{
		ID:       u.EffectiveMessage.ID,
		Message:  text.GetString(),
		Entities: text.GetEntities(),
	})
	return dispatcher.EndGroups
}

func dotsave(ctx *ext.Context, u *ext.Update) error {
	msg := u.EffectiveMessage
	if msg.ReplyTo.ReplyToMsgID == 0 {
		ctx.EditMessage(u.EffectiveChat().GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: "Reply to a message!",
		})
		return dispatcher.EndGroups
	}
	ctx.DeleteMessages(u.EffectiveChat().GetID(), []int{u.EffectiveMessage.ID})
	ctx.ForwardMessage(u.EffectiveChat().GetID(), gotgproto.Self.ID,
		&tg.MessagesForwardMessagesRequest{ID: []int{u.EffectiveMessage.ReplyTo.ReplyToMsgID}},
	)
	return dispatcher.EndGroups
}

func ping(ctx *ext.Context, u *ext.Update) error {
	timeThen := time.Now()
	utils.TelegramClient.Ping(ctx)
	timeNow := time.Since(timeThen)
	text := entityhelper.Plain("PONG\n").Code(strconv.FormatInt(timeNow.Milliseconds(), 10) + "ms")
	ctx.EditMessage(u.EffectiveChat().GetID(), &tg.MessagesEditMessageRequest{
		ID:       u.EffectiveMessage.ID,
		Message:  text.String,
		Entities: text.Entities,
	})
	return dispatcher.EndGroups
}
