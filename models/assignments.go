package models

import (
	"aytp/graph/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Assignment struct {
	Base
	CourseID    uuid.UUID `gorm:"type:uuid"`
	Title       string
	Description string
	DueDate     time.Time
	Files       []*AssignmentFile `gorm:"foreignKey:AssignmentID"`
	Submissions []*Submission     `gorm:"foreignKey:AssignmentID"`
}

type SubmissionFile struct {
	SubmissionID uuid.UUID `gorm:"type:uuid"`
	Name         string
	Url          string
}

type AssignmentFile struct {
	AssignmentID uuid.UUID `gorm:"type:uuid"`
	Name         string
	Url          string
}

type Submission struct {
	Base
	CourseID       uuid.UUID `gorm:"type:uuid"`
	AssignmentID   uuid.UUID `gorm:"type:uuid"`
	StudentName    string
	StudentID      uuid.UUID `gorm:"type:uuid"`
	SubmissionDate time.Time
	Note           string
	Remark         string
	Marks          int
	Files          []*SubmissionFile `gorm:"foreignKey:SubmissionID"`
}

func (a Assignment) CreateToGraphData() *model.Assignment {
	gqlFiles := func(files []*AssignmentFile) []*model.File {
		gql := make([]*model.File, len(files))
		for idx, file := range files {
			gql[idx] = &model.File{
				Name: file.Name,
				URL:  file.Url,
			}
		}
		return gql
	}(a.Files)

	return &model.Assignment{
		ID:          a.ID.String(),
		CourseID:    a.CourseID.String(),
		Title:       a.Title,
		Description: a.Description,
		DueDate:     a.DueDate,
		Files:       gqlFiles,
	}
}

func (s Submission) CreateToGraphData() *model.Submission {
	gqlFiles := func(files []*SubmissionFile) []*model.File {
		gql := make([]*model.File, len(files))
		for idx, file := range files {
			gql[idx] = &model.File{
				Name: file.Name,
				URL:  file.Url,
			}
		}
		return gql
	}(s.Files)

	return &model.Submission{
		ID:             s.ID.String(),
		CourseID:       s.CourseID.String(),
		StudentName:    s.StudentName,
		Note:           &s.Note,
		Remark:         &s.Remark,
		Marks:          &s.Marks,
		SubmissionDate: s.SubmissionDate,
		Files:          gqlFiles,
	}
}
