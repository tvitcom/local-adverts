package util

import(
	"html/template"
	"net/smtp"
	"strings"
	"context"
	"bytes"
	"time"
)

func SendEmailWithTempate(ctx context.Context, smtpHost, smtpPort, fromMail, fromPassword, toMail, tplFile string, mailData map[string]string) error {
	cmdTimeout := time.Second * 30
	to := []string{toMail}
	t, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, mailData); err != nil {
		return err
	}
	auth := smtp.PlainAuth("", fromMail, fromPassword, smtpHost)
	addr := smtpHost + ":" + smtpPort
	_, cancel := context.WithTimeout(ctx, cmdTimeout)
    defer cancel()
	return smtp.SendMail(addr, auth, fromMail, to, buf.Bytes())
}

// import("strings")
func GetEmailsAlias(email string) string {
	splitted := strings.Split(email, "@")
	return splitted[0]
}
