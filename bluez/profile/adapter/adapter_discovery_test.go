package adapter

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestDiscovery(t *testing.T) {
	a := getDefaultAdapter(t)

	err := a.StartDiscovery()
	if err != nil {
		t.Fatal(err)
	}
	defer a.StopDiscovery()

	discovery, cancel, err := a.OnDeviceDiscovered()
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	wait := make(chan error)

	go func() {
		for dev := range discovery {
			if dev == nil {
				return
			}
			wait <- nil
		}
	}()

	go func() {
		sleep := 5
		time.Sleep(time.Duration(sleep) * time.Second)
		log.Debugf("Discovery timeout exceeded (%ds)", sleep)
		wait <- nil
	}()

	err = <-wait
	if err != nil {
		t.Fatal(err)
	}

}
