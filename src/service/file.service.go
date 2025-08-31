package service

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

/*
Saves the given 'imageUri' into a file in the local temp folder and returns its name.
*/
func (s FileService) DownloadImageFromWeb(imageUri string) (string, error) {
	imageUrl, err := url.Parse(imageUri)
	if err != nil {
		return "", err
	}

	response, err := http.Get(imageUri)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	fileExtension := "jpg"
	contentType := response.Header.Get("Content-Type")
	if contentType != "" {
		if extractedExtension, found := strings.CutSuffix(contentType, "image/"); found {
			fileExtension = extractedExtension
		}
	}

	purgedFilename := s.PurgeFileName(imageUrl.Path)
	if purgedFilenameExt := pathModule.Ext(purgedFilename); purgedFilenameExt != "" {
		if tmpPurgedFilename, found := strings.CutSuffix(purgedFilename, purgedFilenameExt); found {
			purgedFilename = tmpPurgedFilename
		}
	}

	filename := fmt.Sprintf("%s.%s", purgedFilename, fileExtension)

	createdImageFile, err := os.Create(pathModule.Join(settings.Global.App.TempLocation, filename))
	if err != nil {
		return "", err
	}
	defer createdImageFile.Close()

	_, err = io.Copy(createdImageFile, response.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

/*
Copies the given 'imagePath' into the local temp folder and returns its name.
*/
func (s FileService) CopyImageFromFS(imagePath string) (string, error) {
	filename := pathModule.Base(imagePath)
	if f, err := os.Open(pathModule.Join(settings.Global.App.TempLocation, filename)); err == nil {
		// Már megtalálható a fájl az ideiglenes mappába, térjünk vissza a nevével
		f.Close()

		return filename, nil
	}

	ogFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer ogFile.Close()

	if f, err := os.Create(pathModule.Join(settings.Global.App.TempLocation, filename)); err != nil {
		return "", err
	} else {
		defer f.Close()

		io.Copy(f, ogFile)
	}

	return filename, nil
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
	case ".webp":
		config, err := webp.DecodeConfig(picFile)
		if err != nil {
			return 0, 0, ext, err
		}

		width = config.Width
		height = config.Height
	case ".jpeg":
		fallthrough
	case ".jpg":
		fallthrough
	case ".png":
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
	var newWidth int = width
	var newHeight int = height
	if ratio < 1 {
		newWidth = max(int(float32(width)*ratio), THUMBNAIL_SIZE)
		newHeight = max(int(float32(height)*ratio), THUMBNAIL_SIZE)
	}

	return newWidth, newHeight, ext, nil
}

func (s FileService) HasExecutable() bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*4)
	defer cancelCtx()

	return exec.CommandContext(ctx, settings.Global.App.FFMPEGLocation, "-version").Run() == nil
}
