package yaml

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const BASE_PATH = "./data"

func LoadFile[T any](filePath string) (*T, error) {
	file, err := os.Open(path.Join(BASE_PATH, filePath))
	if err != nil {
		return nil, err
	}

	var e T
	dec := yaml.NewDecoder(file)

	return &e, dec.Decode(&e)
}

func CreateFile[T any](filePath string, data T) error {
	createdFile, err := os.Create(path.Join(BASE_PATH, filePath))
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(createdFile)
	defer enc.Close()
	enc.SetIndent(4)

	return enc.Encode(data)
}
