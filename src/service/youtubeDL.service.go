package service

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"
)

var ytRegexp = regexp.MustCompile(`^(https?://)?(www.)?(youtube.com|youtu.be)/(watch\?v=\S{11})`)

type YoutubeDLService struct {
	fileService FileService
}

func NewYoutubeDLService(fileService FileService) *YoutubeDLService {
	return &YoutubeDLService{fileService: fileService}
}

const (
	FILE_EXTENSION        = "webm"
	FINAL_IMAGE_EXTENSION = "jpeg"
	FINAL_MUSIC_EXTENSION = "mp3"
	YT_THUMBNAIL_URL      = "https://i.ytimg.com/vi/%s/0.jpg"
)

func (s YoutubeDLService) PrepareUrl(url string, stripUnnecessaryParameters bool) (string, error) {
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

func (s YoutubeDLService) HasExecutable() bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelCtx()

	res := make(chan bool)
	go func() {
		cmd := exec.CommandContext(ctx, settings.Global.App.YTDLLocation, "--version")
		if settings.USE_CMD_HIDE_WINDOW {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow:    true,
				CreationFlags: 0x08000000,
			}
		}

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

func (s YoutubeDLService) DownloadVideo(ctxArg context.Context, music *model.Music) (*model.Music, error) {
	baseFilename := s.fileService.CreateFilename(music)
	filename := fmt.Sprintf("%s.%s", baseFilename, FILE_EXTENSION)
	tempFilePath := path.Join(settings.Global.App.TempLocation, filename)

	music.Filename = &filename
	music.Status = 1

	// Build command
	commandArgs := []string{
		"--no-playlist",
		"-R", "3",
		"--windows-filenames",
		"-f", "bestaudio",
		"-o", path.Join(settings.Global.App.TempLocation, filename),
		music.Url,
	}

	ctx, cancelCtx := context.WithTimeout(ctxArg, time.Second*60)
	defer cancelCtx()

	// Download to Temp
	cmd := exec.CommandContext(ctx, settings.Global.App.YTDLLocation, commandArgs...)
	if settings.USE_CMD_HIDE_WINDOW {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error("YoutubeService.DownloadVideo [Something happend during command execution]:", err.Error()+": ", strings.ReplaceAll(stderr.String(), "\n", " | "))
		return music, err
	}

	// Move file to music dir
	if err := s.fileService.MoveTo(tempFilePath, settings.Global.App.MusicsLocation); err != nil {
		logger.Warning("YoutubeService.DownloadVideo [Couldn't move music to its directory]", err)

		return music, err
	}

	// Mark music as downloaded
	music.Status = 2

	return music, nil
}

func (s YoutubeDLService) GetVideoThumbnailUrl(rawUrl string) (string, error) {
	videoUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	vId := videoUrl.Query().Get("v")
	return fmt.Sprintf(YT_THUMBNAIL_URL, vId), nil
}
