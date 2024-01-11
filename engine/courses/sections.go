package courses

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func GetChapterSections(token, chapterId string) ([]*model.Section, error) {
	chapter, err := engine.FetchChapterByID(chapterId)
	if err != nil {
		return nil, err
	}

	if !CheckEnrollMent(token, chapter.CourseID.String()) {
		return nil, errors.New("please enroll to learn this course")
	}

	var gqlSections = make([]*model.Section, len(chapter.Sections))
	for i, section := range chapter.Sections {
		gqlSections[i] = section.CreateToGraphData()
	}
	return gqlSections, nil
}

func AddSection(input model.NewSection) (bool, error) {
	//fetch tutor by token
	//ensure tutor teaches is assigned course
	_, err := engine.FetchSectionsByChapterID(input.ChapterID)
	if err != nil {
		return false, errors.New("chapter does not exist")
	}

	var section models.Section
	err = utils.DB.Model(&models.Section{}).Where(&models.Section{Heading: input.Heading}).First(&section).Error
	if err == nil {
		section.Content = fmt.Sprint(section.Content, " ", input.Content)
	} else {
		section = models.Section{
			ChapterID: uuid.FromStringOrNil(input.ChapterID),
			Heading:   input.Heading,
			Content:   input.Content,
		}
	}
	err = utils.DB.Save(&section).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
