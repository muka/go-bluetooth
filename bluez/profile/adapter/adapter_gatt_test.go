package adapter

import "testing"

//GetGattManager return a GattManager1 instance
func TestGetGattManager(t *testing.T) {

	a := getDefaultAdapter(t)

	_, err := a.GetGattManager()
	if err != nil {
		t.Fatal(err)
	}

}
