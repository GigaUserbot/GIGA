package modules

import (
	"fmt"
	"strings"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/stylisehelper"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/db"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

func (m *module) LoadAfk(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("AFK")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	dispatcher.AddHandler(handlers.NewCommand("afk", authorised(afk)))
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.All, checkAfk), 1)
}

func afk(ctx *ext.Context, u *ext.Update) error {
	args := strings.Fields(u.EffectiveMessage.Message)
	chat := u.EffectiveChat()
	if len(args) > 1 {
		switch args[1] {
		case "on", "true":
			reason := ""
			if len(args) > 2 {
				reason = strings.Join(args[2:], " ")
			}
			go db.UpdateAFK(true, reason)
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID: u.EffectiveMessage.ID,
				Message: fmt.Sprintf("Turned on AFK mode.%s", func() string {
					if reason != "" {
						return fmt.Sprintf("\nReason: %s", reason)
					}
					return reason
				}()),
			})
		case "off", "false":
			go db.UpdateAFK(false, "")
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      u.EffectiveMessage.ID,
				Message: "Turned off AFK mode.",
			})
		default:
			ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
				ID:      u.EffectiveMessage.ID,
				Message: "AFK: Invalid Arguments",
			})
		}
	} else {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: "AFK: No arguments were provided.",
		})
	}
	return dispatcher.EndGroups
}

func checkAfk(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	user := u.EffectiveUser()
	if user != nil && user.Bot {
		// Don't reply to bots ffs
		return nil
	}
	if !(u.EffectiveMessage.Mentioned || (chat.IsAUser() && chat.GetID() != gotgproto.Self.ID)) {
		return nil
	}
	afk := db.GetAFK()
	if !afk.Toggle {
		return nil
	}
	text := stylisehelper.Start(styling.Plain("I'm currently AFK"))
	if afk.Reason != "" {
		text.Plain("\nReason: ").Code(afk.Reason)
	}
	ctx.Reply(u, text.StoArray, nil)
	return nil
}
