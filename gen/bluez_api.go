package gen

import (
	"encoding/json"
	"io/ioutil"

	"github.com/muka/go-bluetooth/gen/types"
)

type BluezAPI struct {
	Version string
	Api     []*types.ApiGroup
}

// Serialize store the structure as JSON
func (g *BluezAPI) Serialize(destFile string) error {

	data, err := json.Marshal(g)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destFile, data, 0755)
}

// LoadJSON parse an ApiGroup from JSON definition
func LoadJSON(srcFile string) (*BluezAPI, error) {

	b, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return nil, err
	}

	a := new(BluezAPI)
	err = json.Unmarshal(b, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
