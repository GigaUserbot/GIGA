package modules

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/entityhelper"
	"github.com/anonyindian/gotgproto/types"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/db"
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
	dispatcher.AddHandler(handlers.NewCommand("alive", authorised(alive)))
	dispatcher.AddHandler(handlers.NewCommand("json", authorised(jsonify)))
	dispatcher.AddHandler(handlers.NewCommand("taglogger", authorised(tagLogger)))
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.All, checkTags), -1)
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

func ping(ctx *ext.Context, u *ext.Update) error {
	timeThen := time.Now()
	utils.TelegramClient.Ping(ctx)
	timeNow := time.Since(timeThen)
	text := entityhelper.Plain("🧭PONG🧭\n").Code(strconv.FormatInt(timeNow.Milliseconds(), 10) + "ms")
	ctx.EditMessage(u.EffectiveChat().GetID(), &tg.MessagesEditMessageRequest{
		ID:       u.EffectiveMessage.ID,
		Message:  text.String,
		Entities: text.Entities,
	})
	return dispatcher.EndGroups
}

func tagLogger(ctx *ext.Context, u *ext.Update) error {
	args := strings.Fields(u.EffectiveMessage.Message)
	chat := u.EffectiveChat()
	if len(args) > 1 {
		switch args[1] {
		case "on", "true":
			go db.TagLogger(true)
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      u.EffectiveMessage.ID,
				Message: "All mentions will be logged now.",
			})
		case "off", "false":
			go db.TagLogger(false)
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      u.EffectiveMessage.ID,
				Message: "Mentions will not be logged now.",
			})
		default:
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      u.EffectiveMessage.ID,
				Message: "TagLogger: Invalid Arguments",
			})
		}
	} else {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: "TagLogger: No arguments were provided.",
		})
	}
	return dispatcher.EndGroups
}

func checkTags(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	user := u.EffectiveUser()
	if u.EffectiveMessage.Out {
		return nil
	}
	if user != nil && user.Bot {
		return nil
	}
	if !(u.EffectiveMessage.Mentioned || (chat.IsAUser() && chat.GetID() != gotgproto.Self.ID)) {
		return nil
	}
	if !db.GetTagLogger() {
		return nil
	}
	logsGroup := db.GetSettings().LogsGroup
	if chat.IsAUser() {
		chatUser := chat.(*types.User)
		text := entityhelper.Bold("New Private Message")
		text.Bold("\nBy: ")
		if chatUser.Username != "" {
			text.Plain("@").Plain(chatUser.Username)
		} else {
			text.Code(chatUser.FirstName + " " + chatUser.LastName)
		}
		ctx.SendMessage(logsGroup, &tg.MessagesSendMessageRequest{
			Message:  text.String,
			Entities: text.Entities,
		})
		return nil
	}
	text := entityhelper.Bold("New Chat Mention")
	text.Bold("\nBy: ")
	if user.Username != "" {
		text.Plain("@").Plain(user.Username)
	} else {
		text.Code(user.FirstName + " " + user.LastName)
	}
	if link := getChatMessageLink(chat, u.EffectiveMessage.ID); link != "" {
		text.Bold("\nLink: ").Plain(link)
	}
	ctx.SendMessage(logsGroup, &tg.MessagesSendMessageRequest{
		Message:   text.String,
		Entities:  text.Entities,
		NoWebpage: true,
	})
	return nil
}

func getChatMessageLink(c types.EffectiveChat, msgId int) string {
	if c.IsAChat() {
		return ""
	}
	chat := c.(*types.Channel)
	if chat.Username != "" {
		return fmt.Sprintf("t.me/%s/%d", chat.Username, msgId)
	}
	return fmt.Sprintf("t.me/c/%d/%d", chat.ID, msgId)
}
