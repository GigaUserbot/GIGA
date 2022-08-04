package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/anonyindian/gotgproto/types"
	"github.com/gigauserbot/giga/db"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

var (
	TelegramClient *telegram.Client
	BotSaved       = false
)

// StartupAutomations includes the stuff to be done on each startup
func StartupAutomations(ctx *ext.Context, client *telegram.Client) {
	if group := setupLogsGroup(ctx, client); group != 0 {
		_, err := ctx.SendMessage(group, &tg.MessagesSendMessageRequest{
			Message: "Your GIGA is alive!",
		})
		if err != nil {
			// check err in string because unwrapping didn't work
			if strings.Contains(err.Error(), "PEER_ID_INVALID") {
				db.UpdateLogs(0)
				StartupAutomations(ctx, client)
				return
			}
		}
	}
	_ = setupBot
	// setupBot(ctx, client, nil)
}

var TOKEN_REGEXP = regexp.MustCompile(`(\d+:[a-zA-Z0-9_\-]+)`)

func setupBot(ctx *ext.Context, client *telegram.Client, u types.EffectiveChat) int64 {
	if u == nil {
		u, _ = ctx.ResolveUsername("botfather")
	}
	_, err := ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "/newbot",
	})
	if err != nil && strings.Contains(err.Error(), "YOU_BLOCKED_USER") {
		if ok, _ := ctx.Client.ContactsUnblock(ctx, &tg.InputPeerUser{
			UserID:     u.GetID(),
			AccessHash: u.GetAccessHash(),
		}); ok {
			return setupBot(ctx, client, u)

		}
	}
	time.Sleep(time.Second * 1)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "GIGA Helper Bot",
	})
	time.Sleep(time.Second * 1)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: fmt.Sprintf("@GIGA_%s%dbot", string(gotgproto.Self.FirstName[0]), time.Now().Unix()),
	})
	for !BotSaved {
		time.Sleep(time.Second * 1)
	}

	return 0
}

func setupLogsGroup(ctx *ext.Context, client *telegram.Client) int64 {
	if group := db.GetSettings().LogsGroup; group != 0 {
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
	db.UpdateLogs(group)
	return group
}
