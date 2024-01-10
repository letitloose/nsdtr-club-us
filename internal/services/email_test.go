package services

import (
	"fmt"
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	email := Email{
		Username: "louis.garwood@gmail.com",
		Password: "qhej mkki eexf esgn",
		Host:     "smtp.gmail.com"}

	subject := "Test email from Go!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><h1>Hello World!</h1><p>%s<p></body></html>", time.Now())

	err := email.SendEmail(subject, mime, body)
	if err != nil {
		t.Fatal("failed: ", err)
	}
}
