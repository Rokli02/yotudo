package src

import (
	"context"
	"os"
	"path/filepath"
	"yotudo/src/database"
	"yotudo/src/lib/logger"
	"yotudo/src/service"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx         context.Context
	db          *database.Database
	infoService *service.InfoService
}

func NewApp() *App {
	return &App{}
}

func (a *App) SetDatabaseConnection(db *database.Database) {
	a.db = db
}

func (a *App) SetInfoService(infoService *service.InfoService) {
	a.infoService = infoService
}

func (a *App) Startup(ctx context.Context) {
	logger.Info("Application is starting up...")
	a.Ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	logger.Info("Application is shuting down...")

	if a.db != nil {
		a.db.Close()
	}

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
	if a.infoService != nil {
		if err := a.infoService.SetWindowSize(runtime.WindowGetSize(ctx)); err != nil {
			logger.Error(err)
		}
	}

	return
}
