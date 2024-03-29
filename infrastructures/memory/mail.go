package memory

import (
	"fmt"
	"strings"

	"chico/takeout/usecase/order"
)

type MemorySendOrderMail struct {
	Sent []DummyMailData
}

type DummyMailData struct {
	SendTo   []string
	SendFrom string
	Bcc      string
	Title    string
	Message  string
}

func NewMemorySendOrderMail() *MemorySendOrderMail {
	return &MemorySendOrderMail{}
}

func (m *MemorySendOrderMail) SendComplete(data order.OrderCompleteMailData) error {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("from:%s\n", data.SendFrom))

	toStr := ""
	for _, to := range data.SendTo {
		toStr += to + ","
	}
	b.WriteString(fmt.Sprintf("to:%s\n", toStr))

	b.WriteString(fmt.Sprintf("cc:%s\n", data.Cc))
	b.WriteString(fmt.Sprintf("title:%s\n", data.Title))
	b.WriteString(fmt.Sprintf("message:%s\n", data.Message))

	fmt.Println(b.String())

	mData := &DummyMailData{
		Title:    data.Title,
		Message:  data.Message,
		Bcc:      data.Cc,
		SendTo:   data.SendTo,
		SendFrom: data.SendFrom,
	}
	m.Sent = append(m.Sent, *mData)

	return nil
}

func (m *MemorySendOrderMail) SendCancel(data order.OrderCancelMailData) error {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("from:%s\n", data.SendFrom))

	toStr := ""
	for _, to := range data.SendTo {
		toStr += to + ","
	}
	b.WriteString(fmt.Sprintf("to:%s\n", toStr))

	b.WriteString(fmt.Sprintf("cc:%s\n", data.Cc))
	b.WriteString(fmt.Sprintf("title:%s\n", data.Title))
	b.WriteString(fmt.Sprintf("message:%s\n", data.Message))

	fmt.Println(b.String())

	mData := &DummyMailData{
		Title:    data.Title,
		Message:  data.Message,
		Bcc:      data.Cc,
		SendTo:   data.SendTo,
		SendFrom: data.SendFrom,
	}
	m.Sent = append(m.Sent, *mData)

	return nil
}

func (m *MemorySendOrderMail) SendDailySummary(data order.ReservationSummaryMailData) error {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("from:%s\n", data.SendFrom))

	toStr := ""
	for _, to := range data.SendTo {
		toStr += to + ","
	}
	b.WriteString(fmt.Sprintf("to:%s\n", toStr))

	b.WriteString(fmt.Sprintf("cc:%s\n", data.Cc))
	b.WriteString(fmt.Sprintf("title:%s\n", data.Title))
	b.WriteString(fmt.Sprintf("message:%s\n", data.Message))

	fmt.Println(b.String())

	mData := &DummyMailData{
		Title:    data.Title,
		Message:  data.Message,
		Bcc:      data.Cc,
		SendTo:   data.SendTo,
		SendFrom: data.SendFrom,
	}
	m.Sent = append(m.Sent, *mData)

	return nil
}
