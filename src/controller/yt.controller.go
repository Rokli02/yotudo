package controller

import (
	"context"
	"yotudo/src"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
