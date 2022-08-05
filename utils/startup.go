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
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot"
	"github.com/gigauserbot/giga/config"
	"github.com/gigauserbot/giga/db"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

var (
	TelegramClient *telegram.Client
	botSaved             = false
	BotFatherId    int64 = 93372553
)

// StartupAutomations includes the stuff to be done on each startup
func StartupAutomations(l *logger.Logger, ctx *ext.Context, client *telegram.Client) {
	if group := setupLogsGroup(ctx, client); group != 0 {
		_, err := ctx.SendMessage(group, &tg.MessagesSendMessageRequest{
			Message:      "Your GIGA is alive!",
			ReplyToMsgID: trySendingFile(ctx, group),
		})
		if err != nil {
			// check err in string because unwrapping didn't work
			if strings.Contains(err.Error(), "PEER_ID_INVALID") {
				db.UpdateLogs(0)
				StartupAutomations(l, ctx, client)
				return
			}
		}
	}
	if config.ValueOf.BotToken != "" {
		b, err := bot.MakeBot(config.ValueOf.BotToken)
		if err != nil {
			l.ChangeLevel(logger.LevelCritical).Println("failed to start bot:", err.Error())
			return
		}
		if !b.SupportsInlineQueries {
			l.ChangeLevel(logger.LevelError).Println("Inline Queries are turned off for the bot,")
			l.Println("Please enable them for full functionality of GIGA!")
		}
		bot.Username = b.Username
		bot.StartClient(l, b)
	} else if db.GetSettings().Token == "" {
		uname := setupBot(ctx, client, nil)
		if uname == "BOT_NOT_CREATED" {
			fmt.Println("failed to create bot")
			return
		}
		u, _ := ctx.ResolveUsername(uname)
		ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
			Message: "/start",
		})
	} else {
		b, err := bot.MakeBot(db.GetSettings().Token)
		if err != nil {
			l.ChangeLevel(logger.LevelCritical).Println("failed to start bot:", err.Error())
			return
		}
		if !b.SupportsInlineQueries {
			l.ChangeLevel(logger.LevelError).Println("Inline Queries are turned off for the bot,")
			l.Println("Please enable them for full functionality of GIGA!")
		}
		bot.StartClient(l, b)
	}
}

var TOKEN_REGEXP = regexp.MustCompile(`(\d+:[a-zA-Z0-9_\-]+)`)

func setupBot(ctx *ext.Context, client *telegram.Client, u types.EffectiveChat) string {
	if u == nil {
		u, _ = ctx.ResolveUsername("botfather")
	}
	BotFatherId = u.GetID()
	_, err := ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "/cancel",
	})
	if err != nil && strings.Contains(err.Error(), "YOU_BLOCKED_USER") {
		if ok, _ := ctx.Client.ContactsUnblock(ctx, &tg.InputPeerUser{
			UserID:     u.GetID(),
			AccessHash: u.GetAccessHash(),
		}); ok {
			return setupBot(ctx, client, u)

		}
	}
	time.Sleep(time.Second)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "/newbot",
	})
	time.Sleep(time.Second * 1)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "GIGA Helper Bot",
	})
	time.Sleep(time.Second * 1)
	uname := fmt.Sprintf("@GIGA_%s%dbot", string(gotgproto.Self.FirstName[0]), time.Now().Unix())
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: uname,
	})
	for i := 1; !botSaved; i++ {
		time.Sleep(time.Second * 1)
		if i >= 5 {
			return "BOT_NOT_CREATED"
		}
	}
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "/setinline",
	})
	time.Sleep(time.Second * 1)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: uname,
	})
	time.Sleep(time.Second * 1)
	ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
		Message: "giga query...",
	})
	return uname
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

func GetBotToken(l *logger.Logger) func(ctx *ext.Context, u *ext.Update) error {
	return func(ctx *ext.Context, u *ext.Update) error {
		if botSaved {
			return nil
		}
		chat := u.EffectiveChat()
		if chat.GetID() != BotFatherId {
			return nil
		}
		if !TOKEN_REGEXP.MatchString(u.EffectiveMessage.Message) {
			return nil
		}
		token := TOKEN_REGEXP.FindString(u.EffectiveMessage.Message)
		b, err := bot.MakeBot(token)
		if err != nil {
			uname := setupBot(ctx, TelegramClient, nil)
			if uname == "BOT_NOT_CREATED" {
				return nil
			}
			u, _ := ctx.ResolveUsername(uname)
			ctx.SendMessage(u.GetID(), &tg.MessagesSendMessageRequest{
				Message: "/start",
			})
		}
		bot.StartClient(l, b)
		db.UpdateBot(token)
		botSaved = true
		return nil
	}
}

func trySendingFile(ctx *ext.Context, chatId int64) int {
	upload := uploader.NewUploader(ctx.Client)
	f, err := upload.FromPath(ctx, "assets/giga.webp")
	if err != nil {
		return 0
	}
	m, err := ctx.SendMedia(chatId, &tg.MessagesSendMediaRequest{
		Media: &tg.InputMediaUploadedDocument{
			File:     f,
			MimeType: "image/webp",
		},
	})
	if err != nil {
		return 0
	}
	return m.ID
}
