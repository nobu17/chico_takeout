package mail

import (
	"chico/takeout/common"
	"chico/takeout/usecase/order"
	"fmt"

	"chico/takeout/infrastructures/memory"
	"chico/takeout/infrastructures/smtp"
)

const (
	SendGrid = "SendGrid"
	Smtp     = "Smtp"
	Console  = "Console"
)

func NewSendOrderMailService(cfg common.MailConfig) order.SendOrderMailService {
	switch cfg.Mailer {
	case SendGrid:
		fmt.Println("use send grid mailer.")
		return NewSendGridSendOrderMail()
	case Smtp:
		fmt.Println("use smtp mailer.")
		return smtp.NewSmtpSendOrderMail()
	}
	fmt.Println("use memory mailer.(use for test.)")
	return memory.NewMemorySendOrderMail()
}
