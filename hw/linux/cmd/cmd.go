package cmd

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Exec Execute a command and collect the output
func Exec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Tracef("Exec: %s %s", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)
	res, err := cmd.CombinedOutput()

	return string(res), err
}
