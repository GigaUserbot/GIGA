package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/generic"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gigauserbot/giga/config"
	"github.com/gigauserbot/giga/db"
	"github.com/gigauserbot/giga/modules"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"
)

var (
	delay          = flag.Int("delay", 0, "")
	restartChatId  = flag.Int("chat", 0, "")
	restartMsgId   = flag.Int("msg_id", 0, "")
	restartMsgText = flag.String("msg_text", "", "")
)

func main() {
	flag.Parse()
	l := logger.New(os.Stderr, &logger.LoggerOpts{
		ProjectName: "GIGA-USERBOT",
	})
	if *restartMsgId != 0 {
		// Clean Console
		os.Stderr.Write([]byte("\n"))
	}
	if *delay != 0 {
		l.Println("Delaying start for", *delay, "seconds")
		time.Sleep(time.Second * time.Duration(*delay))
	}
	if config.DEBUG {
		l.ChangeMinimumLevel(logger.LevelDebug)
	}
	config.Load(l)
	handlers.DefaultPrefix = []rune{'.', '$'}
	db.Load(l)
	runClient(l)
}

func runClient(l *logger.Logger) {
	log := l.Create("CLIENT")
	// custom dispatcher handles all the updates
	dp := dispatcher.MakeDispatcher()
	dp.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, utils.GetBotToken(l)), 2)
	gotgproto.StartClient(&gotgproto.ClientHelper{
		// Get AppID from https://my.telegram.org/apps
		AppID: config.ValueOf.AppId,
		// Get ApiHash from https://my.telegram.org/apps
		ApiHash: config.ValueOf.ApiHash,
		// Session of your client
		// sessionName: name of the session / session string in case of TelethonSession or StringSession
		// sessionType: can be any out of Session, TelethonSession, StringSession.
		Session: sessionMaker.NewSession(config.GetSessionString(), sessionMaker.TelethonSession),
		// Make sure to specify custom dispatcher here in order to enjoy gotgproto's update handling
		Dispatcher: dp,
		// Add the handlers, post functions in TaskFunc
		TaskFunc: func(ctx context.Context, client *telegram.Client) error {
			go func() {
				for {
					if gotgproto.Sender != nil {
						log.ChangeLevel(logger.LevelInfo).Println("STARTED")
						break
					}
				}
				ctx := ext.NewContext(ctx, client.API(), gotgproto.Self, gotgproto.Sender, &tg.Entities{})
				utils.TelegramClient = client
				if *restartMsgId == 0 && *restartMsgText == "" {
					utils.StartupAutomations(l, ctx, client)
				} else {
					generic.EditMessage(ctx, *restartChatId, &tg.MessagesEditMessageRequest{
						ID:      *restartMsgId,
						Message: *restartMsgText,
					})
				}
				// Modules shall not be loaded unless the setup is complete
				modules.Load(l, dp)
				helpmaker.MakeHelp()
				if config.ValueOf.TestServer {
					l.ChangeLevel(logger.LevelMain).Println("RUNNING ON TEST SERVER")
				}
				l.ChangeLevel(logger.LevelMain).Println("GIGA HAS BEEN STARTED")
			}()
			return nil
		},
		DisableCopyright: true,
		DCList: func() (dct dcs.List) {
			if config.ValueOf.TestServer {
				dct = dcs.Test()
			}
			return
		}(),
	})
}
