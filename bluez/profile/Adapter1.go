package profile

import (
	"fmt"

	"github.com/muka/go-bluetooth/src/gen/profile/adapter"
)

// NewAdapter1 create a new Adapter1 client
func NewAdapter1(hostID string) (*adapter.Adapter1, error) {
	objectPath := fmt.Sprintf("/org/bluez/%s", hostID)
	a, err := adapter.NewAdapter1(objectPath)
	return a, err
}
