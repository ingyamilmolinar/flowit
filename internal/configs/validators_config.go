package configs

import (
	"os/exec"
	"strings"
)

func shellValidator(shell string) bool {
	cmds := strings.Split(shell, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	_, err := cmd.Output()

	if err != nil {
		return false
	}
	return true
}
