package main

import (
	"embed"
	"yotudo/src"
	"yotudo/src/controller"
	"yotudo/src/database"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/service"
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
	app := src.NewApp()
	app.AddDatabaseConnection(db)

	statusRepository := repository.NewStatusRepository(db.Conn)
	genreRepository := repository.NewGenreRepository(db.Conn)
	authorRepository := repository.NewAuthorRepository(db.Conn)
	contributorRepository := repository.NewContributorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn, contributorRepository)

	fileService := service.NewFileService()
	ytService := service.NewYoutubeService(fileService)

	statusController := controller.NewStatusController(statusRepository)
	genreController := controller.NewGenreController(genreRepository)
	authorController := controller.NewAuthorController(authorRepository)
	musicController := controller.NewMusicController(musicRepository, authorRepository, contributorRepository)
	ytController := controller.NewYtController(app, musicRepository, fileService, ytService)

	if ytService.HasExecutable() {
		logger.Info("Found youtube helper executable")
	} else {
		logger.Info("Couldn't find youtube helper executable")
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "yotudo",
		Width:  1280,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []any{
			app,
			statusController,
			genreController,
			authorController,
			musicController,
			ytController,
		},
		Linux: &linux.Options{},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
