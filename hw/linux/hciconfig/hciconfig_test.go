package hciconfig

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestStatus(t *testing.T) {

	h := NewHCIConfig("hci0")
	_, err := h.Status()
	if err != nil {
		t.Fatal()
	}

}

func TestUpDown(t *testing.T) {

	h := NewHCIConfig("hci0")
	_, err := h.Status()
	if err != nil {
		t.Fatal(err)
	}

	_, err = h.Down()
	if err != nil {
		t.Fatal(err)
	}

	_, err = h.Up()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetAdapters(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	_, err := GetAdapters()
	if err != nil {
		t.Fatal(err)
	}

}
