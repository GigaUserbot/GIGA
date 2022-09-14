package modules

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
	"github.com/gotd/td/tg"
)

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func (m *module) LoadShell(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("SHELL")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("shell", "help")
	dispatcher.AddHandler(handlers.NewCommand("sh", sh))
}

func sh(ctx *ext.Context, u *ext.Update) error {
	chat := u.EffectiveChat()
	// msg := u.EffectiveMessage
	cmd := strings.Fields(u.EffectiveMessage.Message)
	if len(cmd) == 1 {
		_, err := ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
			ID:      u.EffectiveMessage.ID,
			Message: "No command passed",
		})
		return err
	}

	err, out, errout := Shellout(strings.Join(cmd[1:], " "))
	if err != nil {
		logger.Println("error: %v\n", err)
	}
	m := ""
	if out != "" {
		m = "Standard output: \n" + out
	}

	if errout != "" {
		if m != "" {
			m += "\n"
		}
		m += "Standard Output Error: \n" + errout
	}

	_, err = ctx.EditMessage(chat.GetID(), &tg.MessagesEditMessageRequest{
		Message: m,
		ID:      u.EffectiveMessage.ID,
	})
	return err
}
