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

// func SendEmail(email, subject, url string) error {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", os.Getenv("SMTP_EMAIL_SENDER"))
// 	m.SetHeader("To", email)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", `<h1>Email Confirmation</h1>
//                     <h2>Hello `+email+`</h2>
//                     <p>Thank you for joining us. Please confirm your email by clicking on the following link</p>
//                     <a href='`+url+`'> Click here</a>
// 					atau masuk ke link `+url)

// 	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_EMAIL_SENDER"), os.Getenv("SMTP_EMAIL_PASS"))

// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}

// 	return nil
// }

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

	// m.SetBody("text/html", `<div style="font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; background-color: #F5F6FA; border-radius: 20px; width: 100%; margin: auto;">
	// 							<div class="card"
	// 							style="text-align: center; border-radius: 10px; overflow: hidden; width: 100%; max-width: 400px; background-color: white; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); margin: auto;">
	// 							<div class="card-header" style="background-color: #2395FF; padding: 20px; color: white;">
	// 								<h1>Email Verification</h1>
	// 							</div>
	// 							<div class="card-body" style="padding: 20px;">
	// 								<h2 style="margin-bottom: 16px;">Dear Mr./Ms. `+name+`,</h2>
	// 								<p style="margin-bottom: 40px;">Thank you for registering with us. Please click the button below to verify your
	// 								email address and activate your account.</p>
	// 								<a href='`+url+`' class="button"
	// 								style="background-color: #2395FF; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none;">Verify
	// 								Email</a>
	// 								<p style=" margin-top: 40px;">If you did not register for an account, please ignore this email.</p>
	// 							</div>
	// 							<div class="card-footer"
	// 								style="padding: 20px; background-color: #f5f5f5; text-align: center; font-size: 14px; color: #666;">
	// 								@2024 `+os.Getenv("APP_NAME")+`
	// 							</div>
	// 							</div>
	// 						</div>`)

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

	// m.SetBody("text/html", `<div style="font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; background-color: #F5F6FA; border-radius: 20px; width: 100%; margin: auto;">
	// 							<div class="card"
	// 							style="text-align: center; border-radius: 10px; overflow: hidden; width: 100%; max-width: 400px; background-color: white; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); margin: auto;">
	// 							<div class="card-header" style="background-color: #2395FF; padding: 20px; color: white;">
	// 								<h1>Reset Password</h1>
	// 							</div>
	// 							<div class="card-body" style="padding: 20px;">
	// 								<h2 style="margin-bottom: 16px;">Dear Mr./Ms. `+name+`,</h2>
	// 								<p style="margin-bottom: 40px;">Follow this link to reset your password for your account.</p>
	// 								<a href='`+url+`' class="button"
	// 								style="background-color: #2395FF; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none;">Reset Password</a>
	// 								<p style=" margin-top: 40px;">If you didn't ask to reset your password, please ignore this email.</p>
	// 							</div>
	// 							<div class="card-footer"
	// 								style="padding: 20px; background-color: #f5f5f5; text-align: center; font-size: 14px; color: #666;">
	// 								@2024 `+os.Getenv("APP_NAME")+`
	// 							</div>
	// 							</div>
	// 						</div>`)

	d := gomail.NewDialer("smtp.gmail.com", 587, provider, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
