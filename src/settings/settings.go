package settings

import "yotudo/src/lib/logger"

var Global Settings

func LoadSettings() (*Settings, error) {
	Global = Settings{
		App: AppSettings{
			TempLocation:   "./data/tmp",
			ImagesLocation: "./data/imgs",
			MusicsLocation: "./data/mscs",
		},
		Database: DatabaseSettings{
			Location: "./data/agd_01",
			Version:  "0.1.0",
		},
	}

	if config, err := LoadYaml[Settings]("config.yaml"); err != nil {
		logger.Error(err)

		return nil, err
	} else if err := Global.Merge(config, true); err != nil {
		logger.Error(err)

		return nil, err
	}

	return &Global, nil
}
