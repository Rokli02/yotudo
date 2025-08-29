package settings

import (
	"os"
	"yotudo/src/lib/yaml"
)

var Global Settings

const USE_CMD_HIDE_WINDOW = true

func createMultipleDirectories(paths []string) error {
	for _, path := range paths {
		if dir, err := os.Open(path); err == nil {
			dir.Close()

			continue
		}

		if err := os.MkdirAll(path, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}

func CreateEssentialDirectoriesAndFiles() error {
	if err := createMultipleDirectories([]string{"./data/tmp", "./data/imgs", "./data/mscs", "./data/logs"}); err != nil {
		return err
	}

	if config, err := os.Open("./data/config.yaml"); err == nil {
		config.Close()
	} else if err := yaml.CreateFile("config.yaml", settingsYaml{
		Logger: settingsYaml_Logger{
			Level: "info",
			Types: []string{"console"},
		},
		App: settingsYaml_App{
			FFMPEGLocation: "ffmpeg",
			YTDLLocation:   "yt-dlp",
		},
		Database: settingsYaml_Database{
			Location: "./data/yU0dRywKd",
		},
	}); err != nil {
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
			Version:  "1.0.1",
		},
		Logger: LoggerSettings{
			Level: "warning",
			Types: []string{"console"},
		},
	}

	if config, err := yaml.LoadFile[settingsYaml]("config.yaml"); err != nil {
		return nil, err
	} else if err := Global.Merge(config); err != nil {
		return nil, err
	}

	return &Global, nil
}
