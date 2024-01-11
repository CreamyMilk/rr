package emails

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/services"
	"aytp/utils"
	"errors"
	"log"
)

func CreateEmail(input model.EmailSendingInput) (bool, error) {
	userInfo, _, err := engine.FetchUserAndModelByToken(input.Token)
	if err != nil {
		return false, err
	}
	err = engine.FetchUserByEmail(input.To)
	if err != nil {
		return false, err
	}

	newEmail := models.Emailz{
		From:    userInfo.Email,
		To:      input.To,
		Subject: input.Subject,
		Body:    input.Body,
	}

	go func(recipeints []string, message string) {
		err = services.SendSomeEmail(recipeints, message)
		if err != nil {
			log.Println("📧 has failed us", err)
		}
	}([]string{userInfo.Email}, input.Body)

	err = utils.DB.Create(&newEmail).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func FetchOneEmail(token, emailId string) (*model.Email, error) {
	_, err := utils.ValidateJWTForAuthId(token)
	if err != nil {
		return nil, err
	}
	email, err := engine.FetchEmailById(emailId)
	if err != nil {
		return nil, err
	}
	return email.CreateToGraphData(), nil
}

func FetchManyEmails(token string, emailType model.EmailType) ([]*model.Email, error) {
	userInfo, _, err := engine.FetchUserAndModelByToken(token)
	if err != nil {
		return nil, err
	}

	var condition models.Emailz
	if !emailType.IsValid() {
		return nil, errors.New("invalid email type")
	}

	if emailType == model.EmailTypeInbox {
		condition.To = userInfo.Email
	} else if emailType == model.EmailTypeOutbox {
		condition.From = userInfo.Email
	}

	emails, err := engine.FetchAllEmailz(&condition)
	if err != nil {
		return nil, err
	}
	var gqlCourses = make([]*model.Email, len(emails))
	for i, email := range emails {
		gqlCourses[i] = email.CreateToGraphData()
	}
	return gqlCourses, nil
}
