package smtp

import (
	"net/smtp"
	"os"
)

type smtpMail struct {
}

func (s *smtpMail) sendMail(subject, message, from, bcc string, to []string) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASS"),
		os.Getenv("MAIL_HOST"),
	)

	return smtp.SendMail(
		os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"),
		auth,
		from,
		to,
		[]byte(
			"Bcc:"+bcc+"\n"+
				"To:"+to[0]+"\n"+
				"Subject:"+subject+"\r\n"+
				"\r\n"+
				message),
	)
}
