package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"aytp/engine/admins"
	"aytp/engine/assignments"
	"aytp/engine/courses"
	"aytp/engine/emails"
	"aytp/engine/events"
	"aytp/engine/hero"
	"aytp/engine/meetings"
	"aytp/engine/patners"
	"aytp/engine/quiz"
	"aytp/engine/students"
	"aytp/engine/tutors"
	"aytp/graph/generated"
	"aytp/graph/model"
	"context"
	"errors"
	"fmt"
)

// RegisterStudent is the resolver for the RegisterStudent field.
func (r *mutationResolver) RegisterStudent(ctx context.Context, input model.RegisterStudentInput) (*model.AuthPayload, error) {
	return students.ValidateRegistration(input)
}

// Login is the resolver for the Login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	switch input.Role {
	case model.RoleStudent:
		return students.Login(input)
	case model.RoleAdmin:
		return admins.Login(input)
	case model.RolePatner:
		return admins.Login(input)
	case model.RoleTutor:
		return tutors.Login(input)
	default:
		return nil, errors.New("invalid role")
	}
}

// CreateCourse is the resolver for the CreateCourse field.
func (r *mutationResolver) CreateCourse(ctx context.Context, input model.NewCourse) (bool, error) {
	return courses.AddCourse(input)
}

// AddTutor is the resolver for the AddTutor field.
func (r *mutationResolver) AddTutor(ctx context.Context, input model.NewTutor) (bool, error) {
	return tutors.AddNewTutor(input)
}

// AssignTutor is the resolver for the AssignTutor field.
func (r *mutationResolver) AssignTutor(ctx context.Context, input model.AssignCourseToTutorInput) (bool, error) {
	panic(fmt.Errorf("not implemented: AssignTutor - AssignTutor"))
}

// CreateChapter is the resolver for the CreateChapter field.
func (r *mutationResolver) CreateChapter(ctx context.Context, input model.NewChapter) (bool, error) {
	return courses.AddChapter(input)
}

// CreateSection is the resolver for the CreateSection field.
func (r *mutationResolver) CreateSection(ctx context.Context, input model.NewSection) (bool, error) {
	return courses.AddSection(input)
}

// CreateEmail is the resolver for the CreateEmail field.
func (r *mutationResolver) CreateEmail(ctx context.Context, input model.EmailSendingInput) (bool, error) {
	return emails.CreateEmail(input)
}

// CreateRepoForCourse is the resolver for the CreateRepoForCourse field.
func (r *mutationResolver) CreateRepoForCourse(ctx context.Context, input model.RepoCreateionForCourseInput) (bool, error) {
	return courses.CreateRepoForCourse(input)
}

// CreateQuiz is the resolver for the CreateQuiz field.
func (r *mutationResolver) CreateQuiz(ctx context.Context, input model.CreateQuizInput) (bool, error) {
	return quiz.CreateQuiz(input)
}

// SubmitQuiz is the resolver for the SubmitQuiz field.
func (r *mutationResolver) SubmitQuiz(ctx context.Context, input model.SubmitQuizInput) (bool, error) {
	return quiz.MarkQuiz(input)
}

// EnrollCourse is the resolver for the EnrollCourse field.
func (r *mutationResolver) EnrollCourse(ctx context.Context, token string, courseID string) (bool, error) {
	return courses.EnrollCourse(token, courseID)
}

// ChangePassword is the resolver for the ChangePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (bool, error) {
	return students.ChangePassword(input)
}

// ScheduleMeeting is the resolver for the ScheduleMeeting field.
func (r *mutationResolver) ScheduleMeeting(ctx context.Context, input model.ScheduleMeetingInput) (bool, error) {
	return meetings.CreateMeeting(input)
}

// AddFilesToCourse is the resolver for the AddFilesToCourse field.
func (r *mutationResolver) AddFilesToCourse(ctx context.Context, input model.CourseFilesInput) (bool, error) {
	return courses.AddFilesToCourse(input)
}

// DeleteCourse is the resolver for the DeleteCourse field.
func (r *mutationResolver) DeleteCourse(ctx context.Context, token string, courseID string) (bool, error) {
	return courses.DeleteCourse(token, courseID)
}

// DeleteFileUpload is the resolver for the DeleteFileUpload field.
func (r *mutationResolver) DeleteFileUpload(ctx context.Context, token string, fileurl string) (bool, error) {
	return courses.DeleteFileUpload(token, fileurl)
}

// CreateAssignment is the resolver for the CreateAssignment field.
func (r *mutationResolver) CreateAssignment(ctx context.Context, input model.CreateAssignmentInput) (bool, error) {
	return assignments.CreateAssignment(input)
}

// SubmitAssignment is the resolver for the SubmitAssignment field.
func (r *mutationResolver) SubmitAssignment(ctx context.Context, input model.SubmitAssignmentInput) (bool, error) {
	return assignments.SubmitAssignment(input)
}

// MarkAssignent is the resolver for the MarkAssignent field.
func (r *mutationResolver) MarkAssignent(ctx context.Context, input model.MarkAssignentInput) (bool, error) {
	return assignments.MarkStudentAssignment(input)
}

// ChangeLogo is the resolver for the ChangeLogo field.
func (r *mutationResolver) ChangeLogo(ctx context.Context, token string, url string) (bool, error) {
	return admins.ChangeLogo(token, url)
}

// AddSchool is the resolver for the AddSchool field.
func (r *mutationResolver) AddSchool(ctx context.Context, input model.NewSchoolAdmin) (bool, error) {
	return admins.AddRealAdmin(input)
}

// DeleteMeeting is the resolver for the DeleteMeeting field.
func (r *mutationResolver) DeleteMeeting(ctx context.Context, token string, id string) (bool, error) {
	return meetings.DeleteMeeting(token, id)
}

// DeleteEvent is the resolver for the DeleteEvent field.
func (r *mutationResolver) DeleteEvent(ctx context.Context, token string, id string) (bool, error) {
	return events.DeleteEvent(token, id)
}

// AddSchoolAdmin is the resolver for the AddSchoolAdmin field.
func (r *mutationResolver) AddSchoolAdmin(ctx context.Context, input model.NewSchoolAdmin) (bool, error) {
	return admins.AddRealAdmin(input)
}

// ResetPassword is the resolver for the ResetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPasswordInput) (bool, error) {
	return admins.ResetPassword(input)
}

// AddEvents is the resolver for the AddEvents field.
func (r *mutationResolver) AddEvents(ctx context.Context, input model.AddEventsInput) (bool, error) {
	return events.AddEvent(input)
}

// EditEvent is the resolver for the EditEvent field.
func (r *mutationResolver) EditEvent(ctx context.Context, input model.EditEventsInput) (bool, error) {
	return events.EditEvent(input)
}

// EditHero is the resolver for the EditHero field.
func (r *mutationResolver) EditHero(ctx context.Context, input model.EditHeroInput) (bool, error) {
	panic(fmt.Errorf("not implemented: EditHero - EditHero"))
}

// RedeemBadge is the resolver for the RedeemBadge field.
func (r *mutationResolver) RedeemBadge(ctx context.Context, token string, badgeID string) (bool, error) {
	return students.RedeemMyBadge(token, badgeID)
}

// IssueCertificate is the resolver for the IssueCertificate field.
func (r *mutationResolver) IssueCertificate(ctx context.Context, input model.IssueCertificateInput) (*model.Certificate, error) {
	panic(fmt.Errorf("not implemented: IssueCertificate - IssueCertificate"))
}

// AddHero is the resolver for the AddHero field.
func (r *mutationResolver) AddHero(ctx context.Context, input model.AddHeroInput) (bool, error) {
	return hero.AddHero(input)
}

// AddPmeeting is the resolver for the addPmeeting field.
func (r *mutationResolver) AddPmeeting(ctx context.Context, input model.PersonalMeetingInput) (bool, error) {
	return meetings.CreatePersonalMeeting(input)
}

// ChangePMeeting is the resolver for the changePMeeting field.
func (r *mutationResolver) ChangePMeeting(ctx context.Context, input model.ChangeMeetingStatus) (bool, error) {
	return meetings.ChangeMeetingStatus(input)
}

// FetchCourses is the resolver for the FetchCourses field.
func (r *queryResolver) FetchCourses(ctx context.Context, token string, schoolref string) ([]*model.Course, error) {
	return courses.GetAllCourses(token, schoolref)
}

// FetchCourseByID is the resolver for the FetchCourseById field.
func (r *queryResolver) FetchCourseByID(ctx context.Context, token string, id string) (*model.Course, error) {
	return courses.GetCourseByID(token, id)
}

// FetchChaptersByCourseID is the resolver for the FetchChaptersByCourseId field.
func (r *queryResolver) FetchChaptersByCourseID(ctx context.Context, token string, id string) ([]*model.Chapter, error) {
	return courses.GetCourseChapters(token, id)
}

// FetchStudentByID is the resolver for the FetchStudentById field.
func (r *queryResolver) FetchStudentByID(ctx context.Context, id string) (*model.Student, error) {
	panic(fmt.Errorf("not implemented: FetchStudentByID - FetchStudentById"))
}

// FetchStudentByCourseID is the resolver for the FetchStudentByCourseId field.
func (r *queryResolver) FetchStudentByCourseID(ctx context.Context, courseID string) ([]*model.Student, error) {
	return students.GetCourseStudents(courseID)
}

// FetchCourseNameByCourseID is the resolver for the FetchCourseNameByCourseId field.
func (r *queryResolver) FetchCourseNameByCourseID(ctx context.Context, courseID string) (*model.Course, error) {
	return students.GetCourseNames(courseID)
}

// FetchStudents is the resolver for the FetchStudents field.
func (r *queryResolver) FetchStudents(ctx context.Context, token string, page int) ([]*model.Student, error) {
	return students.FetchAllStudents(token, page)
}

// FetchSectionsByChapterID is the resolver for the FetchSectionsByChapterId field.
func (r *queryResolver) FetchSectionsByChapterID(ctx context.Context, token string, id string) ([]*model.Section, error) {
	return courses.GetChapterSections(token, id)
}

// FetchQuiz is the resolver for the FetchQuiz field.
func (r *queryResolver) FetchQuiz(ctx context.Context, token string, courseID string) ([]*model.Quiz, error) {
	return quiz.FetchQuiz(token, courseID)
}

// FetchQuizResults is the resolver for the FetchQuizResults field.
func (r *queryResolver) FetchQuizResults(ctx context.Context, token string, courseID string) ([]*model.QuizResult, error) {
	return quiz.FetchQuizResults(token, courseID)
}

// FetchStudentQuizResults is the resolver for the FetchStudentQuizResults field.
func (r *queryResolver) FetchStudentQuizResults(ctx context.Context, token string, courseID string) ([]*model.StudentQuizResult, error) {
	return quiz.FetchStudentQuizResults(token, courseID)
}

// FetchQuizByID is the resolver for the FetchQuizById field.
func (r *queryResolver) FetchQuizByID(ctx context.Context, token string, quizID string) (*model.Quiz, error) {
	return quiz.FetchOneQuiz(token, quizID)
}

// FetchQuizResultByID is the resolver for the FetchQuizResultById field.
func (r *queryResolver) FetchQuizResultByID(ctx context.Context, token string, quizID string) (*model.QuizResult, error) {
	return quiz.FetchQuizResult(token, quizID)
}

// FetchStudentCourses is the resolver for the FetchStudentCourses field.
func (r *queryResolver) FetchStudentCourses(ctx context.Context, token string) ([]*model.Course, error) {
	return students.GetStudentCourses(token)
}

// FetchProfile is the resolver for the FetchProfile field.
func (r *queryResolver) FetchProfile(ctx context.Context, token string) (*model.User, error) {
	return students.FetchProfile(token)
}

// FetchAllEmails is the resolver for the FetchAllEmails field.
func (r *queryResolver) FetchAllEmails(ctx context.Context, token string, emailType model.EmailType) ([]*model.Email, error) {
	return emails.FetchManyEmails(token, emailType)
}

// FetchEmail is the resolver for the FetchEmail field.
func (r *queryResolver) FetchEmail(ctx context.Context, token string, emailID string) (*model.Email, error) {
	return emails.FetchOneEmail(token, emailID)
}

// FetchStudentCourses is the resolver for the FetchStudentCourses field.
func (r *queryResolver) FetchTutorCourses(ctx context.Context, token string) ([]*model.Course, error) {
	return tutors.GetTutorCourses(token)
}

// FetchMeetings is the resolver for the FetchMeetings field.
func (r *queryResolver) FetchMeetings(ctx context.Context, token string, courseID *string) ([]*model.Meeting, error) {
	return meetings.FetchScheduledMeetings(token, courseID)
}

// FetchParticipantsInMeeeting is the resolver for the FetchParticipantsInMeeeting field.
func (r *queryResolver) FetchParticipantsInMeeeting(ctx context.Context, token string, vidSDKRoomID string) ([]*model.Participant, error) {
	return meetings.FetchParticipantsInMeeeting(token, vidSDKRoomID)
}

// FetchLayoutDetails is the resolver for the FetchLayoutDetails field.
func (r *queryResolver) FetchLayoutDetails(ctx context.Context, token string) (*model.LayoutDetails, error) {
	return patners.FetchSidebarData(token)
}

// FetchCourseFiles is the resolver for the FetchCourseFiles field.
func (r *queryResolver) FetchCourseFiles(ctx context.Context, token string, courseID string) ([]*model.File, error) {
	return courses.GetCourseFiles(token, courseID)
}

// FetchAssignments is the resolver for the FetchAssignments field.
func (r *queryResolver) FetchAssignments(ctx context.Context, token string, courseID string) ([]*model.Assignment, error) {
	return assignments.FetchAllAssignments(token, courseID)
}

// FetchAssignment is the resolver for the FetchAssignment field.
func (r *queryResolver) FetchAssignment(ctx context.Context, token string, assignmentID string) (*model.Assignment, error) {
	return assignments.FetchOneAssignments(token, assignmentID)
}

// FetchSubmissions is the resolver for the FetchSubmissions field.
func (r *queryResolver) FetchSubmissions(ctx context.Context, token string, assignmentID string) ([]*model.Submission, error) {
	return assignments.FetchAllSubmissions(token, assignmentID)
}

// FetchSubmission is the resolver for the FetchSubmission field.
func (r *queryResolver) FetchSubmission(ctx context.Context, token string, submissionID string) (*model.Submission, error) {
	return assignments.FetchOneSubmission(token, submissionID)
}

// FetchLogo is the resolver for the FetchLogo field.
func (r *queryResolver) FetchLogo(ctx context.Context, href string) (string, error) {
	return admins.FetchLogo(href)
}

// FetchSchools is the resolver for the FetchSchools field.
func (r *queryResolver) FetchSchools(ctx context.Context) ([]*model.School, error) {
	return admins.FetchAllSchools()
}

// FetchAllQuizResult is the resolver for the FetchAllQuizResult field.
func (r *queryResolver) FetchAllQuizResult(ctx context.Context, quizID string, courseID string) ([]*model.Result, error) {
	return quiz.FetchResult(quizID, courseID)
}

// FetchPointsByStudentID is the resolver for the FetchPointsByStudentId field.
func (r *queryResolver) FetchPointsByStudentID(ctx context.Context, token string) (int, error) {
	return students.FetchPoints(token)
}

// FetchBadges is the resolver for the FetchBadges field.
func (r *queryResolver) FetchBadges(ctx context.Context, token string) ([]*model.Badge, error) {
	return students.FetchBadges(token)
}

// FetchPublicCourses is the resolver for the FetchPublicCourses field.
func (r *queryResolver) FetchPublicCourses(ctx context.Context, schoolref string) ([]*model.Course, error) {
	return courses.GetAllPublicCourses(schoolref)
}

// FetchSchoolStudentsNumber is the resolver for the FetchSchoolStudentsNumber field.
func (r *queryResolver) FetchSchoolStudentsNumber(ctx context.Context, schoolref string) (int, error) {
	return students.FetchStudentsNumber(schoolref)
}

// FetchEvents is the resolver for the FetchEvents field.
func (r *queryResolver) FetchEvents(ctx context.Context, ref string) ([]*model.Event, error) {
	return events.FetchAllEvents(ref)
}

// FetchEventByID is the resolver for the FetchEventById field.
func (r *queryResolver) FetchEventByID(ctx context.Context, id string) (*model.Event, error) {
	return events.FetchOneEvent(id)
}

// FetchHero is the resolver for the FetchHero field.
func (r *queryResolver) FetchHero(ctx context.Context, ref string) (*model.Hero, error) {
	return hero.FetchHeroByRef(ref)
}

// FetchHeroByID is the resolver for the FetchHeroById field.
func (r *queryResolver) FetchHeroByID(ctx context.Context, id string) (*model.Hero, error) {
	return hero.FetchHeroById(id)
}

// FetchPMeetings is the resolver for the fetchPMeetings field.
func (r *queryResolver) FetchPMeetings(ctx context.Context, token string) ([]*model.PersonalMeeting, error) {
	return meetings.FetchPMeeting(token)
}

// FetchPMeetingSlots is the resolver for the fetchPMeetingSlots field.
func (r *queryResolver) FetchPMeetingSlots(ctx context.Context, token string, courseID string) ([]int, error) {
	return meetings.FetchAvailableMeetingTimes(token, courseID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
