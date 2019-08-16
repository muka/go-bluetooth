package api

import (
	"fmt"
	"testing"

	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/util"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	IgnoreFlag        bool                   `dbus:"ignore"`
	ToOmit            map[string]interface{} `dbus:"omitEmpty"`
	Ignored           string                 `dbus:"ignore"`
	IgnoredByProperty []string               `dbus:"ignore=IgnoreFlag"`
}

func (s testStruct) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := util.StructToMap(s, m)
	return m, err
}

func TestParseTag(t *testing.T) {

	s := testStruct{
		IgnoreFlag:        true,
		Ignored:           "foo",
		IgnoredByProperty: []string{"bar"},
	}

	prop := &DBusProperties{
		props:       make(map[string]bluez.Properties),
		propsConfig: make(map[string]map[string]*prop.Prop),
	}

	prop.AddProperties("test", s)

	err := prop.parseProperties()
	if err != nil {
		t.Fatal(err)
	}

	cfg := prop.propsConfig["test"]

	for field, cfg := range cfg {
		fmt.Printf("%s: %++v\n", field, cfg)
	}

	_, ok := cfg["ToOmit"]
	assert.True(t, !ok)

	_, ok = cfg["Ignored"]
	assert.True(t, !ok)

	_, ok = cfg["IgnoredByProperty"]
	assert.True(t, !ok)

}
