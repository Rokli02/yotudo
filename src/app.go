package src

import (
	"context"
	"yotudo/src/database"
	"yotudo/src/lib/logger"
	"yotudo/src/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	Ctx         context.Context
	db          *database.Database
	infoService *service.InfoService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) SetDatabaseConnection(db *database.Database) {
	a.db = db
}

func (a *App) SetInfoService(infoService *service.InfoService) {
	a.infoService = infoService
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	if a.infoService != nil {
		if err := a.infoService.SetWindowSize(runtime.WindowGetSize(ctx)); err != nil {
			logger.Error(err)
		}
	}

	return
}
