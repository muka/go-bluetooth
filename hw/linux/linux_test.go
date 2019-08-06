package linux

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAdapters(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	list, err := GetAdapters()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, list)
}

func TestGetAdapter(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	_, err := GetAdapter("hci0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAdapterNotFound(t *testing.T) {
	_, err := GetAdapter("hci999")
	if err == nil {
		t.Fatal("adapter should not exists")
	}
}

func TestUp(t *testing.T) {
	err := Up("hci0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDown(t *testing.T) {
	err := Down("hci0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestReset(t *testing.T) {
	err := Reset("hci0")
	if err != nil {
		t.Fatal(err)
	}
}
