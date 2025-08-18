package settings

import (
	"os"
	"path"
	"yotudo/src/lib/logger"

	"gopkg.in/yaml.v3"
)

func LoadYaml[T any](filePath string) (*T, error) {
	file, err := os.Open(path.Join("./data", filePath))
	if err != nil {
		logger.WarningF("Couldn't read yaml file (%s) due to: %s", filePath, err.Error())

		return nil, err
	}

	logger.DebugF("Open '%s' and pass to yaml decoder", filePath)

	var e T
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&e); err != nil {
		logger.Error(err)

		return nil, err
	}

	return &e, nil
}

func CreateYaml[T any](filePath string, data T) error {
	logger.Debug("CreateYaml was called")
	createdFile, err := os.Create(path.Join("./data", filePath))
	if err != nil {
		return err
	}

	logger.DebugF("Created '%s' and passed to yaml encoder", filePath)

	enc := yaml.NewEncoder(createdFile)
	enc.SetIndent(4)
	defer enc.Close()
	if err := enc.Encode(data); err != nil {
		logger.Error(err)

		return err
	}

	return nil
}
