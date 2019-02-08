package linux

import "testing"
import log "github.com/sirupsen/logrus"

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
