package gen

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	api, err := Parse("../src/bluez/doc", []string{})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, api.Version)
	assert.NotEmpty(t, api.Api)
}

func TestSerialization(t *testing.T) {

	api, err := Parse("../src/bluez/doc", []string{})
	if err != nil {
		t.Fatal(err)
	}

	destDir := "../test/"
	jsonFile := path.Join(destDir, "test.json")

	err = Mkdir(destDir)
	if err != nil {
		t.Fatal(err)
	}

	err = api.Serialize(jsonFile)
	if err != nil {
		t.Fatal(err)
	}

	api1, err := LoadJSON(jsonFile)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, api.Version, api1.Version)
}
