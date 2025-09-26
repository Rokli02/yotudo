package src

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"yotudo/src/lib/logger"
	"yotudo/src/service"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx       context.Context
	serverCmd *exec.Cmd
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

	a.StopServer()

	return
}

func (a *App) StartServer() error {
	a.serverCmd = exec.CommandContext(context.Background(), "./yotudo-server.exe",
		"--localhost",
		"--host-frontend",
		"--port", fmt.Sprintf("%d", settings.Global.Server.Port),
	)

	go func() {
		stdout, _ := a.serverCmd.StdoutPipe()
		buf := make([]byte, 4096)

		for {
			if a.serverCmd == nil {
				return
			}

			read, err := stdout.Read(buf)
			if err != nil {
				logger.Warning(err)
				return
			}

			if read != 0 {
				var prefix string
				out, _ := strings.CutSuffix(string(buf), "\n")

				if indexOfColon := strings.Index(out, ":"); indexOfColon != -1 {
					prefix = out[:indexOfColon]
					out = strings.TrimSpace(out[indexOfColon+1:])
				}

				switch prefix {
				case "INFO":
					logger.Info(out)
				case "ERR":
					logger.Error(out)
				case "WARN":
					logger.Warning(out)
				case "EVENT":
					// TODO: Process events
				default:
					logger.Debug(out)
				}
			}
		}
	}()

	return a.serverCmd.Start()
}

func (a *App) StopServer() error {
	if a.serverCmd == nil {
		return nil
	}

	cmd := a.serverCmd
	a.serverCmd = nil

	cmd.Stdin.Read([]byte("exit"))
	return cmd.Wait()
}
