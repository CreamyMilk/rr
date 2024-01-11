package meetings

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 19-20
const (
	defLink        = "https://app.videosdk.live/rooms/aytp/"
	startTimeLimit = 8
	endTimeLimit   = 20
)

func CreatePersonalMeeting(input model.PersonalMeetingInput) (bool, error) {
	student, err := engine.FetchStudentByAuthToken(input.Token)
	if err != nil {
		return false, err
	}
	course, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, err
	}
	courseAssign, err := engine.FetchCourseAssign(&models.CourseAssign{CourseID: uuid.FromStringOrNil(input.CourseID)})
	if err != nil {
		return false, err
	}
	tutor, err := engine.FetchTutorByID(courseAssign.TutorID.String())
	if err != nil {
		return false, err
	}
	if input.Start.Year() != time.Now().Year() || input.Start.Month() != time.Now().Month() || input.Start.Day() != time.Now().Day() {
		return false, errors.New("invalid date")
	}
	fmt.Println(input.Start.Hour(), time.Now().Hour())
	if input.Start.Hour() <= time.Now().Hour() {
		return false, errors.New("select a later time")
	}
	if input.Start.Hour() < startTimeLimit || input.Start.Hour() > endTimeLimit-1 {
		return false, errors.New("time is out of scheduled time")
	}
	if input.Start.Minute() != 0 {
		return false, errors.New("personal meetings are only allocated by hours")
	}

	err = utils.DB.Model(&models.PersonalMeetings{}).
		Where(&models.PersonalMeetings{Year: input.Start.Year(), Month: input.Start.Month(), Day: input.Start.Day(), Hour: input.Start.Hour()}).First(nil).Error
	if err == nil {
		return false, errors.New("select a different time")
	}
	newMeeting := models.PersonalMeetings{
		TutorId:     tutor.ID,
		TutorName:   tutor.Name,
		Tutorcourse: course.Name,
		CourseId:    course.ID,
		Year:        input.Start.Year(), Month: input.Start.Month(), Day: input.Start.Day(), Hour: input.Start.Hour(),
		UserId:   student.ID,
		UserName: student.Name,
		Status:   model.PersonalMeetingStatusRequested,
	}
	err = utils.DB.Create(&newMeeting).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func ChangeMeetingStatus(input model.ChangeMeetingStatus) (bool, error) {
	tutor, err := engine.FetchTutorByToken(input.Token)
	if err != nil {
		return false, err
	}
	var meeting models.PersonalMeetings
	err = utils.DB.Where("id = ? AND tutor_id = ?", input.ID, tutor.ID).First(&meeting).Error
	if err != nil {
		return false, err
	}

	if meeting.Status != model.PersonalMeetingStatusRequested {
		return false, errors.New("meeting is already " + meeting.Status.String())
	}

	roomId := tutor.RoomId
	if roomId == "" {
		meetingRoom, err := CreateCourseMeetingRoom()
		if err != nil {
			return false, err
		}
		roomId = meetingRoom.ID
		tutor.RoomId = roomId
		err = utils.DB.Save(tutor).Error
		if err != nil {
			return false, err
		}
	}
	meeting.Status = input.NewStatus
	meeting.TutorLink = fmt.Sprintf("%sTutor_63fcdd48abbbfd3acc18381c/%s", defLink, roomId)
	meeting.StudentLink = fmt.Sprintf("%sStudent_63fcdd48abbbfd5d6318381e/%s", defLink, roomId)

	err = utils.DB.Save(meeting).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FetchPMeeting(token string) ([]*model.PersonalMeeting, error) {
	userId, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	var meetings []models.PersonalMeetings
	if role == model.RoleTutor {
		err := utils.DB.Where(&models.PersonalMeetings{TutorId: userId}).Order("hour asc").Find(&meetings).Error
		if err != nil {
			return nil, err
		}
	} else if role == model.RoleStudent {
		err := utils.DB.Where(&models.PersonalMeetings{UserId: userId}).Order("hour asc").Find(&meetings).Error
		if err != nil {
			return nil, err
		}
	}

	var gqlMeetings []*model.PersonalMeeting
	for _, v := range meetings {
		if v.Hour < time.Now().Hour() {
			continue
		}
		gqlMeeting := &model.PersonalMeeting{
			ID:        v.ID.String(),
			Status:    v.Status,
			Tutorname: v.TutorName,
			Tutorid:   v.TutorId.String(),
			UserName:  v.UserName,
			UserID:    v.UserId.String(),
			CreatedAt: v.CreatedAt,
			StartHour: v.Hour,
			EndHour:   v.Hour + 1,
		}
		if role == model.RoleTutor {
			gqlMeeting.Link = v.TutorLink
		} else if role == model.RoleStudent {
			gqlMeeting.Link = v.StudentLink
		}
		gqlMeetings = append(gqlMeetings, gqlMeeting)
	}
	return gqlMeetings, nil
}

func FetchAvailableMeetingTimes(token, courseid string) ([]int, error) {
	_, _, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	var todayMeetings []models.PersonalMeetings
	err = utils.DB.Model(&models.PersonalMeetings{}).
		Where(&models.PersonalMeetings{
			Year:     time.Now().Year(),
			Month:    time.Now().Month(),
			Day:      time.Now().Day(),
			CourseId: uuid.FromStringOrNil(courseid)}).
		Find(&todayMeetings).Error

	if err != nil {
		return nil, err
	}

	var freeHours []int
loop1:
	for i := time.Now().Hour(); i < endTimeLimit; i++ {
		for _, v := range todayMeetings {
			fmt.Println(i, v.Hour)
			if v.Hour == i {
				continue loop1
			}
		}
		freeHours = append(freeHours, i)
	}
	return freeHours, nil
}
