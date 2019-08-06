package hcitool

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHcitoolGetAdapters(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	list, err := GetAdapters()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, list)

}

func TestHcitoolGetAdapter(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	a, err := GetAdapter("hci0")
	if err != nil {
		t.Fatal(err)
	}

	if a == nil {
		t.Fatal("An adapter should be available")
	}

}

func TestHcitoolGetAdapterNotfound(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	list, err := GetAdapters()
	if err != nil {
		t.Fatal(err)
	}

	size := len(list)
	devID := fmt.Sprintf("hci%d", (size + 1))
	a, err := GetAdapter(devID)
	if err != nil {
		t.Fatal(err)
	}

	if a != nil {
		t.Fatal(fmt.Sprintf("%s should not be avail", devID))
	}
}
