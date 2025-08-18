package service

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	pathModule "path"
	"regexp"
	"strings"
	"time"
	"yotudo/src/model"
	"yotudo/src/settings"

	"golang.org/x/image/webp"
)

type FileService uint8

const THUMBNAIL_SIZE = 512

var (
	validFilenameRegexp       = regexp.MustCompile(`^([^\\\:\*\?\"\/\<\>\|\0]{1,255})$`)
	invalidFilenameCharacters = regexp.MustCompile(`[\\\:\*\?\"\/\<\>\|\0]`)
)

func NewFileService() FileService {
	return FileService(0)
}

func (s FileService) ValidName(fileName string) bool {
	return !strings.HasPrefix(fileName, " ") &&
		!strings.HasSuffix(fileName, " ") &&
		validFilenameRegexp.Match([]byte(fileName))
}

func (s FileService) PurgeFileName(fileName string) string {
	trimedFileName := strings.TrimSpace(fileName)

	return invalidFilenameCharacters.ReplaceAllString(trimedFileName, "_")
}

func (s FileService) CreateFilename(music *model.Music) string {
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(music.Author.Name)

	if len(music.Contributors) != 0 {
		nameBuilder.WriteString(" (")

		for i, contributors := range music.Contributors {
			nameBuilder.WriteString(contributors.Name)

			if len(music.Contributors)-1 != i {
				nameBuilder.WriteString(", ")
			}
		}

		nameBuilder.WriteString(")")
	}

	nameBuilder.WriteString(fmt.Sprintf(" - %s [%x]", music.Name, time.Now().UnixMicro()))

	return s.PurgeFileName(nameBuilder.String())
}

func (s FileService) MoveTo(from, to string) error {
	filename := pathModule.Base(from)
	newPath := pathModule.Join(to, filename)

	if f, err := os.Open(newPath); err == nil {
		f.Close()

		return fmt.Errorf("file (%s) already exists in (%s)", filename, to)
	}

	return os.Rename(from, newPath)
}

func (s FileService) IsExists(path string) bool {
	if f, err := os.Open(path); err == nil {
		f.Close()

		return true
	}

	return false
}

func (s FileService) DownloadImageToTemp(imageUri string) {
	if file, err := os.Open(imageUri); err == nil {
		defer file.Close()
		// Image found on the PC locally
	} else {
		// Image must be found on the web

	}
}

func (s FileService) GetImageConfig(imagePath string) (int, int, string, error) {
	ext := pathModule.Ext(imagePath)
	var width, height int
	picturePath := pathModule.Join(settings.Global.App.ImagesLocation, imagePath)
	picFile, err := os.Open(picturePath)
	if err != nil {
		return 0, 0, ext, err
	}

	// Get Image's dimensions, so it can be scaled down, if necessary
	switch ext {
	case "webp":
		config, err := webp.DecodeConfig(picFile)
		if err != nil {
			return 0, 0, ext, err
		}

		width = config.Width
		height = config.Height
	case "jpeg":
		fallthrough
	case "jpg":
		fallthrough
	case "png":
		config, _, err := image.DecodeConfig(picFile)
		if err != nil {
			return 0, 0, ext, err
		}

		width = config.Width
		height = config.Height
	}

	// Get max size - actual size ration
	var ratio float32
	if width < height {
		ratio = THUMBNAIL_SIZE / float32(width)
	} else {
		ratio = THUMBNAIL_SIZE / float32(height)
	}

	// Calculate new size based on the ratio
	var newWidth int
	var newHeight int
	if ratio < 1 {
		newWidth = max(int(float32(width)*ratio), THUMBNAIL_SIZE)
		newHeight = max(int(float32(height)*ratio), THUMBNAIL_SIZE)
	}

	return newWidth, newHeight, ext, nil
}
