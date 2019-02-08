package linux

import "github.com/muka/ble/linux/hci/socket"

//Down turn down an HCI device
func Down(adapterID int) error {
	return socket.Down(adapterID)
}

//Up turn up an HCI device
func Up(adapterID int) error {
	return socket.Up(adapterID)
}

//List list available HCI devices
func List() ([]int, error) {
	return socket.List()
}
