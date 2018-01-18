package linux

import "testing"

func TestGetAdapters(t *testing.T) {

	list, err := GetAdapters()
	if err != nil {
		t.FailNow()
	}

	if len(list) == 0 {
		t.Fatal("At least an adapter should be available")
	}

}
