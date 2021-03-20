package config

import (
	"gopkg.in/gomail.v2"
	"strings"
)

const MailHost = "smtp.gmail.com"
const MailPort = 587
const ConfigAuthEmail = "email@example.com"
const ConfigAuthPassword = "herepasswordemail"
const MailFromName = "khaiz badaru tammam <khaiz@my.id>"

func SendMail(email, otp string) (string, error) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", MailFromName)
	mailer.SetHeader("To", email, email)
	mailer.SetHeader("Subject", "OTP Login Go Stisla")
	mailer.SetBody("text/html", strings.Join([]string{"Your Otp is <b>", otp, "</b> \n Please don't share this code to anyone"}, ""))

	dialer := gomail.NewDialer(
		MailHost,
		MailPort,
		ConfigAuthEmail,
		ConfigAuthPassword,
	)

	err := dialer.DialAndSend(mailer)

	if err != nil {
		return "unable to send email", err
	}

	return "mail sent!", nil
}
