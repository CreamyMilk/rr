package quiz

import (
	"aytp/engine"
	"aytp/engine/tutors"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	"fmt"

	"strings"

	uuid "github.com/satori/go.uuid"
)

func CreateQuiz(input model.CreateQuizInput) (bool, error) {
	courses, err := tutors.GetTutorCourses(input.Token)
	if err != nil {
		return false, err
	}

	course, err := engine.FetchCourseById(input.CourseID)
	if err != nil {
		return false, err
	}

	if !tutors.CheckIFTutorTeaches(courses, course) {
		return false, errors.New("you are not assigned to teach this course")
	}

	newQuiz := models.Quiz{
		Title:    input.Title,
		CourseID: input.CourseID,
	}
	err = utils.DB.Create(&newQuiz).Error
	if err != nil {
		return false, err
	}

	for i, q := range input.QuestionsInput {
		answerFound := make(map[string]bool)
		for _, answer := range q.Answer {
			for _, choice := range q.Choices {
				if answer == choice {
					answerFound[answer] = true
				}
			}
		}

		for _, answerCheck := range answerFound {
			if !answerCheck {
				return false, fmt.Errorf("an answer for question %d is missing in choices", i+1)
			}
		}

		newQuestion := &models.Question{
			QuizID:  newQuiz.ID,
			Text:    q.Text,
			Answer:  strings.Join(q.Answer, ","),
			Choices: strings.Join(q.Choices, ","),
		}

		err = utils.DB.Create(&newQuestion).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func FetchQuiz(token, courseId string) ([]*model.Quiz, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	quizes, err := engine.FetchQuizes(courseId)
	if err != nil {
		return nil, err
	}
	gqlQuizes := make([]*model.Quiz, len(quizes))
	for i, quiz := range quizes {
		gqlQuizes[i] = quiz.CreateToGraphData()
	}
	return gqlQuizes, nil
}

func FetchOneQuiz(token, quizId string) (*model.Quiz, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	quiz, err := engine.FetchQuizById(quizId)
	if err != nil {
		return nil, err
	}
	return quiz.CreateToGraphData(), nil
}

func MarkQuiz(input model.SubmitQuizInput) (bool, error) {
	student, err := engine.FetchStudentByAuthToken(input.Token)
	if err != nil {
		return false, err
	}

	quiz, err := engine.FetchQuizById(input.QuizID)
	if err != nil {
		return false, err
	}

	var compl models.QuizComplete
	err = utils.DB.Where(&models.QuizComplete{QuizID: quiz.ID, StudentID: student.ID}).First(&compl).Error
	if err == nil {
		return false, errors.New("quiz already taken")
	}

	marks := 0
	for i, q := range quiz.Questions {
		answers := strings.Split(q.Answer, ",")
		for _, answer := range answers {
			if answer == input.Answers[i] {
				marks += 1
				break
			}
		}
	}
	newCompletion := models.QuizComplete{
		QuizID:    quiz.ID,
		Marks:     marks,
		Total:     len(quiz.Questions),
		Answers:   strings.Join(input.Answers, ","),
		StudentID: student.ID,
	}

	err = utils.DB.Create(newCompletion).Error
	if err != nil {
		return false, err
	}

	var current_points = student.Points
	err = utils.DB.Model(&models.Student{}).Where("id = ?", student.ID).Update("points", current_points+marks).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func FetchQuizResult(token, quizId string) (*model.QuizResult, error) {
	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return nil, err
	}
	quiz, err := engine.FetchQuizById(quizId)
	if err != nil {
		return nil, err
	}

	result, err := engine.GetQuizResult(student.ID, uuid.FromStringOrNil(quizId))
	if err != nil {
		return nil, err
	}

	return &model.QuizResult{
		Quiz:  quiz.CreateToGraphData(),
		Marks: &result.Marks,
		RealAnswers: func(questions []*models.Question) []string {
			ans := make([]string, len(questions))
			for i, q := range questions {
				ans[i] = q.Answer
			}
			return ans
		}(quiz.Questions),
		Total:   &result.Total,
		Answers: strings.Split(result.Answers, ","),
	}, nil
}

func FetchQuizResults(token, courseId string) ([]*model.QuizResult, error) {
	student, err := engine.FetchStudentByAuthToken(token)
	if err != nil {
		return nil, err
	}
	quizes, err := engine.FetchQuizes(courseId)
	if err != nil {
		return nil, err
	}

	gqlRes := make([]*model.QuizResult, len(quizes))

	for i, quiz := range quizes {
		var resModel model.QuizResult
		resModel.Quiz = quiz.CreateToGraphData()

		result, err := engine.GetQuizResult(student.ID, uuid.FromStringOrNil(quiz.ID.String()))
		if err == nil {
			resModel.Marks = &result.Marks
			resModel.Total = &result.Total
		}
		gqlRes[i] = &resModel
	}
	return gqlRes, nil
}

func FetchStudentQuizResults(token, courseId string) ([]*model.StudentQuizResult, error) {
	_, err := engine.FetchTutorByToken(token)
	if err != nil {
		return nil, err
	}

	quizes, err := engine.FetchQuizes(courseId)
	if err != nil {
		return nil, err
	}

	enrolls, err := engine.FetchEnrolls(&models.Enroll{CourseID: uuid.FromStringOrNil(courseId)})
	if err != nil {
		return nil, err
	}

	var quizCompletes []models.QuizComplete
	var averageMark int
	var gqlResults = make([]*model.StudentQuizResult, len(quizes))

	for i, quiz := range quizes {
		err := utils.DB.Where("quiz_id = ?", quiz.ID).Find(&quizCompletes).Error
		if err != nil {
			return nil, err
		}
		for _, quizComplete := range quizCompletes {
			averageMark += quizComplete.Marks
		}

		if len(quizCompletes) > 0 {
			averageMark = averageMark / len(quizCompletes)
		}

		gqlResults[i] = &model.StudentQuizResult{
			Quiz:          quiz.CreateToGraphData(),
			Attempts:      len(quizCompletes),
			TotalEnrolled: len(enrolls),
			AverageMark:   averageMark,
		}
	}
	return gqlResults, nil
}

func FetchResult(quizID string, courseID string) ([]*model.Result, error) {
	query := `
		SELECT e.student_id, s.name, COALESCE(qc.marks, 0) AS marks, COALESCE(qc.total, 0) AS totalmarks, 
		CASE WHEN qc.student_id IS NULL THEN 'UNDONE' ELSE 'DONE' END AS quiz_status
		FROM enrolls e
		LEFT JOIN students s ON e.student_id = s.id
		LEFT JOIN quiz_completes qc ON e.student_id = qc.student_id AND qc.quiz_id = ?
		WHERE e.course_id = ?
	`
	rows, err := utils.DB.Raw(query, quizID, courseID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		var studentID, name, quizStatus string
		var marks, totalMarks int
		rows.Scan(&studentID, &name, &marks, &totalMarks, &quizStatus)

		status := model.StudentQuizStatusUndone
		if quizStatus == "DONE" {
			status = model.StudentQuizStatusDone
		}

		result := &model.Result{
			StudentID:   studentID,
			StudentName: name,
			Marks:       marks,
			Total:       totalMarks,
			Done:        status,
			QuizID:      quizID,
		}
		results = append(results, result)
	}
	return results, nil
}

// func FetchResult(quizID string, courseID string) ([]*model.Result, error) {
// 	var complete []models.QuizComplete
// 	err := utils.DB.Where("quiz_id = ?", quizID).Find(&complete).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	var notDone []models.Enroll
// 	rows1, err := utils.DB.Raw("SELECT course_id,student_id from enrolls WHERE enrolls.course_id=? AND enrolls.student_id not in (SELECT student_id FROM quiz_completes WHERE quiz_completes.quiz_id=?)", courseID, quizID).Rows()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows1.Next() {
// 		var cid, sid string
// 		rows1.Scan(&cid, &sid)
// 		notDone = append(notDone, models.Enroll{CourseID: uuid.FromStringOrNil(cid), StudentID: uuid.FromStringOrNil(sid)})
// 	}

// 	var hasDone []models.Enroll
// 	rows2, err := utils.DB.Raw("SELECT course_id,student_id from enrolls WHERE enrolls.course_id=? AND enrolls.student_id in (SELECT student_id FROM quiz_completes WHERE quiz_completes.quiz_id=?)", courseID, quizID).Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for rows2.Next() {
// 		var cid, sid string
// 		rows2.Scan(&cid, &sid)
// 		hasDone = append(hasDone, models.Enroll{CourseID: uuid.FromStringOrNil(cid), StudentID: uuid.FromStringOrNil(sid)})
// 	}

// SELECT e.student_id, COALESCE(qc.marks, 0) As marks FROM enrolls e LEFT JOIN ( SELECT student_id, SUM(marks) AS marks FROM quiz_completes GROUP BY student_id) qc ON e.student_id=qc.student_id WHERE e.course_id = 'ca1a2498-4bf0-4a94-94da-1cfb64619f42';

// 	var stIDs []string
// 	for _, not := range hasDone {
// 		stIDs = append(stIDs, not.StudentID.String())
// 	}
// 	for _, has := range notDone {
// 		stIDs = append(stIDs, has.StudentID.String())
// 	}
// 	var students []models.Student
// 	err = utils.DB.Where("id in (?)", stIDs).Find(&students).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	var studentlist = make(map[string]models.Student)
// 	var hasDoneQuiz = make(map[string]models.QuizComplete)

// 	for _, comp := range complete {
// 		hasDoneQuiz[comp.StudentID.String()] = comp
// 	}

// 	for _, s := range students {
// 		studentlist[s.ID.String()] = s
// 	}

// 	var results []*model.Result

// 	for k, v := range studentlist {
// 		_, found := hasDoneQuiz[k]
// 		results = append(results, &model.Result{
// 			StudentName: v.Name,
// 			StudentID:   k,
// 			Marks:       hasDoneQuiz[k].Marks,
// 			Done: func(b bool) model.StudentQuizStatus {
// 				if b {
// 					return model.StudentQuizStatusDone
// 				}
// 				return model.StudentQuizStatusUndone
// 			}(found),
// 			Total:  hasDoneQuiz[k].Total,
// 			QuizID: hasDoneQuiz[k].QuizID.String(),
// 		})
// 	}

// 	return results, nil
// }
