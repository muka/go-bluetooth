package cmd

import (
	"fmt"
	"log"
	"os"
)

func failArg(arg string) {
	failArgs([]string{arg})
}

func failArgs(args []string) {
	fail(fmt.Errorf("Missing arguments: %s", args))
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
