package controller

import (
	"context"
	"time"
	"yotudo/src"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type YtController struct {
	app             *src.App
	ytService       *service.YoutubeService
	musicRepository *repository.Music
}

func NewYtController(app *src.App, ytService *service.YoutubeService, musicRepository *repository.Music) *YtController {
	return &YtController{
		app:             app,
		ytService:       ytService,
		musicRepository: musicRepository,
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

		if err := c.ytService.DownloadVideo(music); err != nil {
			logger.Warning(err)

			runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, -1, eventResultFailed, err))
			return
		}

		/* DELETE start*/
		for i := range 3 {
			runtime.EventsEmit(ctx, eventName, musicId, i)
			logger.Debug("download-progress event")

			time.Sleep(time.Second * 1)
		}
		/* DELETE end*/

		// Mark event as finished
		runtime.EventsEmit(ctx, eventName, c.createEventData(musicId, 100, eventResultCompleted, nil))
	}()

	return nil
}

func (c *YtController) createEventData(musicId int64, progress float32, status eventResult, err error) []any {
	return []any{musicId, progress, status, err}
}
