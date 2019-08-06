package adapter

import (
	"fmt"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	a := getDefaultAdapter(t)

	err := a.StartDiscovery()
	if err != nil {
		t.Fatal(err)
	}

	discovery, cancel, err := a.DeviceDiscovered()
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
		sleep := 30
		time.Sleep(time.Duration(sleep) * time.Second)
		wait <- fmt.Errorf("Discovery timeout exceeded (%ds)", sleep)
	}()

	err = <-wait
	if err != nil {
		t.Fatal(err)
	}

}
