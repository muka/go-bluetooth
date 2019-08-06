package api

import (
	"fmt"
	"testing"
	"time"
)

func TestDiscoverDevice(t *testing.T) {

	adapterID := GetDefaultAdapterID()

	// err := ResetController(adapterID)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	discovery, cancel, err := Discover(adapterID, nil)
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
