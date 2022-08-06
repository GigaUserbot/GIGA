package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
)

func helpInline(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery
	data := strings.Fields(query.Query)
	if len(data) == 1 {
		b.AnswerInlineQuery(query.Id, []gotgbot.InlineQueryResult{
			&gotgbot.InlineQueryResultArticle{
				Id:    query.Id,
				Title: "Help Menu",
				InputMessageContent: &gotgbot.InputTextMessageContent{
					MessageText: helpmaker.GetMainHelp(),
					ParseMode:   helpmaker.GetParseMode(),
				},
				ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: helpmaker.GetPageHelp(1),
				},
			},
		}, &gotgbot.AnswerInlineQueryOpts{
			CacheTime: 1,
		})
		return ext.EndGroups
	}
	return ext.EndGroups
}

func helpCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	if query.From.Id != gotgproto.Self.ID {
		query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Url: fmt.Sprintf("t.me/%s?start=deploy_own_via_help", strings.TrimPrefix(b.Username, "@")),
		})
		return ext.EndGroups
	}
	go query.Answer(b, nil)
	if query.Data == "help_" {
		b.EditMessageText(helpmaker.GetMainHelp(), &gotgbot.EditMessageTextOpts{
			InlineMessageId: query.InlineMessageId,
			ParseMode:       helpmaker.GetParseMode(),
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: helpmaker.GetPageHelp(1),
			},
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}
	data := strings.Split(query.Data, "_")
	switch data[1] {
	case "next":
		// let it panic if index error
		currPage, _ := strconv.Atoi(data[2])
		b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
			InlineMessageId: query.InlineMessageId,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: helpmaker.GetPageHelp(currPage + 1),
			},
		})
	case "prev":
		// let it panic if index error
		currPage, _ := strconv.Atoi(data[2])
		b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
			InlineMessageId: query.InlineMessageId,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: helpmaker.GetPageHelp(currPage - 1),
			},
		})
	case "close":
		b.EditMessageText("Help Menu Closed", &gotgbot.EditMessageTextOpts{
			InlineMessageId:       query.InlineMessageId,
			DisableWebPagePreview: true,
		})
	default:
		b.EditMessageText(helpmaker.GetModuleHelp(data[1]), &gotgbot.EditMessageTextOpts{
			InlineMessageId:       query.InlineMessageId,
			ParseMode:             helpmaker.GetParseMode(),
			DisableWebPagePreview: true,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Back", CallbackData: helpmaker.HelpPrifex}, {Text: "Close", CallbackData: helpmaker.HelpPrifex + "close"}},
				},
			},
		})
	}
	return ext.EndGroups
}

func (m *module) LoadInline(dp *ext.Dispatcher) {
	log := m.Logger.Create("INLINE")
	defer log.ChangeLevel(logger.LevelMain).Println("LOADED")
	dp.AddHandler(handlers.NewInlineQuery(inlinequery.QueryPrefix("help"), helpInline))
	dp.AddHandler(handlers.NewCallback(callbackquery.Prefix("help_"), helpCallback))
}
