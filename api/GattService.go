package api

// NewGattService creates a new GATT Service
func NewGattService(path string) *GattService {
	s := GattService{path}
	return &s
}

//GattService a GATT service for a Device
type GattService struct {
	Path string
}
