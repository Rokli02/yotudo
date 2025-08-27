package settings

import "fmt"

type Settings struct {
	App      AppSettings
	Database DatabaseSettings
}

func (s *Settings) Merge(other *settingsYaml) error {
	errStr := "Didn't merge property '%s', because it was empty"

	if other.App.DownloadLocation != "" {
		s.App.DownloadLocation = other.App.DownloadLocation
	} else {
		return fmt.Errorf(errStr, "App.DownloadLocation")
	}

	if other.App.YTDLLocation != "" {
		s.App.YTDLLocation = other.App.YTDLLocation
	}

	if other.App.FFMPEGLocation != "" {
		s.App.FFMPEGLocation = other.App.FFMPEGLocation
	}

	if other.Database.Location != "" {
		s.Database.Location = other.Database.Location
	}

	return nil
}

type DatabaseSettings struct {
	Location string
	Version  string
}

type AppSettings struct {
	TempLocation     string
	MusicsLocation   string
	ImagesLocation   string
	DownloadLocation string
	YTDLLocation     string
	FFMPEGLocation   string
}

type settingsYaml struct {
	App      settingsYaml_App      `yaml:"app"`
	Database settingsYaml_Database `yaml:"database"`
	Logger   settingsYaml_Logger   `yaml:"logger"`
}

type settingsYaml_App struct {
	DownloadLocation string `yaml:"downloadLocation"`
	YTDLLocation     string `yaml:"ytdlLocation"`
	FFMPEGLocation   string `yaml:"ffmpegLocation"`
}

type settingsYaml_Database struct {
	Location string `yaml:"location"`
}

type settingsYaml_Logger struct {
	// Enum that may be one of the following: [error, warning, info, debug]
	Level string `yaml:"level"`
	// Logger types that may be used during the app runtime and must be at least one or more of the followings: [console, file, ???database???]
	Types []string `yaml:"types"`
}

//TODO: Ellenőrizni, HA nem töltődne be valami oknál fogva a configuráció a 'yaml' fájlból, lehetséges hogy a privát típus miatt történik
