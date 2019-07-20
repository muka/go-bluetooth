package zoeetree

import "testing"

var address = "F3:AC:E9:41:7B:AE"

func TestNewZooeTree(t *testing.T) {

	zt := NewZoeeTree(address)

	err := zt.Connect(10)
	if err != nil {
		t.Fatal(err)
	}
}
