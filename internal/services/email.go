package services

import (
	"fmt"
	"net/smtp"
)

type Email struct {
	Username string
	Password string
	Host     string
}

func (email *Email) SendEmail(subject, mime, body string) error {

	auth := smtp.PlainAuth("", email.Username, email.Password, email.Host)

	if mime == "" {
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}
	subject = fmt.Sprintf("Subject: %s\n", subject)
	message := []byte(subject + mime + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "tollerjones@gmail.com", []string{"louis.garwood@gmail.com"}, message)
	if err != nil {
		return err
	}
	return nil
}
