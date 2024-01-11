package students

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
)

func ValidateRegistration(input model.RegisterStudentInput) (*model.AuthPayload, error) {
	if !(input.Password == input.ConfirmPassword) {
		return nil, errors.New("passwords do not match")
	}

	err := engine.FetchUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	password, err := utils.HashString(input.Password)
	if err != nil {
		return nil, errors.New("internal encryption error")
	}

	adm, err := engine.FetchSchoolByRef(input.School)
	if err != nil {
		return nil, err
	}

	newStudent := models.Student{
		Name:       input.Name,
		Email:      input.Email,
		SchoolID:   adm.SchoolID,
		SchoolName: adm.SchoolName,
		Status:     model.StudentStatusInactive,
		Password:   password,
	}

	err = utils.DB.Create(&newStudent).Error
	if err != nil {
		return nil, err
	}

	return engine.GenerateToken(newStudent.ID, newStudent.Name)
}

func Login(input model.LoginInput) (*model.AuthPayload, error) {
	student, err := engine.FetchStudent(&models.Student{Email: input.Email})
	if err != nil {
		return nil, engine.ErrWrongCreds
	}

	if !utils.CompareHashedString(student.Password, input.Password) {
		return nil, engine.ErrWrongCreds
	}
	return engine.GenerateToken(student.ID, student.Name)
}

func FetchAllStudents(token string, page int) ([]*model.Student, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err == nil {
		_, err := engine.FetchStudentByAuthToken(token)
		if err == nil {
			return nil, errors.New("invalid token")
		}
	}
	students, err := engine.FetchStudents(page)
	if err != nil {
		return nil, err
	}
	var gqlStudents = make([]*model.Student, len(students))
	for i, student := range students {
		gqlStudents[i] = student.CreateToGraphData()
	}
	return gqlStudents, nil
}

func GetStudentCourses(token string) ([]*model.Course, error) {
	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return nil, err
	}
	var courses []models.Course
	// err = utils.DB.Joins("enrolls", utils.DB.Where(&models.Enroll{StudentID: student.ID})).Find(&courses).Error
	err = utils.DB.Model(&models.Course{}).Joins("INNER JOIN enrolls ON enrolls.student_id = ? AND courses.id = enrolls.course_id", student.ID).Scan(&courses).Error

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

func GetCourseStudents(courseID string) ([]*model.Student, error) {
	var students []models.Student

	err := utils.DB.Raw("SELECT students.* FROM students JOIN enrolls ON students.id = enrolls.student_id WHERE enrolls.course_id = ?", courseID).Scan(&students).Error
	if err != nil {
		return nil, err
	}

	var result []*model.Student
	for _, student := range students {
		result = append(result, student.CreateToGraphData())
	}
	return result, nil
}

func FetchStudentsNumber(schoolref string) (int, error) {
	admin, err := engine.FetchAdmin(&models.Admin{
		SchoolRef: schoolref,
	})
	if err != nil {
		return 0, err
	}
	var students []models.Student

	err = utils.DB.Where("school_name=?", admin.SchoolName).Find(&students).Error
	if err != nil {
		return 0, err
	}
	
	return len(students),nil
}

func GetCourseNames(courseID string) (*model.Course, error) {
	var course models.Course
	err := utils.DB.Where("id = ?", courseID).First(&course).Error
	if err != nil {
		return nil, err
	}
	return course.CreateToGraphData(), nil
}

func ChangePassword(input model.ChangePasswordInput) (bool, error) {
	if input.ConfirmNewPassword != input.NewPassword {
		return false, errors.New("confirm password does not match")
	}
	if input.OldPassword == input.NewPassword {
		return false, errors.New("please update to a different password")
	}
	user, userModel, err := engine.FetchUserAndModelByToken(input.Token)
	if err != nil {
		return false, err
	}
	if !utils.CompareHashedString(user.Password, input.OldPassword) {
		return false, errors.New("wrong old password")
	}
	hashedNewPass, err := utils.HashString(input.NewPassword)
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	err = utils.DB.Model(&userModel).Where("id = ?", user.ID).Update("password", hashedNewPass).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FetchProfile(token string) (*model.User, error) {
	user, _, err := engine.FetchUserAndModelByToken(token)
	if err != nil {
		return nil, err
	}
	return user.CreateToGraphData(), nil
}

func FetchPoints(token string) (int, error) {
	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return 0, err
	}
	go CheckBagdes(student)
	return student.Points, nil
}

func CheckBagdes(student *models.Student) bool {
	badges, err := engine.FetchBadges(&models.Badge{StudentID: student.ID})
	if err != nil {
		return false
	}

	for i, badgeNum := range engine.BadgeThreshold {
		if student.Points >= badgeNum {
			badgeFound := false
			for _, badge := range badges {
				if badge.PointsIndex == badgeNum {
					badgeFound = true
					break
				}
			}
			if badgeFound {
				continue
			} else if !badgeFound {
				newBadge := models.Badge{
					StudentID:   student.ID,
					PointsIndex: badgeNum,
					Redeemed:    false,
					PhotoUrl:    engine.BadgePhoto[i],
				}
				err := utils.DB.Create(&newBadge).Error
				if err != nil {
					return false
				}
			}
		} else {
			break
		}

	}
	return true
}

func FetchBadges(token string) ([]*model.Badge, error) {
	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return nil, err
	}
	badges, err := engine.FetchBadges(&models.Badge{StudentID: student.ID})
	if err != nil {
		return nil, err
	}
	gqlBadges := make([]*model.Badge, len(badges))
	for i, b := range badges {
		gqlBadges[i] = &model.Badge{
			ID:          b.ID.String(),
			PhotoURL:    b.PhotoUrl,
			PointsIndex: b.PointsIndex,
			Redeemed:    b.Redeemed,
		}
	}
	return gqlBadges, nil
}

func RedeemMyBadge(token, badgeId string) (bool, error) {
	_, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return false, err
	}
	err = utils.DB.Model(&models.Badge{}).Where("id = ?", badgeId).Update("redeemed", true).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
