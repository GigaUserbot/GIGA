package utils

import (
	"strings"

	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gigauserbot/giga/sql"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

var TelegramClient *telegram.Client

// StartupAutomations includes the stuff to be done on each startup
func StartupAutomations(ctx *ext.Context, client *telegram.Client) {
	if group := setupLogsGroup(ctx, client); group != 0 {
		_, err := ctx.SendMessage(group, &tg.MessagesSendMessageRequest{
			Message: "Your GIGA is alive!",
		})
		if err != nil {
			// check err in string because unwrapping didn't work
			if strings.Contains(err.Error(), "PEER_ID_INVALID") {
				sql.UpdateSettings(0)
				StartupAutomations(ctx, client)
			}
		}
	}
}

func setupLogsGroup(ctx *ext.Context, client *telegram.Client) int64 {
	if group := sql.GetSettings().LogsGroup; group != 0 {
		return group
	}
	u, _ := ctx.ResolveUsername("GIGAubot")
	upd, _ := client.API().MessagesCreateChat(ctx, &tg.MessagesCreateChatRequest{
		Users: []tg.InputUserClass{&tg.InputUser{
			UserID:     u.GetID(),
			AccessHash: u.GetAccessHash(),
		}},
		Title: "GIGA Userbot Logs",
	})
	update, ok := upd.(*tg.Updates)
	if !ok {
		return 0
	}
	group := update.Chats[0].GetID()
	// Add created group's peer to storage coz gotgproto still doesn't do that :P
	storage.AddPeer(group, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
	sql.UpdateSettings(group)
	return group
}
