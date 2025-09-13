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
	Ctx        context.Context
	httpServer *Server
}

func NewApp(httpServer *Server) *App {
	app := &App{
		httpServer: httpServer,
	}

	return app
}

func (a *App) Startup(ctx context.Context) {
	logger.Info("Application is starting up...")
	a.Ctx = ctx
	a.StartServer()
}

func (a *App) Shutdown(ctx context.Context) {
	logger.Info("Application is shuting down...")

	a.StopServer()

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

func (a *App) StartServer() error {
	return a.httpServer.Start()
}

func (a *App) StopServer() error {
	return a.httpServer.Stop(a.Ctx)
}
