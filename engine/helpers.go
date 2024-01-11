package engine

import (
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
)

var (
	ErrWrongCreds  = errors.New("invalid email or password")
	ErrWrongToken  = errors.New("invalid token")
	BadgeThreshold = [17]int{500, 1000, 1500, 2000, 3000, 4000, 5000, 6000, 8000, 10000, 15000, 20000, 25000, 30000, 40000, 50000, 100000}
	BadgePhoto     = [17]string{
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
		"https://png.pngtree.com/png-clipart/20190604/original/pngtree-badge-png-image_996483.jpg",
	}
)

func GenerateToken(authid uuid.UUID, name string) (*model.AuthPayload, error) {
	authToken, err := utils.GenerateJWTForAuthId(authid)
	if err != nil {
		return nil, err
	}
	return &model.AuthPayload{
		Token: authToken,
		Name:  name,
	}, nil
}

func FetchStudentByID(userId string) (*models.Student, error) {
	var student models.Student
	err := utils.DB.Preload(clause.Associations).Where("id = ?", userId).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func FetchStudentByAuthToken(jwt string) (*models.Student, error) {
	studentId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, ErrWrongToken
	}
	return FetchStudentByID(studentId)
}

func FetchStudent(cond *models.Student) (*models.Student, error) {
	var student models.Student
	err := utils.DB.Preload(clause.Associations).Where(cond).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}
func FetchStudents(page int) ([]models.Student, error) {
	var students []models.Student
	err := utils.DB.Preload(clause.Associations).Limit(50).Offset(50 * (page - 1)).Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func FetchAdminByID(id string) (*models.Admin, error) {
	var admin models.Admin
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
func FetchAdminByToken(jwt string) (*models.Admin, error) {
	adminId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, ErrWrongToken
	}
	return FetchAdminByID(adminId)
}
func FetchAdmin(cond *models.Admin) (*models.Admin, error) {
	var admin models.Admin
	err := utils.DB.Where(cond).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FetchPatnerByID(id string) (*models.Patner, error) {
	var patner models.Patner
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&patner).Error
	if err != nil {
		return nil, err
	}
	return &patner, nil
}
func FetchPatner(cond *models.Patner) (*models.Patner, error) {
	var patner models.Patner
	err := utils.DB.Where(cond).First(&patner).Error
	if err != nil {
		return nil, err
	}
	return &patner, nil
}
func FetchPatnerByToken(jwt string) (*models.Patner, error) {
	patnerId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, ErrWrongToken
	}
	return FetchPatnerByID(patnerId)
}

func FetchTutorByID(id string) (*models.Tutor, error) {
	var tutor models.Tutor
	err := utils.DB.Where("id = ?", id).First(&tutor).Error
	if err != nil {
		return nil, err
	}
	return &tutor, nil
}
func FetchTutorByToken(jwt string) (*models.Tutor, error) {
	tutorId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, ErrWrongToken
	}
	return FetchTutorByID(tutorId)
}
func FetchTutor(cond *models.Tutor) (*models.Tutor, error) {
	var tutor models.Tutor
	err := utils.DB.Where(cond).First(&tutor).Error
	if err != nil {
		return nil, err
	}
	return &tutor, nil
}

func DeleteCourseByID(courseId string) (bool, error) {
	err := utils.DB.Model(&models.Course{}).Where("id = ?", courseId).Update("deleted_at", time.Now()).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func FetchStudentsDoingCourseById(courseId string) ([]models.Student, error) {
	var enrollments []models.Enroll

	err := utils.DB.Where("course_id=?", courseId).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	var studentIds []string
	for _, enroll := range enrollments {
		studentIds = append(studentIds, enroll.StudentID.String())
	}

	var students []models.Student
	err = utils.DB.Where("id in (?)", studentIds).Find(&students).Error
	if err != nil {
		return nil, err
	}

	return students, nil
}

func UpdateStudentGHToken(userId, ghToken string) (bool, error) {
	err := utils.DB.Model(&models.Student{}).Where("id = ?", userId).Update("github_token", ghToken).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
func DeleteFileByUrl(fileUrl string) (bool, error) {
	err := utils.DB.Where("url = ?", fileUrl).Delete(&models.File{}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
func FetchEnrolls(cond *models.Enroll) ([]models.Enroll, error) {
	var enrolls []models.Enroll
	err := utils.DB.Where(cond).Find(&enrolls).Error
	if err != nil {
		return nil, err
	}
	return enrolls, nil
}
func FetchBadges(cond *models.Badge) ([]models.Badge, error) {
	var badges []models.Badge
	err := utils.DB.Where(cond).Find(&badges).Error
	if err != nil {
		return nil, err
	}
	return badges, nil
}

func FetchEnroll(cond *models.Enroll) (*models.Enroll, error) {
	var enroll models.Enroll
	err := utils.DB.Where(cond).First(&enroll).Error
	if err != nil {
		return nil, err
	}
	return &enroll, nil
}

func FetchPMeeting(id string) (*models.PersonalMeetings, error) {
	var pm models.PersonalMeetings
	err := utils.DB.Where("id = ?", id).First(&pm).Error
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func FetchCourseAssign(cond *models.CourseAssign) (*models.CourseAssign, error) {
	var assigned models.CourseAssign
	err := utils.DB.Where(cond).First(&assigned).Error
	if err != nil {
		return nil, err
	}
	return &assigned, nil
}

func FetchAllCourses(ref string) ([]models.Course, error) {
	var courses []models.Course
	err := utils.DB.Where(&models.Course{SchoolRef: ref}).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
func FetchAllCoursesByCond(cond *models.Course) ([]models.Course, error) {
	var courses []models.Course
	err := utils.DB.Where(cond).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func FetchAllEmailz(cond *models.Emailz) ([]models.Emailz, error) {
	var emails []models.Emailz
	err := utils.DB.Where(cond).Find(&emails).Error
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func FetchEmailById(id string) (*models.Emailz, error) {
	var email models.Emailz
	err := utils.DB.Where("id = ?", id).First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}

func FetchCourse(cond *models.Course) (*models.Course, error) {
	var course models.Course
	err := utils.DB.Preload(clause.Associations).Where(cond).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func FetchChapterByID(id string) (*models.Chapter, error) {
	var chapter models.Chapter
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func FetchCourseById(id string) (*models.Course, error) {
	var course models.Course
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func FetchChaptersByCourseID(id string) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := utils.DB.Where("course_id = ?", id).Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	return chapters, nil
}

func FetchSectionsByChapterID(id string) ([]models.Section, error) {
	var sections []models.Section
	err := utils.DB.Where("chapter_id = ?", id).Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func FetchDBMeetings(courseIds []string) ([]models.Meeting, error) {
	log.Println(courseIds)
	var meetings []models.Meeting
	err := utils.DB.Where("course_id in (?)", courseIds).Order("time desc").Find(&meetings).Error
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func FetchQuizes(courseId string) ([]models.Quiz, error) {
	cond := models.Quiz{CourseID: courseId}
	var quizes []models.Quiz
	err := utils.DB.Preload(clause.Associations).Where(&cond).Find(&quizes).Error
	if err != nil {
		return nil, err
	}
	return quizes, nil
}

func FetchAssignments(courseId uuid.UUID) ([]models.Assignment, error) {
	cond := models.Assignment{CourseID: courseId}
	var assignments []models.Assignment
	err := utils.DB.Preload(clause.Associations).Where(&cond).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

func FetchSubmissions(assignmentId uuid.UUID) ([]models.Submission, error) {
	cond := models.Submission{AssignmentID: assignmentId}
	var submissions []models.Submission
	err := utils.DB.Preload(clause.Associations).Where(&cond).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func FetchSubmissionById(id string) (*models.Submission, error) {
	var submission models.Submission
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func FetchAssignmentById(id string) (*models.Assignment, error) {
	var assignments models.Assignment
	err := utils.DB.Preload(clause.Associations).Where("id = ?", id).First(&assignments).Error
	if err != nil {
		return nil, err
	}
	return &assignments, nil
}

func FetchQuizById(id string) (*models.Quiz, error) {
	var quiz models.Quiz
	err := utils.DB.Preload(clause.Associations).Where("id = ?", uuid.FromStringOrNil(id)).First(&quiz).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func GetQuizResult(studentId, quizId uuid.UUID) (*models.QuizComplete, error) {
	result := models.QuizComplete{
		StudentID: studentId,
		QuizID:    quizId,
	}
	err := utils.DB.Preload(clause.Associations).Where(&result).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func FetchUserByEmail(em string) error {
	_, err := FetchStudent(&models.Student{Email: em})
	if err != nil {
		_, err := FetchTutor(&models.Tutor{Email: em})
		if err != nil {
			_, err := FetchAdmin(&models.Admin{Email: em})
			if err != nil {
				_, err := FetchPatner(&models.Patner{Email: em})
				if err != nil {
					return errors.New("invalid email address")
				}
				return nil
			}
			return nil
		}
		return nil
	}
	return nil
}

func FetchIdAndRoleByToken(token string) (uuid.UUID, model.Role, error) {
	student, err := FetchStudentByAuthToken(token)
	if err != nil {
		tutor, err := FetchTutorByToken(token)
		if err != nil {
			admin, err := FetchAdminByToken(token)
			if err != nil {
				partner, err := FetchPatnerByToken(token)
				if err != nil {
					return uuid.Nil, "", errors.New("invalid token")
				}
				return partner.ID, model.RolePatner, nil
			}
			return admin.ID, model.RoleAdmin, nil
		}
		return tutor.ID, model.RoleTutor, nil
	}
	return student.ID, model.RoleStudent, nil
}

func FetchUserAndModelByToken(token string) (*models.UserInfo, interface{}, error) {
	userId, role, err := FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, nil, err
	}
	if !role.IsValid() {
		return nil, nil, errors.New("invalid role")
	}
	var userInfo models.UserInfo
	userModel := models.RoleModels[role]

	err = utils.DB.Model(&userModel).Where("id = ?", userId).First(&userInfo).Error
	if err != nil {
		return nil, nil, err
	}
	return &userInfo, userModel, nil
}

func FetchSchoolById(id string) (*models.Admin, error) {
	var admin models.Admin
	err := utils.DB.Where("school_id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, err
}

func FetchSchoolByRef(ref string) (*models.Admin, error) {
	var admin models.Admin
	err := utils.DB.Where("school_ref = ?", ref).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, err
}

func FetchSchoolAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := utils.DB.Find(&admins).Error
	if err != nil {
		return nil, err
	}
	return admins, nil
}
