package infrastructure

import (
	"github.com/badoux/checkmail"
)

type EmailChecker struct{}

func (EmailChecker) IsValid(email string) bool {
	return checkmail.ValidateFormat(email) != nil
}

func (EmailChecker) IsExists(email string) (exists bool, errCode string, errMeg string) {
	err := checkmail.ValidateHost(email)
	errCode = ""
	errMeg = ""
	exists = false
	if SMTPErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		errCode = SMTPErr.Code()
		errMeg = SMTPErr.Error()
	}
	return
}
