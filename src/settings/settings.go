package settings

import (
	"os"
	"yotudo/src/lib/logger"
)

var Global Settings

const USE_CMD_HIDE_WINDOW = true

func createMultipleDirectories(paths []string) error {
	for _, path := range paths {
		if dir, err := os.Open(path); err == nil {
			dir.Close()

			continue
		}

		if err := os.Mkdir(path, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}

func CreateEssentialDirectoriesAndFiles() error {
	if err := createMultipleDirectories([]string{"./data", "./data/tmp", "./data/imgs", "./data/mscs"}); err != nil {
		logger.Error(err)

		return err
	}

	if config, err := os.Open("./data/config.yaml"); err == nil {
		config.Close()
	} else if err := CreateYaml("config.yaml", ExternalConfig{}); err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func LoadSettings() (*Settings, error) {
	Global = Settings{
		App: AppSettings{
			TempLocation:   "./data/tmp",
			ImagesLocation: "./data/imgs",
			MusicsLocation: "./data/mscs",
			YTDLLocation:   "yt-dlp",
			FFMPEGLocation: "ffmpeg",
		},
		Database: DatabaseSettings{
			Location: "./data/agd_01",
			Version:  "0.1.0",
		},
	}

	if config, err := LoadYaml[ExternalConfig]("config.yaml"); err != nil {
		logger.Error(err)

		return nil, err
	} else if err := Global.Merge(config); err != nil {
		logger.Error(err)

		return nil, err
	}

	return &Global, nil
}
