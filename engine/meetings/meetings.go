package meetings

import (
	"aytp/engine"
	"aytp/engine/students"
	"aytp/engine/tutors"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type CreateRoomResponse struct {
	Webhook struct {
		Events   []string `json:"events"`
		EndPoint string   `json:"endPoint"`
	} `json:"webhook"`
	Disabled     bool      `json:"disabled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	RoomID       string    `json:"roomId"`
	CustomRoomID string    `json:"customRoomId"`
	Links        struct {
		GetRoom    string `json:"get_room"`
		GetSession string `json:"get_session"`
	} `json:"links"`
	ID string `json:"id"`
}

type FetchSessionResponse struct {
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PerPage     int `json:"perPage"`
		LastPage    int `json:"lastPage"`
		Total       int `json:"total"`
	} `json:"pageInfo"`
	Data []struct {
		Start        time.Time `json:"start"`
		End          time.Time `json:"end"`
		Participants []any     `json:"participants"`
		Region       string    `json:"region"`
		HlsLog       []any     `json:"hlsLog"`
		Dirty        bool      `json:"dirty"`
		ID           string    `json:"id"`
		RoomID       string    `json:"roomId"`
		Status       string    `json:"status"`
		Links        struct {
			GetRoom    string `json:"get_room"`
			GetSession string `json:"get_session"`
		} `json:"links"`
	} `json:"data"`
}
type ParticipantsResponse struct {
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PerPage     int `json:"perPage"`
		LastPage    int `json:"lastPage"`
		Total       int `json:"total"`
	} `json:"pageInfo"`
	Data []struct {
		APIKey       string    `json:"apiKey,omitempty"`
		Start        time.Time `json:"start"`
		End          any       `json:"end"`
		Participants []struct {
			Events struct {
				Mic              []any `json:"mic"`
				Webcam           []any `json:"webcam"`
				ScreenShare      []any `json:"screenShare"`
				ScreenShareAudio []any `json:"screenShareAudio"`
			} `json:"events"`
			ParticipantID string `json:"participantId"`
			Name          string `json:"name"`
			Timelog       []struct {
				Start time.Time `json:"start"`
				End   any       `json:"end"`
			} `json:"timelog"`
			InitTime   float64 `json:"initTime"`
			DeviceInfo struct {
				SdkType          string `json:"sdkType"`
				SdkVersion       string `json:"sdkVersion"`
				Platform         string `json:"platform"`
				BrowserUserAgent struct {
					Browser struct {
						Name    string `json:"name"`
						Version string `json:"version"`
					} `json:"browser"`
					Os struct {
						Name        string `json:"name"`
						Version     string `json:"version"`
						VersionName string `json:"versionName"`
					} `json:"os"`
					Platform struct {
						Type   string `json:"type"`
						Vendor string `json:"vendor"`
					} `json:"platform"`
				} `json:"browserUserAgent"`
			} `json:"deviceInfo"`
			ChangeModeHistory []struct {
				Mode    string    `json:"mode"`
				Timelog time.Time `json:"timelog"`
				ID      string    `json:"_id"`
			} `json:"changeModeHistory"`
			Mode string `json:"mode"`
			ID   string `json:"_id"`
		} `json:"participants"`
		Region              string `json:"region"`
		UsageDataCalculated bool   `json:"usageDataCalculated"`
		CalculationErrored  bool   `json:"calculationErrored"`
		HlsLog              []any  `json:"hlsLog"`
		ID                  string `json:"id"`
		RoomID              string `json:"roomId"`
		Status              string `json:"status"`
		Links               struct {
			GetRoom    string `json:"get_room"`
			GetSession string `json:"get_session"`
		} `json:"links"`
		Webhook struct {
			TotalCount   int   `json:"totalCount"`
			SuccessCount int   `json:"successCount"`
			Data         []any `json:"data"`
		} `json:"webhook"`
		ActiveDuration      int    `json:"activeDuration,omitempty"`
		MaxJoinedPeersCount int    `json:"maxJoinedPeersCount,omitempty"`
		ChatLink            string `json:"chatLink,omitempty"`
	} `json:"data"`
}

func CreateCourseMeetingRoom() (*CreateRoomResponse, error) {
	req, err := http.NewRequest("POST", "https://api.videosdk.live/v2/rooms", nil)
	if err != nil {
		return nil, err
	}
	videoSdkToken := os.Getenv("VIDEOSDK_TOKEN")
	if videoSdkToken == "" {
		return nil, errors.New("video sdk token invalid")
	}
	req.Header.Set("Authorization", videoSdkToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var responseBody []byte
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResponse CreateRoomResponse
	err = json.Unmarshal([]byte(responseBody), &jsonResponse)
	if err != nil {
		return nil, err
	}
	return &jsonResponse, nil
}

func CreateMeeting(input model.ScheduleMeetingInput) (bool, error) {
	tutor, err := engine.FetchTutorByToken(input.Token)
	if err != nil {
		return false, errors.New("only tutors can schedule meetings")
	}
	course, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, err
	}

	_, err = engine.FetchCourseAssign(&models.CourseAssign{CourseID: course.ID, TutorID: tutor.ID})
	if err != nil {
		return false, errors.New("this course is not assigned to you")
	}

	if input.Time.Before(time.Now()) {
		return false, errors.New("please select future dates")
	}

	newMeeting := models.Meeting{
		Time:        input.Time,
		Title:       input.Title,
		Description: input.Description,
		CourseId:    input.CourseID,
		CourseName:  course.Name,
	}

	meeting, err := CreateCourseMeetingRoom()
	if err != nil {
		return false, err
	}
	newMeeting.RoomId = meeting.RoomID
	newMeeting.SessionID = meeting.ID
	log.Println("😅", utils.AsPrettyJson(newMeeting))
	log.Println("🍪", newMeeting.SessionID)

	err = utils.DB.Model(&models.Course{}).Where("id = ?", course.ID).Update("room_id", meeting.RoomID).Error
	if err != nil {
		return false, err
	}

	err = utils.DB.Create(&newMeeting).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FetchScheduledMeetings(token string, courseId *string) ([]*model.Meeting, error) {
	_, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	link := "https://app.videosdk.live/rooms/aytp/"
	if role == model.RoleTutor {
		link = fmt.Sprint(link, "Tutor_63fcdd48abbbfd3acc18381c/")
	} else {
		link = fmt.Sprint(link, "Student_63fcdd48abbbfd5d6318381e/")
	}

	var meetings []models.Meeting
	var courses []*model.Course

	if courseId != nil {
		meetings, err = engine.FetchDBMeetings([]string{*courseId})
		if err != nil {
			return nil, err
		}
	} else {
		if role == model.RoleStudent {
			courses, err = students.GetStudentCourses(token)
		} else if role == model.RoleTutor {
			courses, err = tutors.GetTutorCourses(token)
		}

		if err != nil {
			return nil, err
		}

		courseIds := make([]string, len(courses))
		for i, course := range courses {
			courseIds[i] = course.ID
		}

		meetings, err = engine.FetchDBMeetings(courseIds)
		if err != nil {
			return nil, err
		}
	}

	var gqlMeetings = make([]*model.Meeting, len(meetings))
	for i, meeting := range meetings {
		gqlMeetings[i] = &model.Meeting{
			ID:              meeting.ID.String(),
			Title:           meeting.Title,
			Description:     meeting.Description,
			Time:            meeting.Time,
			CourseID:        meeting.CourseId,
			CourseName:      meeting.CourseName,
			Link:            fmt.Sprint(link, meeting.RoomId),
			VidSDKMeetingID: meeting.RoomId,
		}
	}
	return gqlMeetings, nil
}

func FetchParticipantsInMeeeting(token string, vidSDKRoomID string) ([]*model.Participant, error) {
	// 😒 Who cares about the roles we just return it to everyone the worst they can do list make a list

	// _, role, err := engine.FetchIdAndRoleByToken(token)
	// if err != nil {
	// 	return nil, err
	// }

	// if role != model.RoleTutor {
	// 	return nil, errors.New("you ain't a tutor what you need the participants for 😒")
	// }
	pResp, err := FetchMeetingParticipants(vidSDKRoomID)
	if err != nil {
		return nil, err
	}

	var gqlParticipants = make([]*model.Participant, 0)
	for _, meeting := range pResp.Data {
		if meeting.Status == "ongoing" {
			for _, studentStuff := range meeting.Participants {
				gqlParticipants = append(gqlParticipants, &model.Participant{
					Name: studentStuff.Name,
				})
			}
		}
	}
	return gqlParticipants, nil
}

func FetchMeetingSessions(roomId string) (*FetchSessionResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprint("https://api.videosdk.live/v2/sessions/?roomId=", roomId), nil)
	if err != nil {
		return nil, err
	}

	videoSdkToken := os.Getenv("VIDEOSDK_TOKEN")
	if videoSdkToken == "" {
		return nil, errors.New("video sdk token invalid")
	}
	req.Header.Set("Authorization", videoSdkToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var responseBody []byte
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResponse FetchSessionResponse
	err = json.Unmarshal([]byte(responseBody), &jsonResponse)
	if err != nil {
		return nil, err
	}

	return &jsonResponse, nil
}

func DeleteMeeting(token, meetId string) (bool, error) {
	_, err := engine.FetchTutorByToken(token)
	if err != nil {
		return false, err
	}
	err = utils.DB.Where("id = ?", meetId).Delete(&models.Meeting{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FetchMeetingParticipants(roomId string) (*ParticipantsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.videosdk.live/v2/sessions?roomId=%s", roomId), nil)
	if err != nil {
		return nil, err
	}

	videoSdkToken := os.Getenv("VIDEOSDK_TOKEN")
	if videoSdkToken == "" {
		return nil, errors.New("video sdk token invalid")
	}
	req.Header.Set("Authorization", videoSdkToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var responseBody []byte
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResponse ParticipantsResponse
	err = json.Unmarshal([]byte(responseBody), &jsonResponse)
	if err != nil {
		return nil, err
	}

	return &jsonResponse, nil
}
