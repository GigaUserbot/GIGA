package modules

import (
	"os"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gotd/td/tg"
)

func (m *module) LoadProg(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("PROG")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	dispatcher.AddHandler(handlers.NewCommand("killub", authorised(killub)))
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
