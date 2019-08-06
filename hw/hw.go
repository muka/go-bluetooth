package hw

import "github.com/muka/go-bluetooth/hw/linux"

func GetAdapter(adapterID string) (a linux.AdapterInfo, err error) {
	return linux.GetAdapter(adapterID)
}

func GetAdapters() ([]linux.AdapterInfo, error) {
	return linux.GetAdapters()
}

func Up(adapterID string) error {
	return linux.Up(adapterID)
}

func Down(adapterID string) error {
	return linux.Down(adapterID)
}

func Reset(adapterID string) error {
	return linux.Reset(adapterID)
}
