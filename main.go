package main

import (
	"embed"
	"yotudo/src"
	"yotudo/src/database"
	"yotudo/src/handler"
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
	if err := settings.CreateEssentialDirectoriesAndFiles(); err != nil {
		panic(err)
	}

	if _, err := settings.LoadSettings(); err != nil {
		panic(err)
	}

	_, closeLoggers := logger.InitializeLogger(settings.Global.Logger.Level, settings.Global.Logger.Types)
	defer closeLoggers()

	db := database.LoadDatabase()
	defer db.Close()

	app := src.NewApp()

	statusService := &service.StatusService{}
	genreService := &service.GenreService{}
	authorService := &service.AuthorService{}
	musicService := &service.MusicService{}
	youtubeService := &service.YoutubeService{}
	dialogService := service.NewDialogService(&app.Ctx)

	if !service.GlobalFileService.HasExecutable() {
		logger.Error("Couldn't find ffmpeg executable")

		panic("Couldn't find ffmpeg executable")
	}

	if !service.GlobalYoutubeDLService.HasExecutable() {
		logger.Error("Couldn't find youtube helper executable")

		panic("Couldn't find youtube helper executable")
	}

	windowWidth, windowHeight := service.GlobalInfoService.GetWindowSize()

	assetsHandler := handler.NewAssetsHandler()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "yotudo",
		MinWidth:  400,
		MinHeight: 280,
		Width:     windowWidth,
		Height:    windowHeight,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: assetsHandler,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		OnBeforeClose:    app.BeforeClose,
		Bind: []any{
			app,
			statusService,
			genreService,
			authorService,
			musicService,
			youtubeService,
			dialogService,
		},
		Linux: &linux.Options{},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
