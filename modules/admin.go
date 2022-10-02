package modules

import (
	"fmt"
	"html"

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
	helpmaker.SetModuleHelp("admin", `
This module provides help for the basic admin moderation rights like banning, unbanning etc.

<b>Commands</b>:
 • <code>.ban `+html.EscapeString("<username/reply_to_message>")+`</code>: Use this command to ban a user.   
 • <code>.unban `+html.EscapeString("<username/reply_to_message>")+`</code>: Use this command to unban a user.   
 • <code>.del `+html.EscapeString("<reply_to_message>")+`</code>: Use this command to delete the replied to message.   
 • <code>.purge `+html.EscapeString("<reply_to_message>")+`</code>: Use this command to delete messages from the current one to the replied to message.   
	`)
	dispatcher.AddHandler(handlers.NewCommand("ban", authorised(ban)))
	dispatcher.AddHandler(handlers.NewCommand("unban", authorised(unban)))
	dispatcher.AddHandler(handlers.NewCommand("del", authorised(del)))
	dispatcher.AddHandler(handlers.NewCommand("purge", authorised(purge)))
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
	_, err = ctx.UnbanChatMember(chat.GetID(), target)
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

func del(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	reply := u.EffectiveMessage.ReplyTo.ReplyToMsgID
	if reply == 0 {
		ctx.DeleteMessages(chat.GetID(), []int{u.EffectiveMessage.ID})
		return dispatcher.EndGroups
	}
	ctx.DeleteMessages(chat.GetID(), []int{u.EffectiveMessage.ID, reply})
	return dispatcher.EndGroups
}

func purge(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	reply := u.EffectiveMessage.ReplyTo.ReplyToMsgID
	if reply == 0 {
		ctx.DeleteMessages(chat.GetID(), []int{u.EffectiveMessage.ID})
		return dispatcher.EndGroups
	}
	toDel := []int{u.EffectiveMessage.ID, reply}
	for i := reply; i < u.EffectiveMessage.ID; i++ {
		toDel = append(toDel, i)
	}
	ctx.DeleteMessages(chat.GetID(), toDel)
	return dispatcher.EndGroups
}
