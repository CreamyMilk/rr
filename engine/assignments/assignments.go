package assignments

import (
	"aytp/engine"
	"aytp/engine/courses"

	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

func CreateAssignment(input model.CreateAssignmentInput) (bool, error) {
	if input.DueDate.Before(time.Now()) {
		return false, errors.New("please select future dates")
	}
	course, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, err
	}

	if !courses.CheckEnrollMent(input.Token, input.CourseID) {
		return false, errors.New("you are not assigned to this course")
	}

	newAssignment := models.Assignment{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		CourseID:    course.ID,
	}

	err = utils.DB.Create(&newAssignment).Error
	if err != nil {
		return false, err
	}

	for _, gqlFile := range input.Files {
		newFile := &models.AssignmentFile{
			AssignmentID: newAssignment.ID,
			Name:         gqlFile.Name,
			Url:          gqlFile.URL,
		}
		err = utils.DB.Create(&newFile).Error
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func FetchAllAssignments(token, courseId string) ([]*model.Assignment, error) {
	assignments, err := engine.FetchAssignments(uuid.FromStringOrNil(courseId))
	if err != nil {
		return nil, err
	}
	if !courses.CheckEnrollMent(token, courseId) {
		return nil, errors.New("you are not enrolled or assigned to this course")
	}

	id, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	gqlAssigns := make([]*model.Assignment, len(assignments))
	for i, assign := range assignments {
		gqlAssignment := assign.CreateToGraphData()
		
		if role == model.RoleStudent {
			var submission models.Submission
			err := utils.DB.Model(&models.Submission{}).Where(&models.Submission{AssignmentID: assign.ID, StudentID: id}).First(&submission).Error
			if err == nil {
				gqlAssignment.Marks = &submission.Marks
				gqlAssignment.Remark = &submission.Remark
			}
		}
		gqlAssigns[i] = gqlAssignment
	}
	return gqlAssigns, nil
}

func FetchOneAssignments(token, assignmentId string) (*model.Assignment, error) {
	assignment, err := engine.FetchAssignmentById(assignmentId)
	if err != nil {
		return nil, err
	}
	if !courses.CheckEnrollMent(token, assignment.CourseID.String()) {
		return nil, errors.New("you are not enrolled or assigned to this course")
	}
	gqlAssignment := assignment.CreateToGraphData()

	id, role, err := engine.FetchIdAndRoleByToken(token)
	if err != nil {
		return nil, err
	}

	if role == model.RoleStudent {
		var submission models.Submission
		err := utils.DB.Model(&models.Submission{}).Where(&models.Submission{AssignmentID: assignment.ID, StudentID: id}).First(&submission).Error
		if err == nil {
			gqlAssignment.Marks = &submission.Marks
			gqlAssignment.Remark = &submission.Remark
		}
	}
	return gqlAssignment, nil
}

func SubmitAssignment(input model.SubmitAssignmentInput) (bool, error) {
	assignment, err := engine.FetchAssignmentById(input.AssignmentID)
	if err != nil {
		return false, err
	}
	student, err := engine.FetchStudentByAuthToken(input.Token)
	if err != nil {
		return false, err
	}
	if !courses.CheckEnrollMent(input.Token, assignment.CourseID.String()) {
		return false, errors.New("you are not enrolled in this course")
	}

	err = utils.DB.Model(&models.Submission{}).Where(&models.Submission{StudentID: student.ID, AssignmentID: assignment.ID}).First(nil).Error
	if err == nil {
		return false, errors.New("you have already submitted this assignment")
	}

	newSubmission := models.Submission{
		CourseID:       assignment.CourseID,
		AssignmentID:   assignment.ID,
		StudentName:    student.Name,
		SubmissionDate: time.Now(),
		StudentID:      student.ID,
		Note:           *input.Note,
	}

	err = utils.DB.Create(&newSubmission).Error
	if err != nil {
		return false, err
	}

	for _, gqlFile := range input.Files {
		newFile := &models.SubmissionFile{
			SubmissionID: newSubmission.ID,
			Name:         gqlFile.Name,
			Url:          gqlFile.URL,
		}
		err = utils.DB.Create(&newFile).Error
		if err != nil {
			return false, err
		}

	}
	return true, nil
}

func FetchAllSubmissions(token, assignmentId string) ([]*model.Submission, error) {
	assignment, err := engine.FetchAssignmentById(assignmentId)
	if err != nil {
		return nil, err
	}

	if !courses.CheckEnrollMent(token, assignment.CourseID.String()) {
		return nil, errors.New("you are not enrolled or assigned to this course")
	}

	submissions, err := engine.FetchSubmissions(uuid.FromStringOrNil(assignmentId))
	if err != nil {
		return nil, err
	}

	gqlSubs := make([]*model.Submission, len(submissions))
	for i, sub := range submissions {
		gqlSub := sub.CreateToGraphData()
		gqlSub.Assignment = assignment.CreateToGraphData()
		gqlSubs[i] = gqlSub
	}
	return gqlSubs, nil
}

func FetchOneSubmission(token, submissionId string) (*model.Submission, error) {
	submission, err := engine.FetchSubmissionById(submissionId)
	if err != nil {
		return nil, err
	}

	assignment, err := engine.FetchAssignmentById(submission.AssignmentID.String())
	if err != nil {
		return nil, err
	}

	if !courses.CheckEnrollMent(token, assignment.CourseID.String()) {
		return nil, errors.New("you are not enroleed or assigned to this course")
	}

	gqlSub := submission.CreateToGraphData()
	gqlSub.Assignment = assignment.CreateToGraphData()
	return gqlSub, nil
}

func MarkStudentAssignment(input model.MarkAssignentInput) (bool, error) {
	_, err := engine.FetchTutorByToken(input.Token)
	if err != nil {
		return false, err
	}
	submission, err := engine.FetchSubmissionById(input.SubmissionID)
	if err != nil {
		return false, err
	}
	assignment, err := engine.FetchAssignmentById(submission.AssignmentID.String())
	if err != nil {
		return false, err
	}
	if !courses.CheckEnrollMent(input.Token, assignment.CourseID.String()) {
		return false, errors.New("you are not enrolled or assigned to this course")
	}
	err = utils.DB.Model(&models.Submission{}).Where("id = ?", input.SubmissionID).Updates(&models.Submission{Marks: input.Marks, Remark: *input.Remark}).Error
	if err != nil {
		return false, err
	}
	
	student, err := engine.FetchStudentByID(submission.StudentID.String())
	if err != nil {
		return false, err
	}

	err = utils.DB.Model(&models.Student{}).Where("id = ?", student.ID).Update("points", student.Points + input.Marks).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
