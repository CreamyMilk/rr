package models

import (
	"aytp/graph/model"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Quiz struct {
	Base
	Title     string
	CourseID  string
	Questions []*Question `gorm:"foreignKey:QuizID"`
}

type QuizComplete struct {
	QuizID    uuid.UUID `gorm:"type:uuid"`
	StudentID uuid.UUID `gorm:"type:uuid"`
	CourseID  uuid.UUID `gorm:"type:uuid"`
	SchoolID  uuid.UUID `gorm:"type:uuid"`
	Marks     int
	Total     int
	Answers   string
	QuizStatus string
}

type Question struct {
	// Base
	QuizID  uuid.UUID `gorm:"type:uuid"`
	Text    string
	Choices string
	Answer  string
}

func (q Quiz) CreateToGraphData() *model.Quiz {
	gqlQuestions := func(questions []*Question) []*model.Question {
		gql := make([]*model.Question, len(questions))
		for idx, question := range questions {
			gql[idx] = &model.Question{
				Text:    question.Text,
				Choices: strings.Split(question.Choices, ","),
			}
		}
		return gql
	}(q.Questions)

	return &model.Quiz{
		ID:        q.ID.String(),
		Title:     q.Title,
		CourseID:  q.CourseID,
		Questions: gqlQuestions,
	}
}
