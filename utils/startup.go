package utils

import (
	"context"

	"github.com/gotd/td/telegram"
)

var TelegramClient *telegram.Client

// StartupAutomations includes the stuff to be done on each startup
func StartupAutomations(ctx context.Context, client *telegram.Client) {
	// user, _ := client.Self(ctx)
	// gotgproto.Api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
	// 	Peer: &tg.InputPeerUser{
	// 		UserID:     user.ID,
	// 		AccessHash: user.AccessHash,
	// 	},
	// 	Message:  "Giga is up!",
	// 	RandomID: time.Now().Unix(),
	// })
}
