package courses

import (
	"aytp/engine"
	"aytp/engine/meetings"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func GetAllCourses(token, schoolref string) ([]*model.Course, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err != nil {
		return nil, engine.ErrWrongToken
	}
	courses, err := engine.FetchAllCourses(schoolref)
	if err != nil {
		return nil, err
	}

	var gqlCourses = make([]*model.Course, len(courses))
	for i, course := range courses {
		var course_assign models.CourseAssign
		tutorsName := ""

		err = utils.DB.Where(&models.CourseAssign{CourseID: course.ID}).First(&course_assign).Error
		if err == nil {
			var tutor models.Tutor
			err = utils.DB.Where("id = ?", course_assign.TutorID).First(&tutor).Error
			if err == nil {
				tutorsName = tutor.Name
			}
		}

		gqlCourses[i] = course.CreateToGraphData()
		gqlCourses[i].Tutor = tutorsName
	}
	return gqlCourses, nil
}

func GetAllPublicCourses(schoolref string) ([]*model.Course, error) {
	courses, err := engine.FetchAllCoursesByCond(&models.Course{
		SchoolRef:  schoolref,
		CourseType: model.CourseTypePublic,
	})
	if err != nil {
		return nil, err
	}
	var gqlCourses = make([]*model.Course, len(courses))
	for i, course := range courses {
		gqlCourses[i] = course.CreateToGraphData()
	}
	return gqlCourses, nil
}

func GetCourseByID(token, courseId string) (*model.Course, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	course, err := engine.FetchCourseById(courseId)
	if err != nil {
		return nil, err
	}
	return course.CreateToGraphData(), nil
}

func AddCourse(input model.NewCourse) (bool, error) {
	tutor, err := engine.FetchTutorByToken(input.Token)
	if err != nil {
		return false, err
	}

	_, err = engine.FetchSchoolByRef(input.Schoolref)
	if err != nil {
		return false, err
	}

	_, err = engine.FetchCourse(&models.Course{Name: input.Name})
	if err == nil {
		return false, errors.New("course already exists")
	}
	meetingData, err := meetings.CreateCourseMeetingRoom()
	if err != nil {
		return false, errors.New("failed to create room")
	}

	newCourse := models.Course{
		Name:             input.Name,
		Icon:             input.Icon,
		Description:      input.Description,
		Duration:         input.Duration,
		SchoolRef:        input.Schoolref,
		RoomId:           meetingData.RoomID,
		DefaultSessionId: meetingData.ID,
		Cost:             input.Cost,
		CourseLevel:      input.Courselevel,
		CourseType:       input.Coursetype,
	}

	err = utils.DB.Create(&newCourse).Error
	if err != nil {
		return false, errors.New("could not add new course")
	}

	newCourseAssign := models.CourseAssign{
		TutorID:  tutor.ID,
		CourseID: newCourse.ID,
	}

	err = utils.DB.Create(&newCourseAssign).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func EnrollCourse(token, courseId string) (bool, error) {
	if CheckEnrollMent(token, courseId) {
		return false, errors.New("course already applied")
	}

	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return false, err
	}

	course, err := engine.FetchCourseById(courseId)
	if err != nil {
		return false, err
	}

	newEnroll := models.Enroll{
		StudentID: student.ID,
		CourseID:  course.ID,
	}

	err = utils.DB.Model(&models.Course{}).Where("id = ?", course.ID).Update("enrolled", course.Enrolled+1).Error
	if err != nil {
		return false, err
	}

	err = utils.DB.Create(&newEnroll).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckEnrollMent(token, courseId string) bool {
	user, userModel, err := engine.FetchUserAndModelByToken(token)
	if err != nil {
		return false
	}
	_, isStudent := userModel.(models.Student)
	_, isTutor := userModel.(models.Tutor)

	if isStudent {
		_, err := engine.FetchEnroll(&models.Enroll{StudentID: user.ID, CourseID: uuid.FromStringOrNil(courseId)})
		return err == nil
	} else if isTutor {
		_, err := engine.FetchCourseAssign(&models.CourseAssign{TutorID: user.ID, CourseID: uuid.FromStringOrNil(courseId)})
		return err == nil
	}
	return false
}
