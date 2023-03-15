package mail

import (
	"fmt"

	"chico/takeout/common"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendGridMail struct {
}

func (s *sendGridMail) sendMail(subject, message, from, cc string, to []string) error {
	cfg := common.GetConfig().Mail

	fromAd := mail.NewEmail("CHICO SPICE管理", from)
	tos := []*mail.Email{}
	for _, t := range to {
		toAd := mail.NewEmail("", t)
		tos = append(tos, toAd)
	}
	p := mail.NewPersonalization()
	p.AddFrom(fromAd)
	p.AddTos(tos...)
	p.Subject = subject
	c := mail.NewContent("text/plain", message)

	if cc != "" {
		ccAd := mail.NewEmail("", cc)
		p.AddCCs(ccAd)
	}

	m := mail.NewV3Mail()
	m.SetFrom(fromAd)
	m.AddPersonalizations(p)
	m.AddContent(c)
	client := sendgrid.NewSendClient(cfg.SendGridKey)
	response, err := client.Send(m)
	if err != nil {
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return nil
	}
}
