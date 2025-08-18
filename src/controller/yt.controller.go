package controller

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"yotudo/src"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/service"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/image/webp"
)

type YtController struct {
	app             *src.App
	musicRepository *repository.Music
	fileService     service.FileService
	ytService       *service.YoutubeService
}

func NewYtController(
	app *src.App,
	musicRepository *repository.Music,
	fileService service.FileService,
	ytService *service.YoutubeService,
) *YtController {
	return &YtController{
		app:             app,
		musicRepository: musicRepository,
		fileService:     fileService,
		ytService:       ytService,
	}
}

type eventResult string

const (
	eventResultStart       eventResult = "start"
	eventResultDownloading eventResult = "downloading"
	eventResultCompleted   eventResult = "completed"
	eventResultFailed      eventResult = "failed"
)

func (c *YtController) DownloadByMusicId(musicId int64, eventName string) error {
	music, err := c.musicRepository.FindById(musicId)
	if err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithCancel(c.app.Ctx)
		defer cancel()

		// Starting download event
		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 0, eventResultStart, nil))

		// Give feedback about download event
		// If anything goes wrong send failed event result
		if url, err := c.ytService.PrepareUrl(music.Url, true); err != nil {
			logger.Warning(err)

			runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, -1, eventResultFailed, err))
			return
		} else {
			music.Url = url
		}

		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 10, eventResultDownloading, nil))
		savedMusic, err := c.ytService.DownloadVideo(ctx, music)
		if err != nil {
			logger.Warning(err)

			runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, -1, eventResultFailed, err))
			return
		}

		// shouldDownloadThumbnail := music.PicFilename != nil && *music.PicFilename != "thumbnail"
		// TODO: A FileService segítségével lementeni a 'PicFilename'-ban található képet, majd hozzáadni a rekordhoz

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

func (c *YtController) createEventData(musicId int64, progress float32, status eventResult, err error) []any {
	var _err *string
	if err != nil {
		tmp := err.Error()
		_err = &tmp
	}

	return []any{musicId, progress, status, _err}
}

func (c *YtController) MoveToDownloadDir(musicId int64) error {
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
		picturePath := path.Join(settings.Global.App.ImagesLocation, *music.PicFilename)
		picFile, err := os.Open(picturePath)
		if err != nil {
			goto leave_music_picfile
		}

		// Get Image's dimensions, so it can be scaled down, if necessary
		config, err := webp.DecodeConfig(picFile)
		if err != nil {
			logger.Warning("YtController.MoveToDownloadDir.decodePicFilename", err)

			goto leave_music_picfile
		}

		if tempPictureBase, found := strings.CutSuffix(*music.PicFilename, service.IMAGE_EXTENSION); found {
			tempPicturePath = tempPictureBase + "jpeg"
		} else {
			logger.WarningF("Couldn't find file extension (%s) in filename (%s)", service.IMAGE_EXTENSION, *music.PicFilename)

			goto leave_music_picfile
		}

		tempPicturePath = path.Join(settings.Global.App.TempLocation, tempPicturePath)
		var ratio float32
		if config.Width < config.Height {
			ratio = service.THUMBNAIL_SIZE / float32(config.Width)
		} else {
			ratio = service.THUMBNAIL_SIZE / float32(config.Height)
		}

		var width int
		var height int
		if ratio < 1 {
			width = max(int(float32(config.Width)*ratio), service.THUMBNAIL_SIZE)
			height = max(int(float32(config.Height)*ratio), service.THUMBNAIL_SIZE)
		}

		ctx, cancelCtx := context.WithTimeout(c.app.Ctx, time.Second*10)
		defer cancelCtx()

		logger.Debug("Creating temporary image file for thumbnail")
		if err := exec.CommandContext(ctx, settings.Global.App.FfmpegLocation, "-i", picturePath,
			"-vf", fmt.Sprintf("scale=%d:%d,crop=%d:%d:%d:%d", width, height, service.THUMBNAIL_SIZE, service.THUMBNAIL_SIZE, (width-service.THUMBNAIL_SIZE)/2, height-service.THUMBNAIL_SIZE/2),
			tempPicturePath,
		).Run(); err != nil {
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

	ctx, cancelCtx := context.WithTimeout(c.app.Ctx, time.Second*30)
	defer cancelCtx()

	filename := *music.Filename
	if _filename, found := strings.CutSuffix(filename, service.FILE_EXTENSION); found {
		filename = _filename + "mp3"
	}

	ffmpegArguments = append(ffmpegArguments, path.Join(settings.Global.App.DownloadLocation, filename))

	if err := exec.CommandContext(ctx, settings.Global.App.FfmpegLocation, ffmpegArguments...).Run(); err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (c *YtController) addMetadatas(arguments *[]string, music *model.Music) {
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
