package linux

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// CmdExec Execute a command
func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	path, err := exec.LookPath(baseCmd)
	if err != nil {
		return "", err
	}

	log.Tracef("Exec: %s %s", path, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs[0])
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
