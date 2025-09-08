package model

import "yotudo/src/database/entity"

type Image struct {
	Id   int64
	Name string
}

func (i *Image) FromEntity(e *entity.Image) *Image {
	i.Id = e.Id
	i.Name = e.Name

	return i
}

func (i *Image) ToPossiblyNewImage() *PossiblyNewImage {
	return &PossiblyNewImage{Id: &i.Id, Name: i.Name}
}

type PossiblyNewImage struct {
	Id     *int64
	Name   string
	Source *string
}

func (i *PossiblyNewImage) FromImage(e *Image) *PossiblyNewImage {
	i.Id = &e.Id
	i.Name = e.Name

	return i
}

func (i *PossiblyNewImage) HasId() bool {
	return i != nil && i.Id != nil
}

func (i *PossiblyNewImage) HasName() bool {
	return i != nil && i.Name != ""
}

func (i *PossiblyNewImage) HasSource() bool {
	return i != nil && i.Source != nil
}

func (i *PossiblyNewImage) IsNew() bool {
	return i != nil && i.Id == nil
}
