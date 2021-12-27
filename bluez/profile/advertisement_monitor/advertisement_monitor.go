package advertisement_monitor

// array{(uint8, uint8, array{byte})}
type Pattern struct {
	v1 uint8
	v2 uint8
	v3 []byte
}

type Patterns = []Pattern
