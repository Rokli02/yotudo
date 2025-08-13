package service

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"
)

var ytRegexp = regexp.MustCompile(`^(https?://)?(www.)?(youtube.com|youtu.be)/(watch\?v=\S{11})`)

type YoutubeService struct{}

func NewYoutubeService() *YoutubeService {
	return &YoutubeService{}
}

const (
	FILE_EXTENSION  = "mp3"
	IMAGE_EXTENSION = "webp"
)

func (s YoutubeService) PrepareUrl(url string, stripUnnecessaryParameters bool) (string, error) {
	if !ytRegexp.Match([]byte(url)) {
		return "", fmt.Errorf("the given url is not a youtube video link")
	}

	if stripUnnecessaryParameters {
		indexOfQueryStringStart := strings.Index(url, "?")
		sb := strings.Builder{}
		sb.WriteString(url[:indexOfQueryStringStart+1])

		queryStrings := strings.Split(url[indexOfQueryStringStart+1:], "&")
		for _, queryString := range queryStrings {
			if strings.HasPrefix(queryString, "v=") {
				sb.WriteString(queryString)

				break
			}
		}

		builtUrl := sb.String()

		return builtUrl, nil
	}

	return url, nil
}

func (s YoutubeService) HasExecutable() bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelCtx()

	res := make(chan bool)
	go func() {
		cmd := exec.CommandContext(ctx, settings.Global.App.YtdlLocation, "--version")

		if stdout, err := cmd.Output(); err != nil {
			logger.Error(err)

			res <- false
		} else if len(stdout) == 0 {
			logger.Warning("Didn't write any result to stdout")

			res <- false
		}

		res <- true
	}()

	select {
	case <-ctx.Done():
		logger.Error("Context timed out after 5 seconds")

		return false
	case v := <-res:
		return v
	}
}

func (s YoutubeService) DownloadVideo(music *model.Music) error {
	// HA "music.PicFilename" === "thumbnail" ---> Használja a youtube-os képet, különben a megadottat töltse le, ha még nincs
	var useThumbnail bool
	var hasPic bool
	if music.PicFilename != nil {
		hasPic = true

		if *music.PicFilename == "thumbnail" {
			useThumbnail = true
		}
	}

	logger.DebugF("DownloadVideo hasPic=%t, useThumbnail=%t", hasPic, useThumbnail)

	// Download to Temp

	// Link filename to DB record
	// Link thumbnail pic filename to DB record

	// Move file to music dir
	// Move pic to imgs dir

	// yt-dl download
	// --extract-audio
	// --audio-format best
	// --audio-quality 0
	// --ffmpeg-location settings.Global.App.FfmpegLocation

	// if music.saveThumbnail {
	// save thumbnail image IF got downloaded
	// --write-thumbnail
	// }

	// Ha mindent sikeresen letöltött, a music-on belül állítsa át a "FileName" és "PicFileName" mezőket

	return nil
}
