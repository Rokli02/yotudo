package service

import (
	"context"
	"fmt"
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

type YoutubeService struct {
	fileService FileService
}

func NewYoutubeService(fileService FileService) *YoutubeService {
	return &YoutubeService{fileService: fileService}
}

const (
	FILE_EXTENSION  = "webm"
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

func (s YoutubeService) DownloadVideo(ctxArg context.Context, music *model.Music) (*model.Music, error) {
	baseFilename := s.fileService.CreateFilename(music)
	filename := fmt.Sprintf("%s.%s", baseFilename, FILE_EXTENSION)
	tempFilePath := path.Join(settings.Global.App.TempLocation, filename)
	var picFilename, tempPicFilePath string

	var hasThumbnail bool
	if music.PicFilename != nil && *music.PicFilename == "thumbnail" {
		hasThumbnail = true
		picFilename = fmt.Sprintf("%s.%s", baseFilename, IMAGE_EXTENSION)
		tempPicFilePath = path.Join(settings.Global.App.TempLocation, picFilename)

		music.PicFilename = &picFilename
	}

	music.Filename = &filename
	music.Status = 1

	logger.DebugF("DownloadVideo hasThumbnail=(%t)", hasThumbnail)

	// Build command
	commandArgs := []string{
		"--ffmpeg-location", settings.Global.App.FFMPEGLocation,
		"--no-playlist",
		"-R", "3",
		"--windows-filenames",
		"-f", "bestaudio",
		"--audio-quality", "0",
	}

	if hasThumbnail {
		commandArgs = append(commandArgs, "--write-thumbnail")
	}

	logger.Debug("Temp File Path:", tempFilePath)

	commandArgs = append(commandArgs, "-o", path.Join(settings.Global.App.TempLocation, filename))

	if hasThumbnail {
		commandArgs = append(commandArgs, "-o", "thumbnail:"+path.Join(settings.Global.App.TempLocation, baseFilename))
	}

	commandArgs = append(commandArgs, music.Url)
	logger.Debug("Command was built, now executing it")

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
	if err := cmd.Start(); err != nil {
		logger.Error("YoutubeService.DownloadVideo [Couldn't start command]:", err)
		return music, err
	}

	// Logging the messages of YT download process
	if stdout, err := cmd.StdoutPipe(); err != nil {
		var buff []byte = make([]byte, 256)
		var read int

		for err == nil {
			read, err = stdout.Read(buff)

			if read > 0 {
				logger.Debug("YT_CMD:", string(buff[:read]))
			}
		}
	} else {
		logger.Warning("Something during video download:", err)

		return music, err
	}

	if err := cmd.Wait(); err != nil {
		logger.Error("YoutubeService.DownloadVideo [Something happend during command execution]:", err)

		return music, err
	}

	// Move file to music dir
	if err := s.fileService.MoveTo(tempFilePath, settings.Global.App.MusicsLocation); err != nil {
		logger.Warning("YoutubeService.DownloadVideo [Couldn't move music to its directory]", err)

		return music, nil
	}

	// Move thumbnail (if downloaded) to imgs dir
	if hasThumbnail {
		if err := s.fileService.MoveTo(tempPicFilePath, settings.Global.App.ImagesLocation); err != nil {
			logger.Warning("YoutubeService.DownloadVideo [Couldn't move thumbnail to its directory]", err)

			return music, nil
		}
	}

	// Mark music as downloaded
	music.Status = 2

	return music, nil
}
