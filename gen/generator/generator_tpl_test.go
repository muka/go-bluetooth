package generator

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestCastType(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	typedef := "object"
	res := castType(typedef)

	if res != "dbus.ObjectPath" {
		t.Fatal(fmt.Sprintf("%s != %s", typedef, res))
	}

	typedef = "array{objects, properties}"
	res = castType(typedef)

	if res != "[]dbus.ObjectPath, string" {
		t.Fatal(fmt.Sprintf("%s != %s", typedef, res))
	}

}
