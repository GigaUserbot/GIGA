package main

import (
	"context"
	"os"

	"github.com/anonyindian/giga/config"
	"github.com/anonyindian/giga/modules"
	"github.com/anonyindian/giga/utils"
	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/anonyindian/logger"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
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
	runClient(l)
}

func runClient(l *logger.Logger) {
	// custom dispatcher handles all the updates
	dp := dispatcher.MakeDispatcher()
	modules.Load(l, dp)
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
				utils.TelegramClient = client
				utils.StartupAutomations(ctx, client)
			}()
			return nil
		},
		// Uncomment DCList to run Giga on test servers
		DCList: dcs.Test(),
	})
}
