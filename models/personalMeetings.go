package models

import (
	"aytp/graph/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PersonalMeetings struct {
	Base
	TutorId     uuid.UUID `gorm:"type:uuid"`
	TutorName   string
	Tutorcourse string
	CourseId    uuid.UUID `gorm:"type:uuid"`
	UserId      uuid.UUID `gorm:"type:uuid"`
	UserName    string
	Year        int
	Month       time.Month
	Day         int
	Hour        int
	EndTime     int
	TutorLink   string
	StudentLink string
	Status      model.PersonalMeetingStatus
}
