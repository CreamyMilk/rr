package patners

import (
	"aytp/engine"
	"aytp/engine/meetings"
	"aytp/engine/students"
	"aytp/engine/tutors"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
)

func Login(input model.LoginInput) (*model.AuthPayload, error) {
	patner, err := engine.FetchPatner(&models.Patner{Email: input.Email})
	if err != nil {
		return nil, engine.ErrWrongCreds
	}
	if !utils.CompareHashedString(patner.Password, input.Password) {
		return nil, engine.ErrWrongCreds
	}
	return engine.GenerateToken(patner.ID, patner.Name)
}

func FetchSidebarData(token string) (*model.LayoutDetails, error) {
	userInfo, _, err := engine.FetchUserAndModelByToken(token)
	if err != nil {
		return nil, err
	}

	_, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	emails, err := engine.FetchAllEmailz(&models.Emailz{To: userInfo.Email})
	if err != nil {
		return nil, err
	}
	var points int
	var course_len int
	var courses []*model.Course

	if role == model.RoleStudent {
		student, err := engine.FetchStudentByID(userInfo.ID.String())
		if err != nil {
			return nil, err
		}
		points = student.Points
		courses, err = students.GetStudentCourses(token)
		if err != nil {
			return nil, err
		}
		course_len = len(courses)

	} else if role == model.RoleTutor {
		courses, err = tutors.GetTutorCourses(token)
		if err != nil {
			return nil, err
		}
		course_len = len(courses)
	}

	meetings, err := meetings.FetchScheduledMeetings(token, nil)
	if err != nil {
		return nil, err
	}

	return &model.LayoutDetails{
		Emails:   len(emails),
		Points:   points,
		Courses:  course_len,
		Meetings: len(meetings),
	}, nil
}
