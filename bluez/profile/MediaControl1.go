package profile

import (
	"fmt"

	"github.com/muka/go-bluetooth/src/gen/profile/media"
)

// NewMediaControl1 create a new MediaControl1 client
func NewMediaControl1(hostID string) (*media.MediaControl1, error) {
	return media.NewMediaControl1(fmt.Sprintf("/org/bluez/%s", hostID))
}
