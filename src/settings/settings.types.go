package settings

type DatabaseSettings struct {
	Location string
	Version  string
}

type AppSettings struct {
	TempLocation     string
	MusicsLocation   string
	ImagesLocation   string
	DownloadLocation string
	YtdlLocation     string
	FfmpegLocation   string
}

type Settings struct {
	App      AppSettings
	Database DatabaseSettings
}
