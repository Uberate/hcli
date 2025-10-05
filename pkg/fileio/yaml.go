package fileio

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYaml(path string, obj interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, obj)
}
