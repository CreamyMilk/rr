package models

import (
	"aytp/graph/model"
)

type Emailz struct {
	Base
	From    string
	To      string
	Subject string
	Body    string
}

func (e Emailz) CreateToGraphData() *model.Email {
	return &model.Email{
		ID:      e.ID.String(),
		From:    e.From,
		Subject: e.Subject,
		To:      e.To,
		Body:    e.Body,
	}
}
