package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/config"
	"github.com/gigauserbot/giga/db"
	"github.com/gigauserbot/giga/modules"
	"github.com/gigauserbot/giga/utils"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
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
	db.Load(l)
	runClient(l)
}

const BOT_FATHER_ID = 93372553

func runClient(l *logger.Logger) {
	log := l.Create("CLIENT")
	// custom dispatcher handles all the updates
	dp := dispatcher.MakeDispatcher()
	dp.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, func(ctx *ext.Context, u *ext.Update) error {
		chat := u.EffectiveChat()
		if chat.GetID() != BOT_FATHER_ID {
			return nil
		}
		if !utils.TOKEN_REGEXP.MatchString(u.EffectiveMessage.Message) {
			return nil
		}
		token := utils.TOKEN_REGEXP.FindString(u.EffectiveMessage.Message)
		// let it panic if err arise
		resp, _ := http.Get("api.telegram.org/bot%s/getMe")
		b, _ := ioutil.ReadAll(resp.Body)
		getme := struct {
			UserId int64 `json:"user_id"`
		}{}
		if json.Unmarshal(b, &getme) != nil {
			return nil
		}
		db.SetBot(token)
		utils.BotSaved = true
		return nil
	}), -69)
	gotgproto.StartClient(&gotgproto.ClientHelper{
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
						log.ChangeLevel(logger.LevelInfo).Println("STARTED")
						break
					}
				}
				ctx := ext.NewContext(ctx, client.API(), gotgproto.Self, gotgproto.Sender, &tg.Entities{})
				utils.TelegramClient = client

				utils.StartupAutomations(ctx, client)
				// Modules shall not be loaded unless the setup is complete
				modules.Load(l, dp)
				l.Println("GIGA HAS BEEN STARTED")
			}()
			return nil
		},
		DCList: func() (dct dcs.List) {
			if config.ValueOf.TestServer {
				dct = dcs.Test()
			}
			return
		}(),
	})
}
