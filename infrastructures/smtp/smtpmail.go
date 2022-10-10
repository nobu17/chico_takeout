package smtp

import (
	"net/smtp"

	"chico/takeout/common"
)

type smtpMail struct {
}

func (s *smtpMail) sendMail(subject, message, from, bcc string, to []string) error {
	cfg := common.GetConfig().Mail
	auth := smtp.PlainAuth(
		"",
		cfg.User,
		cfg.Pass,
		cfg.Host,
	)

	return smtp.SendMail(
		cfg.Host+":"+cfg.Port,
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
