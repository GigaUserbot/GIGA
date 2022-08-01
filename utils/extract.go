package utils

import (
	"errors"
	"strings"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/types"
	"github.com/gotd/td/tg"
)

func ExtractUser(ctx *ext.Context, msg *tg.Message, chat types.EffectiveChat) (target int64, err error) {
	if id := msg.ReplyTo.ReplyToMsgID; id != 0 {
		var m []tg.MessageClass
		if _, ok := chat.(*types.Chat); ok {
			m, err = ctx.GetMessages([]tg.InputMessageClass{&tg.InputMessageID{
				ID: id,
			}})
		} else {
			var ms tg.MessagesMessagesClass
			ms, err = ctx.Client.ChannelsGetMessages(ctx, &tg.ChannelsGetMessagesRequest{
				Channel: &tg.InputChannel{
					ChannelID:  chat.GetID(),
					AccessHash: chat.GetAccessHash(),
				},
				ID: []tg.InputMessageClass{&tg.InputMessageID{
					ID: id,
				}},
			})
			switch mt := ms.(type) {
			case *tg.MessagesMessages:
				m = mt.Messages
			case *tg.MessagesChannelMessages:
				m = mt.Messages
			}
		}
		if err != nil {
			return
		}
		msg, ok := m[0].(*tg.Message)
		if ok {
			target = msg.FromID.(*tg.PeerUser).UserID
		}
	}
	if target == 0 {
		args := strings.Fields(msg.Message)
		if !(len(args) > 1 && strings.HasPrefix(args[1], "@")) {
			err = errors.New("no user provided")
			return
		}
		msg.MapEntities()
		var c types.EffectiveChat
		c, err = ctx.ResolveUsername(args[1])
		if err != nil {
			return
		}
		target = c.GetID()
	}
	return
}
