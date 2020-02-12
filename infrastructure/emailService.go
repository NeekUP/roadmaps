package infrastructure

import (
	"bytes"
	"encoding/json"

	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"io/ioutil"
	"net/smtp"

	"text/template"
)

func NewEmailSender(webHost, senderEmail, senderName, smtpHost, smtpPass string, smtpPort int, log core.AppLogger) core.EmailSender {
	return &emailService{
		siteHost:  webHost,
		fromEmail: senderEmail,
		fromName:  senderName,
		host:      smtpHost,
		port:      smtpPort,
		pass:      smtpPass,
		log:       log,
	}
}

type emailService struct {
	siteHost  string
	fromEmail string
	fromName  string
	host      string
	port      int
	pass      string
	log       core.AppLogger
}

type emailTemplateSource struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (service *emailService) Send(recipient string, subject string, body string) (bool, error) {
	auth := smtp.PlainAuth(
		"",
		service.fromEmail,
		service.pass,
		service.host,
	)

	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"\r\n"+
		"%s\r\n", service.fromName, service.fromEmail, recipient, subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", service.host, service.port),
		auth,
		service.fromEmail,
		[]string{recipient},
		[]byte(msg),
	)

	return err == nil, err
}

type userRegistrationEmail struct {
	Host      string
	Recipient string
	UserId    string
	Secret    string
}

func newUserRegistrationEmail(recipient, userId, secret, host string) userRegistrationEmail {
	return userRegistrationEmail{
		Host:      host,
		Recipient: recipient,
		UserId:    userId,
		Secret:    secret,
	}
}

func (service *emailService) Registration(recipient, userId, secret string) (bool, error) {
	service.log.Debugw("Start send email", "from", service.fromEmail, "to", recipient, "type", "Registration")

	ok, err := service.registration(recipient, userId, secret)

	if err == nil {
		service.log.Infow("Send email", "from", service.fromEmail, "to", recipient, "type", "Registration", "status", true)
	} else {
		service.log.Errorw("Send email", "from", service.fromEmail, "to", recipient, "type", "Registration", "status", false, "err", err.Error())
	}
	return ok, err
}

func (service *emailService) registration(recipient, userId, secret string) (bool, error) {
	tpl, err := ioutil.ReadFile("static/emails/registration.tpl")
	if err != nil {
		return false, err
	}

	msg := &emailTemplateSource{}
	err = json.Unmarshal(tpl, &msg)
	if err != nil {
		return false, err
	}

	var out bytes.Buffer
	t := template.Must(template.New("letter").Parse(msg.Body))
	if err := t.Execute(&out, newUserRegistrationEmail(recipient, userId, secret, service.siteHost)); err != nil {
		return false, err
	}
	ok, err := service.Send(recipient, msg.Subject, out.String())

	return ok, err
}
