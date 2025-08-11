package main

import (
	"embed"
	"yotudo/src/controller"
	"yotudo/src/database"
	"yotudo/src/database/repository"
	"yotudo/src/settings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	settings.LoadSettings()
	db := database.LoadDatabase()
	// Create an instance of the app structure
	app := NewApp(db)

	statusRepository := repository.NewStatusRepository(db.Conn)
	genreRepository := repository.NewGenreRepository(db.Conn)
	authorRepository := repository.NewAuthorRepository(db.Conn)
	contributorRepository := repository.NewContributorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "yotudo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []any{
			app,
			controller.NewStatusController(statusRepository),
			controller.NewGenreController(genreRepository),
			controller.NewAuthorController(authorRepository),
			controller.NewMusicController(musicRepository, authorRepository, contributorRepository),
		},
		Linux: &linux.Options{},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
