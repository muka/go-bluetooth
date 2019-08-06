package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoveryFilter(t *testing.T) {

	f := NewDiscoveryFilter()
	f.AddUUIDs("AAAA", "BBBB")
	f.RSSI = 10
	f.Pathloss = 1
	f.DuplicateData = false
	f.Transport = DiscoveryFilterTransportLE

	m := f.ToMap()

	assert.EqualValues(t, m["UUIDs"].([]string), f.UUIDs)
	assert.EqualValues(t, m["RSSI"].(int16), f.RSSI)
	assert.EqualValues(t, m["Pathloss"].(uint16), f.Pathloss)
	assert.EqualValues(t, m["DuplicateData"].(bool), f.DuplicateData)
	assert.EqualValues(t, m["Transport"].(string), f.Transport)

}
