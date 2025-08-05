package youtube

import (
	"testing"
	"yotudo/src/lib/logger"
	"yotudo/src/service/youtube"
)

func TestPrepareUrl(t *testing.T) {
	const url = "https://www.youtube.com/watch?v=wRIkfMSnED4&list=PLIpNwAgyqIjkdKHEpCWq4z5FTgLKVcLwY&index=4"

	service := youtube.NewYoutubeService()

	result, err := service.PrepareUrl(url, true)

	if err != nil {
		t.Error(err)
	}

	logger.Info("PrepareUrl result:", result)
}
