package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/linux/hci"
)

func main() {

	log.SetLevel(log.DebugLevel)

	err := hci.Do()
	if err != nil {
		panic(err)
	}

}
