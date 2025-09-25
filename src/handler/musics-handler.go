package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

	filename, err := service.GlobalYoutubeService.MoveToDir(id, settings.Global.App.TempLocation)
	if err != nil {
		logger.Error("Server Error:", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	musicPath := path.Join(settings.Global.App.TempLocation, filename)
	musicFile, err := os.Open(musicPath)
	if err != nil {
		logger.Error("Server Error:", err)
	}

	defer func() {
		err := os.Remove(musicPath)
		if err != nil {
			logger.Warning("Server Warning:", err)
		}
	}()
	defer musicFile.Close()

	dst := base64.StdEncoding.EncodeToString([]byte(filename))
	w.Header().Set("x-file-name", dst)
	w.WriteHeader(http.StatusOK)

	if _, err := io.CopyBuffer(w, musicFile, nil); err != nil {
		logger.Error("Server Error:", err)
	}
}
