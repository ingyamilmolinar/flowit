package config

import (
	"os/exec"
	"strings"
)

// TODO: Validate command and args before executing
func shellValidator(shell string) bool {
	cmds := strings.Split(shell, " ")
	/* #gosec */
	cmd := exec.Command(cmds[0], cmds[1:]...)
	_, err := cmd.Output()
	return err == nil
}
