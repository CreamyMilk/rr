package models

import (
	"aytp/graph/model"

	uuid "github.com/satori/go.uuid"
)

type Course struct {
	Base
	Name             string
	Icon             string
	Duration         string
	Description      string
	SchoolRef        string
	Cost             int
	CourseLevel      model.CourseLevel
	CourseType       model.CourseType
	Enrolled         int
	RoomId           string
	DefaultSessionId string
	Chapters         []*Chapter    `gorm:"foreignKey:CourseID"`
	Files            []*File       `gorm:"foreignKey:CourseID"`
	Assignment       []*Assignment `gorm:"foreignKey:CourseID"`
}

type Chapter struct {
	Base
	CourseID    uuid.UUID `gorm:"type:uuid"`
	Title       string
	Description string
	Sections    []*Section `gorm:"foreignKey:ChapterID"`
}

type Section struct {
	Base
	ChapterID uuid.UUID `gorm:"type:uuid"`
	Heading   string
	Content   string
}

type File struct {
	CourseID uuid.UUID `gorm:"type:uuid"`
	Name     string
	Url      string
}

type Enroll struct {
	StudentID uuid.UUID `gorm:"type:uuid"`
	CourseID  uuid.UUID `gorm:"type:uuid"`
}

type CourseAssign struct {
	TutorID  uuid.UUID `gorm:"type:uuid"`
	CourseID uuid.UUID `gorm:"type:uuid"`
}

func (s Section) CreateToGraphData() *model.Section {
	return &model.Section{
		Heading: s.Heading,
		Content: s.Content,
	}
}

func (ch Chapter) CreateToGraphData() *model.Chapter {
	// gqlSections := func(sections []*Section) []*model.Section {
	// 	gql := make([]*model.Section, len(sections))
	// 	for idx, section := range sections {
	// 		gql[idx] = section.CreateToGraphData()
	// 	}
	// 	return gql
	// }(ch.Sections)

	return &model.Chapter{
		ID:          ch.ID.String(),
		Title:       ch.Title,
		Description: ch.Description,
	}
}

func (f File) CreateToGraphData() *model.File {
	return &model.File{
		Name: f.Name,
		URL:  f.Url,
	}
}

func (c Course) CreateToGraphData() *model.Course {
	return &model.Course{
		ID:          c.ID.String(),
		Name:        c.Name,
		Icon:        c.Icon,
		Duration:    c.Duration,
		Cost:        c.Cost,
		Description: c.Description,
		Enrolled:    c.Enrolled,
		Level:       c.CourseLevel,
		Type:        c.CourseType,
	}
}
