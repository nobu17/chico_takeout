package mail

import (
	"chico/takeout/usecase/order"
)

type SendGridSendOrderMail struct {
	mailer *sendGridMail
}

func NewSendGridSendOrderMail() *SendGridSendOrderMail {
	return &SendGridSendOrderMail{}
}

func (s *SendGridSendOrderMail) SendComplete(data order.OrderCompleteMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}

func (s *SendGridSendOrderMail) SendCancel(data order.OrderCancelMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}

func (s *SendGridSendOrderMail) SendDailySummary(data order.ReservationSummaryMailData) error {
	return s.mailer.sendMail(data.Title, data.Message, data.SendFrom, data.Cc, data.SendTo)
}