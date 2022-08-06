package bot

import (
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/anonyindian/logger"
)

var Username string

func MakeBot(token string) (*gotgbot.Bot, error) {
	return gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
}

func StartClient(l *logger.Logger, b *gotgbot.Bot) {
	log := l.Create("BOT")
	// custom dispatcher handles all the updates
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		args := ctx.Args()
		if len(args) < 2 {
			ctx.EffectiveMessage.Reply(b, "Started", nil)
			return ext.EndGroups
		}
		switch args[1] {
		case "deploy_own_via_help":
			ctx.EffectiveMessage.Reply(b, "Hey! It seems like you were trying use help section of GIGA deployed by someone else, unfortunately, we cannot allow you to use their userbot's help because we respect privacy of our users."+
				"\n\nNo Worries Though! Deploy your own private GIGA using the following button and use all the features:", &gotgbot.SendMessageOpts{
				ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{{Text: "Deploy GIGA", Url: "https://github.com/GigaUserbot/GIGA"}},
					},
				},
			})
		}
		return ext.EndGroups
	}))

	Load(log, dispatcher)

	updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	log.Println("STARTED")
}
