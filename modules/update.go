package modules

import (
	"fmt"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/entityhelper"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/tg"
)

func (m *module) LoadUpdate(dp *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("UPDATER")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("updater", `
	This module provides help for the updater module.
	
	<b>Commands</b>:
	 • <code>.changelog</code>: Get info about new updates.   
	 • <code>.update</code>: Update the userbot.   
`)
	dp.AddHandler(handlers.NewCommand("changelog", changeLog))
	dp.AddHandler(handlers.NewCommand("update", update))
}

func update(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	msg := u.EffectiveMessage
	update, changed := utils.CheckChanges()
	if !changed {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			Message: "You're currently running the latest version.",
			ID:      msg.ID,
		})
		return dispatcher.EndGroups
	}
	text := entityhelper.Plain("Updating to ").Bold("GIGA ")
	text.Code(fmt.Sprintf("v%s", update.Version))
	ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		Message:  text.String,
		Entities: text.Entities,
		ID:       msg.ID,
	})
	if err := utils.DoUpdate(update.Version, chat.GetID(), msg.ID); err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			Message: fmt.Sprintf("failed to update: %s", err.Error()),
			ID:      msg.ID,
		})
	}
	return dispatcher.EndGroups
}

func changeLog(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	newUpdate, changed := utils.CheckChanges()
	if !changed {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			Message: "You're currently running the latest version.",
			ID:      u.EffectiveMessage.ID,
		})
		return dispatcher.EndGroups
	}
	text := entityhelper.Combine("New GIGA Update", entityhelper.BoldEntity, entityhelper.UnderlineEntity)
	text.Bold("\nVersion: ").Code(newUpdate.Version)
	text.Bold("\nChange-log:")
	for _, change := range newUpdate.Changes {
		text.Plain(fmt.Sprintf("\n • %s", change))
	}
	ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		Message:  text.String,
		Entities: text.Entities,
		ID:       u.EffectiveMessage.ID,
	})
	return dispatcher.EndGroups
}
