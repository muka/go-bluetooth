// Example use of the btmgmt wrapper
package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/muka/go-bluetooth/linux"
)

func main() {

	log.SetLevel(log.DebugLevel)

	list, err := linux.GetAdapters()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for i, a := range list {
		log.Infof("%d) %s (%v)", i+1, a.Name, a.Addr)
	}

}
