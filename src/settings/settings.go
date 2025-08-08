package settings

var Global Settings

func LoadSettings() Settings {
	Global = Settings{
		App: AppSettings{
			TempLocation:   "./data/tmp",
			ImagesLocation: "./data/imgs",
			MusicsLocation: "./data/mscs",
			YtdlLocation:   "/usr/local/bin/youtube-dl",
			FfmpegLocation: "/bin/ffmpeg:",
		},
		Database: DatabaseSettings{
			Location: "./data/agd_01",
			Version:  "0.1.0",
		},
	}

	return Global
}
