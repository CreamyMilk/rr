package courses

import (
	"aytp/engine"
	"aytp/engine/students"
	"aytp/engine/tutors"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
)

func DeleteCourse(token, courseID string) (bool, error) {
	if !CheckEnrollMent(token, courseID) {
		return false, errors.New("you are not assigned to this course")
	}
	_, err := engine.DeleteCourseByID(courseID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteFileUpload(token, fileUrl string) (bool, error) {
	_, err := engine.FetchTutorByToken(token)
	if err != nil {
		return false, err
	}
	_, err = engine.DeleteFileByUrl(fileUrl)
	if err != nil {
		return false, err
	}

	return true, nil
}
func AddFilesToCourse(input model.CourseFilesInput) (bool, error) {
	_, err := engine.FetchTutorByToken(input.Token)
	if err != nil {
		return false, err
	}

	course, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, err
	}

	courses, err := tutors.GetTutorCourses(input.Token)
	if err != nil {
		return false, err
	}
	teaches := false
	for _, c := range courses {
		if c.ID == course.ID.String() {
			teaches = true
			break
		}
	}
	if !teaches {
		return false, errors.New("you are not assigned to teach this course")
	}

	for _, file_input := range input.Files {
		newFile := models.File{
			CourseID: course.ID,
			Name:     file_input.Name,
			Url:      file_input.URL,
		}
		err := utils.DB.Create(&newFile).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func GetCourseFiles(token, courseId string) ([]*model.File, error) {
	_, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}
	course, err := engine.FetchCourseById(courseId)
	if err != nil {
		return nil, err
	}

	var user_courses []*model.Course
	related := false

	if role == model.RoleStudent {
		user_courses, err = students.GetStudentCourses(token)
		if err != nil {
			return nil, err
		}
	} else if role == model.RoleTutor {
		user_courses, err = tutors.GetTutorCourses(token)
		if err != nil {
			return nil, err
		}
	}

	for _, c := range user_courses {
		if c.ID == course.ID.String() {
			related = true
			break
		}
	}
	if !related {
		return nil, errors.New("you are not assigned to teach this course")
	}

	var gqlFiles = make([]*model.File, len(course.Files))
	for i, file := range course.Files {
		gqlFiles[i] = file.CreateToGraphData()
	}
	return gqlFiles, nil
}
