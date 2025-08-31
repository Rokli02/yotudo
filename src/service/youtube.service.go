package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type YoutubeService struct {
	ctx             *context.Context
	musicRepository *repository.Music
	fileService     FileService
	ytDLService     *YoutubeDLService
}

func NewYoutubeService(
	ctx *context.Context,
	musicRepository *repository.Music,
	fileService FileService,
	ytDLService *YoutubeDLService,
) *YoutubeService {
	return &YoutubeService{
		ctx:             ctx,
		musicRepository: musicRepository,
		fileService:     fileService,
		ytDLService:     ytDLService,
	}
}

type eventResult string

const (
	eventResultStart       eventResult = "start"
	eventResultDownloading eventResult = "downloading"
	eventResultCompleted   eventResult = "completed"
	eventResultFailed      eventResult = "failed"
)

func (c *YoutubeService) DownloadByMusicId(musicId int64, eventName string) error {
	music, err := c.musicRepository.FindById(musicId)
	if err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithCancel(*c.ctx)
		defer cancel()

		// Starting download event
		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 0, eventResultStart, nil))

		// Give feedback about download event
		// If anything goes wrong send failed event result
		if url, err := c.ytDLService.PrepareUrl(music.Url, true); err != nil {
			logger.Warning(err)

			runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, -1, eventResultFailed, err))
			return
		} else {
			music.Url = url
		}

		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 10, eventResultDownloading, nil))
		savedMusic, err := c.ytDLService.DownloadVideo(ctx, music)
		if err != nil {
			logger.Warning(err)

			runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, -1, eventResultFailed, err))
			return
		}

		_, err = c.musicRepository.UpdateOne(savedMusic.Id, savedMusic.ToUpdateMusic())
		if err != nil {
			logger.Error("CRITICAL ERROR: video got downloaded but couldn't update its record in db:", err)

			return
		}

		// Mark event as finished
		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 100, eventResultCompleted, nil))
	}()

	return nil
}

func (c *YoutubeService) createEventData(musicId int64, progress float32, status eventResult, err error) []any {
	var _err *string
	if err != nil {
		tmp := err.Error()
		_err = &tmp
	}

	return []any{musicId, progress, status, _err}
}

func (c *YoutubeService) MoveToDownloadDir(musicId int64) error {
	music, err := c.musicRepository.FindById(musicId)
	if err != nil {
		return err
	}

	// Check if music wasn't preprocessed, as it should have
	if music.Filename == nil {
		err := fmt.Sprintf("Music didn't get downloaded previously (id=%d)", musicId)

		logger.Warning(err)

		return errors.New(err)
	}

	// Check for music file
	musicPath := path.Join(settings.Global.App.MusicsLocation, *music.Filename)
	if !c.fileService.IsExists(musicPath) {
		err := fmt.Sprintf("Music was not found in its directory (filename=%s)", *music.Filename)

		logger.Warning(err)

		return errors.New(err)
	}

	// Check for music with same name in temp directory
	if c.fileService.IsExists(path.Join(settings.Global.App.TempLocation, *music.Filename)) {
		err := fmt.Sprintf("For some unknown reason music wasn't moved from temp directory (filename=%s)", *music.Filename)

		logger.Warning(err)

		return errors.New(err)
	}

	var ffmpegArguments []string = []string{
		"-i", musicPath,
	}

	var tempPicturePath string
	defer func() {
		if tempPicturePath != "" {
			if err := os.Remove(tempPicturePath); err != nil {
				logger.Warning("Couldn't delete TEMP picture:", err)
			}
		}
	}()

	// Check for a thumbnail that might be attached to the music
	if music.PicFilename != nil {
		imageWidth, imageHeight, imageExt, err := c.fileService.GetImageConfig(*music.PicFilename)
		if err != nil {
			logger.Error(err)

			goto leave_music_picfile
		}

		imageSize := min(THUMBNAIL_SIZE, imageWidth, imageHeight)
		offsetX, offsetY := (imageWidth-imageSize)/2, (imageHeight-imageSize)/2

		if tempPictureBase, found := strings.CutSuffix(*music.PicFilename, imageExt); found {
			tempPicturePath = fmt.Sprintf("%s.%s", tempPictureBase, FINAL_IMAGE_EXTENSION)
		} else {
			logger.WarningF("Couldn't find file extension (%s) in filename (%s)", imageExt, *music.PicFilename)

			goto leave_music_picfile
		}

		picturePath := path.Join(settings.Global.App.ImagesLocation, *music.PicFilename)
		tempPicturePath = path.Join(settings.Global.App.TempLocation, tempPicturePath)

		ctx, cancelCtx := context.WithTimeout(*c.ctx, time.Second*10)
		defer cancelCtx()

		cmd := exec.CommandContext(ctx, settings.Global.App.FFMPEGLocation,
			"-i", picturePath,
			"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:%d:%d", imageWidth, imageHeight, imageSize, imageSize, offsetX, offsetY),
			tempPicturePath,
		)
		if settings.USE_CMD_HIDE_WINDOW {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow:    true,
				CreationFlags: 0x08000000,
			}
		}
		if err := cmd.Run(); err != nil {
			logger.Error(err)

			goto leave_music_picfile
		}

		ffmpegArguments = append(ffmpegArguments,
			"-i", tempPicturePath,
			"-map", "0",
			"-map", "1",
			"-disposition:v:1", "attached_pic",
		)
	}
leave_music_picfile:

	c.addMetadatas(&ffmpegArguments, music)

	ctx, cancelCtx := context.WithTimeout(*c.ctx, time.Second*30)
	defer cancelCtx()

	filename := *music.Filename
	if _filename, found := strings.CutSuffix(filename, FILE_EXTENSION); found {
		filename = _filename + FINAL_MUSIC_EXTENSION
	}

	ffmpegArguments = append(ffmpegArguments, path.Join(settings.Global.App.DownloadLocation, filename))

	cmd := exec.CommandContext(ctx, settings.Global.App.FFMPEGLocation, ffmpegArguments...)
	if settings.USE_CMD_HIDE_WINDOW {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}

	if err := cmd.Run(); err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (c *YoutubeService) addMetadatas(arguments *[]string, music *model.Music) {
	*arguments = append(*arguments,
		"-y",
		"-id3v2_version", "3",
		"-metadata", "title="+music.Name,
		"-metadata", "genre="+music.Genre.Name,
		"-metadata", "album_artist="+music.Author.Name,
	)

	if music.Album != nil {
		*arguments = append(*arguments, "-metadata", "album="+*music.Album)
	}

	if music.Published != nil {
		*arguments = append(*arguments, "-metadata", fmt.Sprintf("date=%d", *music.Published))
	}

	contributorMap := map[int64]model.Author{
		music.Author.Id: music.Author,
	}

	if len(music.Contributors) != 0 {
		for _, contributor := range music.Contributors {
			contributorMap[contributor.Id] = contributor
		}

	}
	contributors := make([]string, 0, len(contributorMap))
	for _, value := range contributorMap {
		contributors = append(contributors, value.Name)
	}

	*arguments = append(*arguments, "-metadata", "artist="+strings.Join(contributors, ";"))
}
