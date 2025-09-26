package handler

import (
	"fmt"
	"net/http"
)

func ApiHandlers() http.Handler {
	mux := http.NewServeMux()

	var musicsHandler *MusicsHandler = nil
	var statusHandler *StatusHandler = nil

	mux.HandleFunc("/api/musics", musicsHandler.GetManyMusics)
	mux.HandleFunc("/api/musics/{id}", musicsHandler.GetMusicById)
	mux.HandleFunc("/api/musics/{id}/download", musicsHandler.DownloadMusicById)

	mux.HandleFunc("/api/status", statusHandler.GetManyStatus)

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Ismeretlen elérési útvonal\n")
	})
	return mux
}
