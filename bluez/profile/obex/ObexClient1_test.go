package obex

import (
	"testing"
)

func TestNewObexClient1(t *testing.T) {
	t.Log("Create ObexClient1")

	a := NewObexClient1()

	t.Log("Start CreateSession")

	temp := map[string]interface{}{}
	//temp := map[string]dbus.Variant{}
	temp["Target"] = "opp"

	_, err := a.CreateSession("98:4E:97:00:3F:3C", temp)
	if err != nil {
		t.Logf("Error on CreateSession")
		t.Fatal(err)
	}
}
