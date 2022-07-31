package modules

import (
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
)

func (m *module) LoadAdmin(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("ADMIN")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	dispatcher.AddHandler(handlers.NewCommand("ban", authorised(ban)))
}

func ban(ctx *ext.Context, u *ext.Update) error {
	// TODO: make it working
	return dispatcher.EndGroups
}
