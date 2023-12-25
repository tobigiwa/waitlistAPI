package http

import (
	"Blockride-waitlistAPI/env"
	"Blockride-waitlistAPI/templ"
	"bytes"
	"fmt"
	"hash/fnv"
	"html/template"

	gomail "gopkg.in/mail.v2"
)

type templData struct {
	Name  string
	Key   string
	email string
}

const (
	from = "tobigiwa@zohomail.com"

	// password = "eBgH2mNqztkZ"
	smtpHost = "smtp.zoho.com"
	smtpPort = 465
)

func sendConfirmationMail(name, email, key string) error {

	var (
		tpl bytes.Buffer
		s   = templData{Name: name, Key: key, email: email}
	)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", s.email)
	m.SetHeader("Subject", "BlockRide Confirmation mail")

	prepareMail(s, &tpl)
	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(smtpHost, smtpPort, from, env.GetEnvVar().EmailPswd)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func prepareMail(data templData, mailBuf *bytes.Buffer) error {

	t, err := template.ParseFS(templ.EmailHTML, "mail.html")
	if err != nil {
		return err
	}

	if err := t.Execute(mailBuf, data); err != nil {
		return err
	}

	return nil
}

// generateRedisKey returns a 10 digit string hash.
func generateRedisKey(m string) string {
	h := fnv.New32a()
	h.Write([]byte(m))
	return fmt.Sprintf("%d", h.Sum32())
}
