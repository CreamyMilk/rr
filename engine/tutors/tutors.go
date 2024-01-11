package tutors

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
)

func Login(input model.LoginInput) (*model.AuthPayload, error) {
	tutor, err := engine.FetchTutor(&models.Tutor{Email: input.Email})
	if err != nil {
		return nil, engine.ErrWrongCreds
	}

	if !utils.CompareHashedString(tutor.Password, input.Password) {
		return nil, engine.ErrWrongCreds
	}
	return engine.GenerateToken(tutor.ID, tutor.Name)
}

func AddNewTutor(input model.NewTutor) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	err = engine.FetchUserByEmail(input.Email)
	if err == nil {
		return false, errors.New("email already registered")
	}

	pass, err := utils.HashString("password")
	if err != nil {
		return false, nil
	}
	newTutor := models.Tutor{
		Email:    input.Email,
		Name:     input.Name,
		Password: pass,
	}
	err = utils.DB.Create(&newTutor).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTutorCourses(token string) ([]*model.Course, error) {
	tutor, err := engine.FetchTutorByToken(token)
	if err != nil {
		return nil, err
	}
	var courses []models.Course
	// err = utils.DB.Joins("enrolls", utils.DB.Where(&models.Enroll{StudentID: student.ID})).Find(&courses).Error
	err = utils.DB.Model(&models.Course{}).Joins("INNER JOIN course_assigns ON course_assigns.tutor_id = ? AND courses.id = course_assigns.course_id", tutor.ID).Scan(&courses).Error

	if err != nil {
		return nil, err
	}
	var gqlCourses []*model.Course
	for _, course := range courses {
		if course.DeletedAt == nil {
			gqlCourses = append(gqlCourses, course.CreateToGraphData())
		}
	}
	return gqlCourses, nil
}

func CheckIFTutorTeaches(courses []*model.Course, course *models.Course) bool {
	for _, c := range courses {
		if c.ID == course.ID.String() {
			return true
		}
	}
	return false
}
