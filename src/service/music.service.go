package service

import (
	"fmt"
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type MusicService struct {
	musicRepository       *repository.Music
	authorRepository      *repository.Author
	contributorRepository *repository.Contributor
}

func NewMusicService(
	musicRepository *repository.Music,
	authorRepository *repository.Author,
	contributorRepository *repository.Contributor,
) *MusicService {
	return &MusicService{
		musicRepository:       musicRepository,
		authorRepository:      authorRepository,
		contributorRepository: contributorRepository,
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

	if newMusic.PicFilename != "" && newMusic.PicFilename != "thumbnail" {
		// TODO: Átmozgatni / Letölteni a képet az URI alapján az 'imgs' mappába
	}

	insertedId, err := c.musicRepository.SaveOne(newMusic)
	if err != nil {
		return nil, err
	}

	return c.musicRepository.FindById(insertedId)
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
	} else if musicEntity.PicFilename == nil || updateMusic.PicFilename != *musicEntity.PicFilename {
		// TODO: Átmozgatni / Letölteni a képet az URI alapján az 'imgs' mappába
	}

	return c.musicRepository.UpdateOne(updateMusic.Id, updateMusic)
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
