package service

import (
	"fmt"
	"os"
	"path"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"
)

type MusicService struct{}

var GlobalMusicService *MusicService = nil

func (c *MusicService) GetManyByPagination(filter string, statusId int, page *model.Page, sort []model.Sort) *model.Pagination[[]model.Music] {
	musics, totalCount := repository.GlobalMusicRepository.FindByPageAndStatus(statusId, filter, page, sort)

	return &model.Pagination[[]model.Music]{
		Data:  musics,
		Count: totalCount,
	}
}

func (c *MusicService) GetById(id int64) (*model.Music, error) {
	if id < 0 {
		return nil, fmt.Errorf("valid id must be given")
	}

	return repository.GlobalMusicRepository.FindById(id)
}

func (c *MusicService) Save(newMusic *model.NewMusic) (*model.Music, error) {
	if err := c.processMusicAuthor(newMusic); err != nil {
		return nil, err
	}

	if err := c.processMusicContributors(newMusic); err != nil {
		return nil, err
	}

	switch newMusic.PicType {
	case "":
		fallthrough
	case "none":
		newMusic.Image = nil
	case "thumbnail":
		thumbnailUrl, err := GlobalYoutubeDLService.GetVideoThumbnailUrl(newMusic.Url)
		if err != nil {
			logger.Error(err)
			newMusic.Image = nil
			break
		}

		if foundImage, err := GlobalImageService.GetBySource(thumbnailUrl); err != nil {
			tmpFilename, err := GlobalFileService.DownloadImageFromWeb(thumbnailUrl)
			if err != nil {
				logger.Error(err)
				newMusic.Image = nil
			} else {
				newMusic.Image = &model.PossiblyNewImage{Name: tmpFilename, Source: &thumbnailUrl}
			}
		} else {
			logger.DebugF("Found an already existing source(id=%d, name=%s)", foundImage.Id, foundImage.Name)
			newMusic.Image = foundImage.ToPossiblyNewImage()
		}
	case "web":
		if !newMusic.Image.HasName() {
			break
		}

		if foundImage, err := GlobalImageService.GetBySource(newMusic.Image.Name); err != nil {
			tmpFilename, err := GlobalFileService.DownloadImageFromWeb(newMusic.Image.Name)
			if err != nil {
				logger.Error(err)
				newMusic.Image = nil
			} else {
				newMusic.Image = &model.PossiblyNewImage{Name: tmpFilename, Source: &newMusic.Image.Name}
			}
		} else {
			logger.DebugF("Found an already existing source(id=%d, name=%s)", foundImage.Id, foundImage.Name)
			newMusic.Image = foundImage.ToPossiblyNewImage()
		}
	case "local":
		if !newMusic.Image.HasName() {
			break
		}

		filename := GlobalFileService.GetFilename(newMusic.Image.Name)

		if image, err := GlobalImageService.GetByName(filename); err != nil {
			tmpFilename, err := GlobalFileService.CopyImageFromFS(newMusic.Image.Name)
			if err != nil {
				logger.Error(err)
				newMusic.Image = nil
				break
			}

			newMusic.Image.Name = tmpFilename
		} else {
			newMusic.Image.FromImage(image)
		}
	}

	insertedId, err := repository.GlobalMusicRepository.SaveOne(newMusic)
	if err != nil {
		return nil, err
	}

	savedMusic, err := repository.GlobalMusicRepository.FindById(insertedId)
	if err != nil {
		return nil, err
	}

	if savedMusic.Image != nil && newMusic.Image.IsNew() {
		err := GlobalFileService.MoveTo(path.Join(settings.Global.App.TempLocation, savedMusic.Image.Name), settings.Global.App.ImagesLocation)
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

	musicEntity, err := repository.GlobalMusicRepository.FindById_Entity(updateMusic.Id)
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

	switch updateMusic.PicType {
	case "":
		if musicEntity.ImageId != nil {
			updateMusic.Image = &model.PossiblyNewImage{Id: musicEntity.ImageId}
		}
	case "none":
		updateMusic.Image = nil
	case "thumbnail":
		thumbnailUrl, err := GlobalYoutubeDLService.GetVideoThumbnailUrl(updateMusic.Url)
		if err != nil {
			logger.Error(err)
			updateMusic.Image = nil
			break
		}

		if foundImage, err := GlobalImageService.GetBySource(thumbnailUrl); err != nil {
			tmpFilename, err := GlobalFileService.DownloadImageFromWeb(thumbnailUrl)
			if err != nil {
				logger.Error(err)
				updateMusic.Image = nil
			} else {
				updateMusic.Image = &model.PossiblyNewImage{Name: tmpFilename, Source: &thumbnailUrl}
			}
		} else {
			logger.DebugF("Found an already existing source(id=%d, name=%s)", foundImage.Id, foundImage.Name)
			updateMusic.Image = foundImage.ToPossiblyNewImage()
		}
	case "web":
		if !updateMusic.Image.HasName() {
			break
		}

		if foundImage, err := GlobalImageService.GetBySource(updateMusic.Image.Name); err != nil {
			tmpFilename, err := GlobalFileService.DownloadImageFromWeb(updateMusic.Image.Name)
			if err != nil {
				logger.Error(err)
				updateMusic.Image = nil
			} else {
				updateMusic.Image = &model.PossiblyNewImage{Name: tmpFilename, Source: &updateMusic.Image.Name}
			}
		} else {
			logger.DebugF("Found an already existing source(id=%d, name=%s)", foundImage.Id, foundImage.Name)
			updateMusic.Image = foundImage.ToPossiblyNewImage()
		}
	case "local":
		if !updateMusic.Image.HasName() {
			break
		}

		filename := GlobalFileService.GetFilename(updateMusic.Image.Name)

		if image, err := GlobalImageService.GetByName(filename); err != nil {
			tmpFilename, err := GlobalFileService.CopyImageFromFS(updateMusic.Image.Name)
			if err != nil {
				logger.Error(err)
				updateMusic.Image = nil
				break
			}

			updateMusic.Image.Name = tmpFilename
		} else {
			updateMusic.Image.FromImage(image)
		}
	}

	isNewImage := updateMusic.Image.IsNew()
	if isNewImage {
		if savedImage, err := GlobalImageService.Save(updateMusic.Image); err != nil {
			logger.Error(err)
			updateMusic.Image = nil
		} else {
			updateMusic.Image.FromImage(savedImage)
		}
	}

	updatedMusic, err := repository.GlobalMusicRepository.UpdateOne(updateMusic.Id, updateMusic)
	if err != nil {
		return nil, err
	}

	if updatedMusic.Image == nil { // Nincs elmentett kép
		logger.Debug("Nincs elmentett kép")
		if musicEntity.ImageId != nil { // Volt elmentett kép
			logger.Debug("Volt elmentett kép")
			if err := GlobalImageService.DeleteImageIfUnused(*musicEntity.ImageId); err != nil {
				logger.Error(err)
			}
		}
	} else { // Van elmentett kép
		logger.Debug("Van elmentett kép")
		if musicEntity.ImageId == nil { // Nem volt elmentett kép
			logger.Debug("Nem volt elmentett kép")
			if isNewImage { // A mentett kép új, át kell mozgatni
				logger.Debug("A mentett kép új, át kell mozgatni")
				if err := GlobalFileService.MoveTo(path.Join(settings.Global.App.TempLocation, updatedMusic.Image.Name), settings.Global.App.ImagesLocation); err != nil {
					logger.Error(err)
				}
			}
		} else if updatedMusic.Image.Id != *musicEntity.ImageId { // Volt elmentett kép ÉS nem ugyanaz
			logger.Debug("Volt elmentett kép ÉS nem ugyanaz")
			if err := GlobalImageService.DeleteImageIfUnused(*musicEntity.ImageId); err != nil {
				logger.Error(err)
			} else if isNewImage { // A mentett kép új, át kell mozgatni
				logger.Debug("A mentett kép új, át kell mozgatni")
				if err := GlobalFileService.MoveTo(path.Join(settings.Global.App.TempLocation, updatedMusic.Image.Name), settings.Global.App.ImagesLocation); err != nil {
					logger.Error(err)
				}
			}
		}
	}

	return updatedMusic, nil
}

func (c *MusicService) Delete(id int64) error {
	var musicEntity *entity.Music

	if _musicEntity, err := repository.GlobalMusicRepository.FindById_Entity(id); err != nil {
		logger.Error(err)
		return err
	} else {
		musicEntity = _musicEntity
	}

	if deleted, err := repository.GlobalMusicRepository.DeleteOne(id); err != nil || !deleted {
		logger.Error(err)
		return errors.ErrUnableToDelete
	}

	if musicEntity.ImageId != nil {
		if err := GlobalImageService.DeleteImageIfUnused(*musicEntity.ImageId); err != nil {
			logger.Error(err)
		}
	}

	if musicEntity.Filename != nil {
		if err := os.Remove(path.Join(settings.Global.App.MusicsLocation, *musicEntity.Filename)); err != nil {
			logger.Error(err)
		}
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

		savedAuthor, err := repository.GlobalAuthorRepository.SaveOne(*author.Name)
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

			if savedNewContributors, err := repository.GlobalAuthorRepository.SaveMany(newContributorAuthors); err != nil {
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
