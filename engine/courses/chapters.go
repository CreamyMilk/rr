package courses

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func AddChapter(input model.NewChapter) (bool, error) {
	_, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, errors.New("course does not exist")
	}

	err = utils.DB.Model(&models.Chapter{}).Where(&models.Chapter{Title: input.Title}).First(nil).Error
	if err == nil {
		return false, errors.New("course already contains that chapter")
	}

	newChapter := models.Chapter{
		CourseID:    uuid.FromStringOrNil(input.CourseID),
		Title:       input.Title,
		Description: input.Description,
	}
	err = utils.DB.Save(&newChapter).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetCourseChapters(token, courseId string) ([]*model.Chapter, error) {
	if !CheckEnrollMent(token, courseId) {
		return nil, errors.New("please enroll to learn this course")
	}

	chapters, err := engine.FetchChaptersByCourseID(courseId)
	if err != nil {
		return nil, err
	}
	var gqlChapters = make([]*model.Chapter, len(chapters))
	for i, chapter := range chapters {
		gqlChapters[i] = chapter.CreateToGraphData()
	}
	return gqlChapters, nil
}
