package mail

import (
	"errors"
	"net/smtp"
	"os"
	"strings"
)

type loginAuth struct {
	username, password string
}

func auth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func Send(to []string, subject string, message string) error {
	//configure
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")

	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	msg := []byte("Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"From: Inspire Bets <" + from + ">\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message + "\r\n")

	// Create authentication
	auth := auth(from, password)

	//send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)

	return err
}

// usage:
// auth := LoginAuth("loginname", "password")
// err := smtp.SendMail(smtpServer + ":25", auth, fromAddress, toAddresses, []byte(message))
// or
// client, err := smtp.Dial(smtpServer)
// client.Auth(LoginAuth("loginname", "password"))
