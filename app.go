package main

import (
	"context"
	"yotudo/src/database"
)

// App struct
type App struct {
	ctx context.Context
	db  *database.Database
}

// NewApp creates a new App application struct
func NewApp(db *database.Database) *App {
	return &App{
		db: db,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}
