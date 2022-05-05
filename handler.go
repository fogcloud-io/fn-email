package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

var (
	SMTP_HOST     = os.Getenv("SMTP_HOST")
	SMTP_PORT     = os.Getenv("SMTP_PORT")
	SMTP_ACCOUNT  = os.Getenv("SMTP_ACCOUNT")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	SMTP_TO       = os.Getenv("SMTP_TO")

	Subject = "fogcloud.io cloud-function"
)

// Handle a serverless request
func Handle(req []byte) string {
	err := sendEmail(formatEmailBody(Subject, req))
	if err != nil {
		return err.Error()
	}
	return ""
}

func Handler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sendEmail(formatEmailBody(Subject, reqBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func sendEmail(msg []byte) error {
	auth := smtp.PlainAuth("", SMTP_ACCOUNT, SMTP_PASSWORD, SMTP_HOST)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", SMTP_HOST, SMTP_PORT), auth, SMTP_ACCOUNT, []string{SMTP_TO}, msg)
	if err != nil {
		log.Printf("sendEmail: %s", err)
	}
	return err
}

func formatEmailBody(subject string, msg []byte) []byte {
	return []byte(fmt.Sprintf(`To: %s
Subject: %s

%s
`, SMTP_TO, subject, string(msg)))
}
