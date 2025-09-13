package settings

import (
	"encoding/json"
	"fmt"
)

type Settings struct {
	App      AppSettings
	Database DatabaseSettings
	Logger   LoggerSettings
	Server   ServerSettings
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

	if other.Logger.Level != "" {
		s.Logger.Level = other.Logger.Level
	}

	if other.Server.Port != 0 {
		s.Server.Port = other.Server.Port
	}

	if len(other.Logger.Types) != 0 {
		s.Logger.Types = other.Logger.Types
	}

	return nil
}

func (s *Settings) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
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

type LoggerSettings struct {
	Level string
	Types []string
}

type ServerSettings struct {
	Port int
}

type settingsYaml struct {
	App    settingsYaml_App    `yaml:"app"`
	Logger settingsYaml_Logger `yaml:"logger"`
	Server settingsYaml_Server `yaml:"server"`
}

type settingsYaml_App struct {
	DownloadLocation string `yaml:"downloadLocation"`
	YTDLLocation     string `yaml:"ytdlLocation"`
	FFMPEGLocation   string `yaml:"ffmpegLocation"`
}

type settingsYaml_Logger struct {
	// Enum that may be one of the following: [error, warning, info, debug]
	Level string `yaml:"level"`
	// Logger types that may be used during the app runtime and must be at least one or more of the followings: [console, file, ???database???]
	Types []string `yaml:"types"`
}

type settingsYaml_Server struct {
	Port int `yaml:"port"`
}
