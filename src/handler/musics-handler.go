package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/service"
	"yotudo/src/settings"
)

type MusicsHandler struct{}

func (m *MusicsHandler) GetManyMusics(w http.ResponseWriter, r *http.Request) {
	page := &model.Page{Page: 0, Size: 10}
	var filter string
	var statusId int = 2
	var sort []model.Sort = []model.Sort{{Key: "updated_at", Dir: -1}}

	for key, values := range r.URL.Query() {
		if len(values) == 0 || values[0] == "" {
			continue
		}

		value := values[0]

		switch key {
		case "p":
			p, err := strconv.Atoi(value)
			if err != nil {
				logger.Warning("Server Error:", err)

				continue
			}
			if p >= 0 {
				page.Page = p
			}
		case "ps":
			ps, err := strconv.Atoi(value)
			if err != nil {
				logger.Warning("Server Error:", err)

				continue
			}
			if ps > 0 {
				page.Size = ps
			}
		case "f":
			filter = value
		}
	}

	result := service.GlobalMusicService.GetManyByPagination(filter, statusId, page, sort)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error("Server Error:", err)
	}
}

func (m *MusicsHandler) GetMusicById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		logger.Error("Server Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Érvénytelen ID került megadásra")

		return
	}

	music, err := service.GlobalMusicService.GetById(id)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Nem talált")
		} else {
			logger.Error("Server Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Érvénytelen ID került megadásra")
		}

		return
	}

	if err = json.NewEncoder(w).Encode(music); err != nil {
		logger.Error("Server Error:", err)
	}
}

func (m *MusicsHandler) DownloadMusicById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		logger.Error("Server Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Érvénytelen ID került megadásra")

		return
	}

	music, err := service.GlobalMusicService.GetById(id)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Nem talált")
		} else {
			logger.Error("Server Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Érvénytelen ID került megadásra")
		}

		return
	}

	if music.Filename == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Még nem feldolgozott zenét próbáltál letölteni")

		return
	}

	var musicFile *os.File

	if _musicFile, err := os.Open(path.Join(settings.Global.App.MusicsLocation, *music.Filename)); err != nil {
		logger.WarningF("Server Warning: Music was not found in its directory (filename=%s)", *music.Filename)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Nem találta meg a zenét a szerver")

		return
	} else {
		musicFile = _musicFile
	}

	defer musicFile.Close()

	filename := service.GlobalFileService.CreateFilename(music)
	w.Header().Set("x-file-name", filename)

	//TODO: Felhasználni a YoutubeService.MoveToDownloadDir függvényt, hogy áthelyezzük a tmp mappába, majd onnan a borítóval és metaadatokkal feltöltött zenét töltesse le
	//TODO: Csak teszt jelleggel kommentáltam ki, utána ki kell szedni
	// if _, err := io.CopyBuffer(w, musicFile, nil); err != nil {
	// 	logger.Error("Server Error:", err)
	// }
}
