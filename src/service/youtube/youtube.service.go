package youtube

import (
	"fmt"
	"regexp"
	"strings"
)

var ytRegexp = regexp.MustCompile(`^(https?://)?(www.)?(youtube.com|youtu.be)/(watch\?v=\S{11})`)

type YoutubeService uint8

func NewYoutubeService() YoutubeService {
	return YoutubeService(0)
}

func (s YoutubeService) PrepareUrl(url string, stripUnnecessaryParameters bool) (string, error) {
	if !ytRegexp.Match([]byte(url)) {
		return "", fmt.Errorf("the given url is not a youtube video link")
	}

	if stripUnnecessaryParameters {
		indexOfQueryStringStart := strings.Index(url, "?")
		sb := strings.Builder{}
		sb.WriteString(url[:indexOfQueryStringStart+1])

		queryStrings := strings.Split(url[indexOfQueryStringStart+1:], "&")
		for _, queryString := range queryStrings {
			if strings.HasPrefix(queryString, "v=") {
				sb.WriteString(queryString)

				break
			}
		}

		builtUrl := sb.String()

		return builtUrl, nil
	}

	return url, nil
}

func (s YoutubeService) DownloadVideo(url string, saveThumbnail bool) {
	// yt-dl download
	// --extract-audio
	// --audio-format best
	// --audio-quality 0
	// --ffmpeg-location settings.Global.App.FfmpegLocation

	if saveThumbnail {
		// save thumbnail image IF got downloaded
		// --write-thumbnail
	}
}
