package settings

import "fmt"

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

type Settings struct {
	App      AppSettings
	Database DatabaseSettings
}

type ExternalConfig struct {
	App struct {
		DownloadLocation string `yaml:"downloadLocation"`
		YTDLLocation     string `yaml:"ytdlLocation"`
		FFMPEGLocation   string `yaml:"ffmpegLocation"`
	} `yaml:"app"`
	Database struct {
		Location string `yaml:"location"`
	} `yaml:"database"`
}

func (s *Settings) Merge(other *ExternalConfig) error {
	errStr := "Didn't merge property '%s', because it was empty"

	if other.App.DownloadLocation != "" {
		s.App.DownloadLocation = other.App.DownloadLocation
	} else {
		return fmt.Errorf(errStr, "App.FFMPEGLocation")
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
