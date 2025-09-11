package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"yotudo/src/lib/logger"
	"yotudo/src/settings"
)

type AssetsHandler struct{}

var _ http.Handler = (*AssetsHandler)(nil)

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

func (*AssetsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if imageName, found := strings.CutPrefix(req.URL.Path, "/image/"); found {
		imagePath := path.Join(settings.Global.App.ImagesLocation, imageName)

		file, err := os.ReadFile(imagePath)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Requesting image with name: \"%s\", but found nothing", imageName)
		}

		res.Write(file)

		return
	}

	res.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(res, "Requesting asset with name: \"%s\", but found nothing", req.URL.Path)
	logger.WarningF("Requesting asset with name: \"%s\", but found nothing", req.URL.Path)
}
