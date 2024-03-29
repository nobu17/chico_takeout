package smtp

import (
	"chico/takeout/usecase/order"
)

type SmtpSendOrderMail struct {
	mailer *smtpMail
}

func NewSmtpSendOrderMail() *SmtpSendOrderMail {
	return &SmtpSendOrderMail{}
}

func (s *SmtpSendOrderMail) SendComplete(data order.OrderCompleteMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}

func (s *SmtpSendOrderMail) SendCancel(data order.OrderCancelMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}

func (s *SmtpSendOrderMail) SendDailySummary(data order.ReservationSummaryMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}