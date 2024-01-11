package events

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
)

func AddEvent(input model.AddEventsInput) (bool, error) {
	admin, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	newEvent := models.Event{
		Title:       input.Title,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Date:        input.Date,
		ImageURL:    input.ImageURL,
		Location:    input.Location,
		Ref:         admin.SchoolRef,
		Link:        input.Link,
	}
	err = utils.DB.Create(&newEvent).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}
func EditEvent(input model.EditEventsInput) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
		
	}
	event, err := FetchOneEvent(input.EventID)
	if err != nil {
		return false, errors.New("no event found")
	}

	event.Title = input.Title
	event.Description = input.Description
	event.StartTime = input.StartTime
	event.EndTime = input.EndTime
	event.Location = input.Location
	event.ImageURL = input.ImageURL
	event.Date = input.Date
	event.Link = input.Link

	err = utils.DB.Model(&event).Updates(event).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func FetchAllEvents(ref string) ([]*model.Event, error) {
	var events []models.Event
	err := utils.DB.Where(&models.Event{Ref: ref}).Find(&events).Error
	if err != nil {
		return nil, err
	}
	gqlEvents := make([]*model.Event, len(events))
	for i, v := range events {
		gqlEvents[i] = v.CreateToGraphData()
	}
	return gqlEvents, nil
}

func FetchOneEvent(id string) (*model.Event, error) {
	var event models.Event
	err := utils.DB.Where("id=?", id).Find(&event).Error
	if err != nil {
		return nil, err
	}
	return event.CreateToGraphData(), nil
}

func DeleteEvent(token string, id string) (bool, error) {
	_, err := engine.FetchAdminByToken(token)
	if err != nil {
		return false, err
	}
	err = utils.DB.Where("id = ?", id).Delete(&models.Event{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
