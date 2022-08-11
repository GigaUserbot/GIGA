package modules

import (
	"fmt"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/parsemode/entityhelper"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/tg"
)

func (m *module) LoadUpdate(dp *dispatcher.CustomDispatcher) {
	dp.AddHandler(handlers.NewCommand("changelog", changeLog))
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
