package profile

import (
	profileGen "github.com/muka/go-bluetooth/src/gen/profile/profile"
)

// NewProfileManager1 create a new ProfileManager1 client
func NewProfileManager1(hostID string) (*profileGen.ProfileManager1, error) {
	return profileGen.NewProfileManager1()
}
