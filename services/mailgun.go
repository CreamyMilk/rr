package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/mailgun/mailgun-go"
)

type EmailData struct {
	Subject     string
	ContentData interface{}
	EmailTo     string
	EmailFrom   string
	Template    string
}

type EmailResponse struct {
	MessageID string
	Response  string
}

// this function parses variables into the respective templates
func processTemplate(templatePath string, dynamicData interface{}) (string, error) {
	var temp bytes.Buffer
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	tmpl.Execute(&temp, dynamicData)
	return temp.String(), nil
}

func SendSomeEmail(emailAddress []string, sommessage string) error {
	templatePath := "services/emailtemplates/basic.html"
	fromEmailAdress := "login-alert@getabacus.app"

	dynamicData := map[string]interface{}{
		"message": sommessage,
	}

	body, err := processTemplate(templatePath, dynamicData)
	if err != nil {
		return err
	}

	for _, emailAddress := range emailAddress {
		data := EmailData{
			Subject:     "🛎️ LMS Email",
			ContentData: body,
			EmailTo:     emailAddress,
			Template:    templatePath,
			EmailFrom:   fromEmailAdress,
		}
		_, err = sendEmail(data)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func sendEmail(object EmailData) (*EmailResponse, error) {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	var err error
	mg := mailgun.NewMailgun(domain, apiKey)

	sender := object.EmailFrom
	from := fmt.Sprintf("AYTP <%s>", sender)
	subject := object.Subject
	recipient := object.EmailTo
	body := fmt.Sprintf("%v", object.ContentData)

	message := mg.NewMessage(from, subject, body, recipient)
	message.SetHtml(body)

	response, messageId, err := mg.Send(message)
	if err != nil {
		log.Printf("Error sending email")
		return nil, err
	}

	return &EmailResponse{
		MessageID: messageId,
		Response:  response,
	}, nil
}
