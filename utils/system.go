package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var tempDir = os.TempDir()

func Restart(delay int, chatId int64, msgId int, msgText string) {
	// command := fmt.Sprintf("run main.go -delay=5 -chat=%d -msg=%d", chat.GetID(), u.EffectiveMessage.ID)
	args := []string{fmt.Sprintf("-delay=%d", delay), fmt.Sprintf("-chat=%d", chatId), fmt.Sprintf("-msg_id=%d", msgId), fmt.Sprintf("-msg_text=%s", msgText)}
	command := []string{"run", "main.go"}
	command = append(command, args...)
	// fmt.Println(command)
	executable, err := os.Executable()
	if strings.Contains(executable, tempDir) || err != nil {
		executable = "go"
	}
	cmd := exec.Command(executable, command...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Start()
	cmd.Process.Release()
	os.Exit(1)
}
