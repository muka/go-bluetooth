package hci

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestHciList(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	list, err := List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("At least an adapter should be available")
	}

}

func TestHciUp(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	list, err := List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("At least an adapter should be available")
	}

	Up(list[0])
	Down(list[0])

}
