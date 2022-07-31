package main

import (
	"context"
	"os"

	"github.com/anonyindian/giga/config"
	"github.com/anonyindian/giga/modules"
	"github.com/anonyindian/giga/sql"
	"github.com/anonyindian/giga/utils"
	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/anonyindian/logger"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	l := logger.New(os.Stderr, &logger.LoggerOpts{
		ProjectName: "GIGA-USERBOT",
	})
	if config.DEBUG {
		l.ChangeMinimumLevel(logger.LevelDebug)
	}
	config.Load(l)
	handlers.DefaultPrefix = []rune{'.', '$'}
	sql.Load(l)
	runClient(l)
}

func runClient(l *logger.Logger) {
	// custom dispatcher handles all the updates
	dp := dispatcher.MakeDispatcher()
	gotgproto.StartClient(gotgproto.ClientHelper{
		// Get AppID from https://my.telegram.org/apps
		AppID: config.ValueOf.AppId,
		// Get ApiHash from https://my.telegram.org/apps
		ApiHash: config.ValueOf.ApiHash,
		// Session of your client
		// sessionName: name of the session / session string in case of TelethonSession or StringSession
		// sessionType: can be any out of Session, TelethonSession, StringSession.
		Session: sessionMaker.NewSession(config.ValueOf.SessionString, sessionMaker.TelethonSession),
		// Make sure to specify custom dispatcher here in order to enjoy gotgproto's update handling
		Dispatcher: dp,
		// Add the handlers, post functions in TaskFunc
		TaskFunc: func(ctx context.Context, client *telegram.Client) error {
			go func() {
				for {
					if gotgproto.Sender != nil {
						l.Println("GIGA HAS BEEN STARTED")
						break
					}
				}
				ctx := ext.NewContext(ctx, client.API(), gotgproto.Self, gotgproto.Sender, &tg.Entities{})
				utils.TelegramClient = client
				utils.StartupAutomations(ctx, client)
				// Modules shall not be loaded unless the setup is complete
				modules.Load(l, dp)
			}()
			return nil
		},
		// Uncomment DCList to run Giga on test servers
		// DCList: dcs.Test(),
	})
}
