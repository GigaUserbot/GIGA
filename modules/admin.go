package modules

import (
	"fmt"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/stylisehelper"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

func (m *module) LoadAdmin(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("ADMIN")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("admin", "help of admin")
	dispatcher.AddHandler(handlers.NewCommand("ban", authorised(ban)))
	dispatcher.AddHandler(handlers.NewCommand("unban", authorised(unban)))
}

func ban(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	if chat.IsAUser() {
		return dispatcher.EndGroups
	}
	target, err := utils.ExtractUser(ctx, u.EffectiveMessage, chat)
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: fmt.Sprintf("failed to ban user: %s", err.Error()),
		})
		return dispatcher.EndGroups
	}
	_, err = ctx.BanChatMember(chat.GetID(), target, 0)
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: fmt.Sprintf("failed to ban user: %s", err.Error()),
		})
		return dispatcher.EndGroups
	} else {
		text := stylisehelper.Start(styling.Plain("Successfully banned "))
		text.Mention("this user", target).Plain(".")
		builder := gotgproto.Sender.Self().Edit(u.EffectiveMessage.ID)
		builder.StyledText(ctx, text.StoArray...)
	}
	return dispatcher.EndGroups
}

func unban(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	if chat.IsAUser() {
		return dispatcher.EndGroups
	}
	target, err := utils.ExtractUser(ctx, u.EffectiveMessage, chat)
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: fmt.Sprintf("failed to unban user: %s", err.Error()),
		})
		return dispatcher.EndGroups
	}
	_, err = ctx.UnbanChatMember(chat.GetID(), target, 0)
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: fmt.Sprintf("failed to unban user: %s", err.Error()),
		})
		return dispatcher.EndGroups
	} else {
		text := stylisehelper.Start(styling.Plain("Successfully unbanned "))
		text.Mention("this user", target).Plain(".")
		builder := gotgproto.Sender.Self().Edit(u.EffectiveMessage.ID)
		builder.StyledText(ctx, text.StoArray...)
	}
	return dispatcher.EndGroups
}
