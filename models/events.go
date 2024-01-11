package models

import (
	"aytp/graph/model"
	"time"
)

type Event struct {
	Base
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Date        time.Time
	ImageURL    string
	Location    string
	Ref         string
	Link        string
}

func (e Event) CreateToGraphData() *model.Event {
	return &model.Event{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		ImageURL:    e.ImageURL,
		Date:        e.Date,
		Link:        e.Link,
	}
}
