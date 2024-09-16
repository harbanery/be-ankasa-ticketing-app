package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	AppName string
	Name    string
	URL     string
}

func SendEmailVerification(email, name, token string, userID int) error {
	appName := os.Getenv("APP_NAME")
	appUrl := os.Getenv("APP_URL")
	provider := os.Getenv("SMTP_EMAIL_SENDER")
	password := os.Getenv("SMTP_EMAIL_PASS")

	if provider == "" || password == "" {
		return fmt.Errorf("SMTP credentials empty")
	}

	if appName == "" {
		appName = "Ankasa"
	}

	uuid := strconv.Itoa(userID)
	if uuid == "" {
		return fmt.Errorf("user ID empty")
	}
	url := appUrl + "auth/email-verification?id=" + uuid + "&token=" + token

	templatePath := filepath.Join("src", "templates", "email_verification.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse html file from templates: %v", err)
	}

	data := EmailData{
		AppName: appName,
		Name:    name,
		URL:     url,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", provider)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification - "+appName)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, provider, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendRequestResetPassword(email, name, token string, userID int) error {
	appName := os.Getenv("APP_NAME")
	appUrl := os.Getenv("APP_URL")
	provider := os.Getenv("SMTP_EMAIL_SENDER")
	password := os.Getenv("SMTP_EMAIL_PASS")

	if provider == "" || password == "" {
		return fmt.Errorf("SMTP credentials empty")
	}

	if appName == "" {
		appName = "Ankasa"
	}

	uuid := strconv.Itoa(userID)
	if uuid == "" {
		return fmt.Errorf("user ID empty")
	}
	url := appUrl + "auth/reset-password?id=" + uuid + "&token=" + token

	templatePath := filepath.Join("src", "templates", "reset_password.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse html file from templates: %v", err)
	}

	data := EmailData{
		AppName: appName,
		Name:    name,
		URL:     url,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", provider)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password - "+appName)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, provider, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
