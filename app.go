package main

import (
	"context"
	"fmt"
	"os"
	"yotudo/src/lib/logger"
)

const APP_ICON_PATH = "./assets/icon.svg"

// App struct
type App struct {
	ctx  context.Context
	Icon []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) loadAssets() {
	if icon, err := os.Open(APP_ICON_PATH); err == nil {
		stat, _ := icon.Stat()
		fileLength := stat.Size()

		a.Icon = make([]byte, fileLength)
		readLength, err := icon.Read(a.Icon)
		if err != nil || int64(readLength) != fileLength {
			logger.Error("Couldn't read the whole icon file:", fileLength)
		}

		icon.Close()
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.loadAssets()
}

func (a *App) shutdown(ctx context.Context) {}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
