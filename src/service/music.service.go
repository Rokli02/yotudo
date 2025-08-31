package service

import (
	"fmt"
	"os"
	"path"
	"strings"
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"
)

type MusicService struct {
	musicRepository       *repository.Music
	authorRepository      *repository.Author
	contributorRepository *repository.Contributor
	youtubeDLService      *YoutubeDLService
	fileService           FileService
}

func NewMusicService(
	musicRepository *repository.Music,
	authorRepository *repository.Author,
	contributorRepository *repository.Contributor,
	youtubeDLService *YoutubeDLService,
	fileService FileService,
) *MusicService {
	return &MusicService{
		musicRepository:       musicRepository,
		authorRepository:      authorRepository,
		contributorRepository: contributorRepository,
		youtubeDLService:      youtubeDLService,
		fileService:           fileService,
	}
}

func (c *MusicService) GetManyByPagination(filter string, statusId int, page *model.Page, sort []model.Sort) *model.Pagination[[]model.Music] {
	musics, totalCount := c.musicRepository.FindByPageAndStatus(statusId, filter, page, sort)

	return &model.Pagination[[]model.Music]{
		Data:  musics,
		Count: totalCount,
	}
}

func (c *MusicService) GetById(id int64) (*model.Music, error) {
	if id < 0 {
		return nil, fmt.Errorf("valid id must be given")
	}

	return c.musicRepository.FindById(id)
}

func (c *MusicService) Save(newMusic *model.NewMusic) (*model.Music, error) {
	if err := c.processMusicAuthor(newMusic); err != nil {
		return nil, err
	}

	if err := c.processMusicContributors(newMusic); err != nil {
		return nil, err
	}

	switch newMusic.PicType {
	case "none":
		newMusic.PicFilename = ""
	case "thumbnail":
		thumbnailUrl, err := c.youtubeDLService.GetVideoThumbnailUrl(newMusic.Url)
		if err != nil {
			logger.Error(err)
		} else {
			tmpFilename, err := c.fileService.DownloadImageFromWeb(thumbnailUrl)
			if err != nil {
				logger.Error(err)
			} else {
				newMusic.PicFilename = tmpFilename
			}
		}
	case "web":
		if newMusic.PicFilename != "" {
			tmpFilename, err := c.fileService.DownloadImageFromWeb(newMusic.PicFilename)
			if err != nil {
				logger.Error(err)
			} else {
				newMusic.PicFilename = tmpFilename
			}
		}
	case "local":
		if newMusic.PicFilename != "" {
			tmpFilename, err := c.fileService.CopyImageFromFS(newMusic.PicFilename)
			if err != nil {
				logger.Error(err)
			} else {
				newMusic.PicFilename = tmpFilename
			}
		}
	}

	insertedId, err := c.musicRepository.SaveOne(newMusic)
	if err != nil {
		return nil, err
	}

	savedMusic, err := c.musicRepository.FindById(insertedId)
	if err != nil {
		return nil, err
	}

	if savedMusic.PicFilename != nil {
		err := c.fileService.MoveTo(path.Join(settings.Global.App.TempLocation, *savedMusic.PicFilename), settings.Global.App.ImagesLocation)
		if err != nil {
			logger.Error(err)
		}
	}

	return savedMusic, nil
}

func (c *MusicService) Update(updateMusic *model.UpdateMusic) (*model.Music, error) {
	if updateMusic == nil {
		return nil, errors.ErrNotReceivedInputs
	}

	musicEntity, err := c.musicRepository.FindById_Entity(updateMusic.Id)
	if err != nil {
		logger.Warning(err)

		return nil, err
	}

	if err := c.processMusicAuthor(updateMusic); err != nil {
		return nil, err
	}

	if err := c.processMusicContributors(updateMusic); err != nil {
		return nil, err
	}

	if musicEntity.Filename != nil {
		updateMusic.Filename = *musicEntity.Filename
	}

	if updateMusic.PicFilename == "" {
		if musicEntity.PicFilename != nil {
			updateMusic.PicFilename = *musicEntity.PicFilename
		}
	}

	oldPicFilename := musicEntity.PicFilename
	switch updateMusic.PicType {
	case "none":
		updateMusic.PicFilename = ""
	case "thumbnail":
		thumbnailUrl, err := c.youtubeDLService.GetVideoThumbnailUrl(updateMusic.Url)
		if err != nil {
			logger.Error(err)
		} else {
			tmpFilename, err := c.fileService.DownloadImageFromWeb(thumbnailUrl)
			if err != nil {
				logger.Error(err)
			} else {
				updateMusic.PicFilename = tmpFilename
			}
		}
	case "web":
		if updateMusic.PicFilename != "" {
			tmpFilename, err := c.fileService.DownloadImageFromWeb(updateMusic.PicFilename)
			if err != nil {
				logger.Error(err)
			} else {
				updateMusic.PicFilename = tmpFilename
			}
		}
	case "local":
		if musicEntity.PicFilename != nil && strings.HasSuffix(updateMusic.PicFilename, *musicEntity.PicFilename) {
			updateMusic.PicFilename = *musicEntity.PicFilename

			break
		}

		if updateMusic.PicFilename != "" {
			tmpFilename, err := c.fileService.CopyImageFromFS(updateMusic.PicFilename)
			if err != nil {
				logger.Error(err)
			} else {
				updateMusic.PicFilename = tmpFilename
			}
		}
	}

	updatedMusic, err := c.musicRepository.UpdateOne(updateMusic.Id, updateMusic)
	if err != nil {
		return nil, err
	}

	if oldPicFilename != nil && (updatedMusic.PicFilename == nil || *oldPicFilename != *updatedMusic.PicFilename) {
		os.Remove(path.Join(settings.Global.App.ImagesLocation, *oldPicFilename))
	}

	if updatedMusic.PicFilename != nil && (oldPicFilename == nil || *oldPicFilename != *updatedMusic.PicFilename) {
		err := c.fileService.MoveTo(path.Join(settings.Global.App.TempLocation, *updatedMusic.PicFilename), settings.Global.App.ImagesLocation)
		if err != nil {
			logger.Error(err)
		}
	}

	return updatedMusic, nil
}

func (c *MusicService) Delete(id int64) error {
	if deleted, err := c.musicRepository.DeleteOne(id); err != nil || !deleted {
		return errors.ErrUnableToDelete
	}

	return nil
}

func (c *MusicService) processMusicAuthor(music model.OptionalAuthorGetter) error {
	author := music.GetOptionalAuthor()

	if author.Id == nil {
		logger.Debug("No author id was passed when attempted to save a Music")

		if author.Name == nil {
			return fmt.Errorf("no author was given")
		}

		savedAuthor, err := c.authorRepository.SaveOne(*author.Name)
		if err != nil {
			return err
		}

		author.Id = &savedAuthor.Id
	}

	return nil
}

func (c *MusicService) processMusicContributors(music model.OptionalContributorsAccessor) error {
	contributors := music.GetOptionalContributors()

	if len(contributors) != 0 {
		var newContributorAuthors []string = nil

		// Gets those contributors which were passed down without an id (new Authors)
		for _, contributor := range contributors {
			if contributor.Id == nil {
				if contributor.Name == nil {
					logger.Error("author without 'name' or 'id' was passed")

					continue
				}

				if newContributorAuthors == nil {
					newContributorAuthors = make([]string, 0, 1)
				}

				newContributorAuthors = append(newContributorAuthors, *contributor.Name)
			}
		}

		// Found at least one contributor that got passed without id
		if newContributorAuthors != nil {
			logger.Debug("Received contributor(s) without id")

			if savedNewContributors, err := c.authorRepository.SaveMany(newContributorAuthors); err != nil {
				logger.Error(err)

				return err
			} else {
				// Blend freshly inserted authors with the ones already existing ones
				newContributors := make([]model.OptionalAuthor, 0, len(contributors))
				for _, contributor := range contributors {
					if contributor.Id != nil {
						newContributors = append(newContributors, contributor)
					}
				}

				for _, contributor := range savedNewContributors {
					newContributors = append(newContributors, model.OptionalAuthor{Id: &contributor.Id, Name: &contributor.Name})
				}

				// Add blended slice to newMusic object
				music.SetOptionalContributors(newContributors)
			}
		}
	}

	return nil
}
