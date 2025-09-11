package src

import (
	"context"
	"os"
	"path/filepath"
	"yotudo/src/lib/logger"
	"yotudo/src/service"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx context.Context
}

func NewApp() *App {
	app := &App{}

	return app
}

func (a *App) Startup(ctx context.Context) {
	logger.Info("Application is starting up...")
	a.Ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	logger.Info("Application is shuting down...")

	// Delete every file from /tmp folder
	if tempDir, err := os.Open(settings.Global.App.TempLocation); err != nil {
		logger.Warning(err)
	} else {
		defer tempDir.Close()

		filenames, err := tempDir.Readdirnames(0)
		if err != nil {
			logger.Warning(err)
			goto skip_tmp_dir_prune
		}

		for _, file := range filenames {
			err = os.RemoveAll(filepath.Join(settings.Global.App.TempLocation, file))
			if err != nil {
				logger.Warning(err)
			}
		}
	}
skip_tmp_dir_prune:
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	if service.GlobalInfoService != nil {
		if err := service.GlobalInfoService.SetWindowSize(runtime.WindowGetSize(ctx)); err != nil {
			logger.Error(err)
		}
	}

	return
}
