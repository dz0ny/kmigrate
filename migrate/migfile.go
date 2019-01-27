package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	createjsonpatch "github.com/mattbaird/jsonpatch"
)

type Config struct {
	filename  string
	Version   string            `json:"version"`
	Kind      string            `json:"kind"`
	Selectors map[string]string `json:"selectors"`
	Patch     interface{}       `json:"patch"`
}

func (c *Config) GetLabelSelector() string {
	return c.Selectors["label"]
}

func (c *Config) GetResource() string {
	return strings.ToLower(c.Kind)
}

func (c *Config) GetAPIPath() string {
	// https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-groups
	if c.Version == "v1" {
		return "/api/v1/"
	}
	return fmt.Sprintf("/apis/%s/", strings.ToLower(c.Version))
}

func (c *Config) UpdatePatch(i []createjsonpatch.JsonPatchOperation) error {
	c.Patch = i
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.filename, data, 0x644)
}

func (c *Config) GetPatch() ([]byte, error) {
	return json.Marshal(c.Patch)
}

// NewConfigFromFile returns a Config object from filepath
func NewConfigFromFile(filename string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	c.filename = filename
	return &c, nil
}
