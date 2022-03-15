package email

import (
	"encoding/json"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

var configured = false
var auth smtp.Auth
var emailAddress string
var hostAddress string

func IsValidEmailAddress(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}

type emailConfig struct {
	Email struct {
		EmailAddress string `json:"emailAddress"`
		HostAddress  string `json:"hostAddress"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	} `json:"email"`
}

func SendEmail(to, subject, body string) error {
	if !configured {
		configContents, err := os.ReadFile("/studier/server/config.json")
		if err != nil {
			return err
		}

		var config emailConfig
		err = json.Unmarshal(configContents, &config)
		if err != nil {
			return err
		}

		auth = smtp.CRAMMD5Auth(config.Email.Username, config.Email.Password)
		emailAddress = config.Email.EmailAddress
		hostAddress = config.Email.HostAddress

		configured = true
	}

	return smtp.SendMail(hostAddress, auth, to, []string{to}, []byte(strings.ReplaceAll(
		("From: "+emailAddress+"\n"+
			"To: "+to+"\n"+
			"Subject: "+subject+"\n\n"+
			body),
		"\n",
		"\r\n")))
}
