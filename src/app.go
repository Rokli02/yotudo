package src

import (
	"context"
	"yotudo/src/database"
	"yotudo/src/database/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	Ctx            context.Context
	db             *database.Database
	infoRepository *repository.Info
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) SetDatabaseConnection(db *database.Database) {
	a.db = db
}

func (a *App) SetInfoRepository(infoRepository *repository.Info) {
	a.infoRepository = infoRepository
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	if a.infoRepository != nil {
		runtime.WindowGetSize(a.Ctx)
	}

	if a.db != nil {
		a.db.Close()
	}
}
