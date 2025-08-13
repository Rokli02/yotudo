package src

import (
	"context"
	"yotudo/src/database"
)

// App struct
type App struct {
	Ctx context.Context
	db  *database.Database
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) AddDatabaseConnection(db *database.Database) {
	a.db = db
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
