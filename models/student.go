package models

import (
	"aytp/graph/model"

	uuid "github.com/satori/go.uuid"
)

type Student struct {
	Base
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Points       int
	Password     string
	SchoolID     uuid.UUID
	SchoolName   string
	ProfilePhoto string
	Status       model.StudentStatus
	Enrolls      []*Enroll `gorm:"foreignKey:StudentID"`
	GithubToken  *string
}
type Badge struct {
	Base
	StudentID   uuid.UUID `gorm:"type:uuid"`
	PhotoUrl    string
	Redeemed    bool
	PointsIndex int
}

func (s Student) CreateToGraphData() *model.Student {
	return &model.Student{
		ID:       s.ID.String(),
		Name:     s.Name,
		Email:    s.Email,
		Points:   s.Points,
		JoinedAt: s.CreatedAt,
		Status:   s.Status,
		Enrolls:  len(s.Enrolls),
	}
}
