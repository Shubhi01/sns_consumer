package smtp

import (
	"fmt"
	"net/smtp"

	"models"

	"github.com/sirupsen/logrus"
)

var logger = logrus.NewEntry(logrus.New())

func Send(smtpObj *models.EmailActionSpec) {
	from := fmt.Sprintf("<%s>", smtpObj.From)
	pass := smtpObj.Password
	to := smtpObj.To

	logger.Info(from)
	logger.Info(pass)

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, smtpObj.Subject, smtpObj.Body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", smtpObj.From, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		logger.Errorf("smtp error: %s", err)
		return
	}

}
