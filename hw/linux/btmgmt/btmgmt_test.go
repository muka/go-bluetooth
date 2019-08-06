package btmgmt

import "testing"
import log "github.com/sirupsen/logrus"

func TestGetAdapters(t *testing.T) {

	log.SetLevel(log.TraceLevel)

	list, err := GetAdapters()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("At least an adapter should be available")
	}

}

func TestGetAdapter(t *testing.T) {

	_, err := GetAdapter("0")
	if err != nil {
		t.Fatal(err)
	}

}
