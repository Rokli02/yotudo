package handler

import (
	"encoding/json"
	"net/http"
	"yotudo/src/lib/logger"
	"yotudo/src/service"
)

type StatusHandler struct{}

func (s *StatusHandler) GetManyStatus(w http.ResponseWriter, r *http.Request) {
	result := service.GlobalStatusService.GetAll()

	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error("Server Error:", err)
	}
}
