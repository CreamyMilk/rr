package models

import (
	"aytp/graph/model"
)

type Hero struct {
	Base
	Title    string
	Subtitle string
	Banner1  string
	Banner2  string
	Ref      string
}

func (e Hero) CreateToGraphData() *model.Hero {
	return &model.Hero{
		ID:       e.ID.String(),
		Title:    e.Title,
		Subtitle: e.Subtitle,
		Banner1:  e.Banner1,
		Banner2:  e.Banner2,
	}
}
