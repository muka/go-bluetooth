package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func failArg(arg string) {
	failArgs([]string{arg})
}

func failArgs(args []string) {
	fail(fmt.Errorf("Missing arguments: %s", args))
}

func fail(err error) {
	if err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
