package settings

import "fmt"

type DatabaseSettings struct {
	Location string `yaml:"location"`
	Version  string
}

type AppSettings struct {
	TempLocation     string
	MusicsLocation   string
	ImagesLocation   string
	DownloadLocation string `yaml:"downloadLocation"`
	YTDLLocation     string `yaml:"ytdlLocation"`
	FFMPEGLocation   string `yaml:"ffmpegLocation"`
}

type Settings struct {
	App      AppSettings      `yaml:"app"`
	Database DatabaseSettings `yaml:"database"`
}

func (s *Settings) Merge(other *Settings, errOnNotFound bool) error {
	errStr := "Didn't merge property '%s', because it was empty"

	if other.App.DownloadLocation != "" {
		s.App.DownloadLocation = other.App.DownloadLocation
	} else if errOnNotFound {
		return fmt.Errorf(errStr, "App.DownloadLocation")
	}

	if other.App.YTDLLocation != "" {
		s.App.YTDLLocation = other.App.YTDLLocation
	} else if errOnNotFound {
		return fmt.Errorf(errStr, "App.YTDLLocation")
	}

	if other.App.FFMPEGLocation != "" {
		s.App.FFMPEGLocation = other.App.FFMPEGLocation
	} else if errOnNotFound {
		return fmt.Errorf(errStr, "App.FFMPEGLocation")
	}

	if other.Database.Location != "" {
		s.Database.Location = other.Database.Location
	}

	return nil
}
