package modules

import (
	"fmt"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gotd/td/tg"
)

func (m *module) LoadHelp(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("HELP")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetMainHelp("<b><u>The GIGA Userbot</u></b>\n\nHere is the help menu:", "html")
	dispatcher.AddHandler(handlers.NewCommand("help", authorised(helpCmd)))
}

func helpCmd(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	result, err := ctx.GetInlineBotResults(chat.GetID(), bot.Username, &tg.MessagesGetInlineBotResultsRequest{
		Query: "help",
	})
	if err != nil {
		ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: fmt.Sprintf("failed to load help: %s", err.Error()),
		})
		return dispatcher.EndGroups
	}
	ctx.SendInlineBotResult(chat.GetID(), &tg.MessagesSendInlineBotResultRequest{
		HideVia: true,
		QueryID: result.QueryID,
		ID:      result.Results[0].GetID(),
	})
	switch {
	case chat.IsAChannel():
		ctx.Client.ChannelsDeleteMessages(ctx, &tg.ChannelsDeleteMessagesRequest{
			Channel: &tg.InputChannel{
				ChannelID:  chat.GetID(),
				AccessHash: chat.GetAccessHash(),
			},
			ID: []int{u.EffectiveMessage.ID},
		})
	default:
		ctx.Client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     []int{u.EffectiveMessage.ID},
		})
	}
	return dispatcher.EndGroups
}
