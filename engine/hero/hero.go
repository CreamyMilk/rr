package hero

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
)

func AddHero(input model.AddHeroInput) (bool, error) {
	admin, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	newHero := models.Hero{
		Title:    input.Title,
		Subtitle: input.Subtitle,
		Banner1:  input.Banner1,
		Banner2:  input.Banner2,
		Ref:      admin.SchoolRef,
	}
	err = utils.DB.Create(&newHero).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func FetchHeroByRef(ref string) (*model.Hero, error) {
	var hero models.Hero
	err := utils.DB.Where(&models.Hero{Ref: ref}).Find(&hero).Error
	if err != nil {
		return nil, err
	}
	gqlHero := hero.CreateToGraphData()
	return gqlHero, nil
}

func FetchHeroById(id string) (*model.Hero, error) {
	var hero models.Hero
	err := utils.DB.Where("id=?", id).Find(&hero).Error
	if err != nil {
		return nil, err
	}
	return hero.CreateToGraphData(), nil
}

func EditHero(input model.EditHeroInput) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	hero, err := FetchHeroById(input.HeroID)
	if err != nil {
		return false, errors.New("no hero info found")
	}

	hero.Title = input.Title
	hero.Subtitle = input.Subtitle
	hero.Banner1 = input.Banner1
	hero.Banner2 = input.Banner2

	err = utils.DB.Model(&hero).Updates(hero).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
