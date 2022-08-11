package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var tempDir = os.TempDir()

func restart(executable string, preArgs []string, delay int, chatId int64, msgId int, msgText string) {
	args := []string{fmt.Sprintf("-delay=%d", delay), fmt.Sprintf("-chat=%d", chatId), fmt.Sprintf("-msg_id=%d", msgId), fmt.Sprintf("-msg_text=%s", msgText)}
	preArgs = append(preArgs, args...)
	cmd := exec.Command(executable, preArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Start()
	cmd.Process.Release()
	os.Exit(1)
}

func Restart(delay int, chatId int64, msgId int, msgText string) {
	// // command := fmt.Sprintf("run main.go -delay=5 -chat=%d -msg=%d", chat.GetID(), u.EffectiveMessage.ID)
	// args := []string{fmt.Sprintf("-delay=%d", delay), fmt.Sprintf("-chat=%d", chatId), fmt.Sprintf("-msg_id=%d", msgId), fmt.Sprintf("-msg_text=%s", msgText)}
	// command := []string{"run", "main.go"}
	// command = append(command, args...)
	// // fmt.Println(command)

	args := []string{}
	executable, err := os.Executable()
	if strings.Contains(executable, tempDir) || err != nil {
		executable = "go"
		args = []string{"run", "main.go"}
	}
	restart(executable, args, delay, chatId, msgId, msgText)
	// cmd := exec.Command(executable, command...)
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// cmd.Stdin = os.Stdin
	// cmd.Start()
	// cmd.Process.Release()
	// os.Exit(1)
}
