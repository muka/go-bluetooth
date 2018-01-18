// Example use of the btmgmt wrapper
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/linux"
)

func main() {

	log.SetLevel(log.DebugLevel)

	list, err := linux.GetAdapters()
	if err != nil {
		panic(err)
	}

	for i, a := range list {
		log.Infof("%d) %s (%v)", i+1, a.Name, a.Addr)
	}

}
