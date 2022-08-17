package memory

import (
	"fmt"
	"strings"

	"chico/takeout/usecase/order"
)

type MemorySendOrderMail struct {
}

func NewMemorySendOrderMail() *MemorySendOrderMail {
	return &MemorySendOrderMail{}
}

func (m *MemorySendOrderMail) SendComplete(data order.OrderCompleteMailData) error {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("from:%s\n", data.SendFrom))

	toStr := ""
	for _, to := range data.SendTo {
		toStr += to +","
	}
	b.WriteString(fmt.Sprintf("to:%s\n", toStr))

	b.WriteString(fmt.Sprintf("bcc:%s\n", data.Bcc))
	b.WriteString(fmt.Sprintf("title:%s\n", data.Title))
	b.WriteString(fmt.Sprintf("message:%s\n", data.Message))

	fmt.Println(b.String())

	return nil
}

func (m *MemorySendOrderMail) SendCancel(data order.OrderCancelMailData) error {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("from:%s\n", data.SendFrom))

	toStr := ""
	for _, to := range data.SendTo {
		toStr += to +","
	}
	b.WriteString(fmt.Sprintf("to:%s\n", toStr))

	b.WriteString(fmt.Sprintf("bcc:%s\n", data.Bcc))
	b.WriteString(fmt.Sprintf("title:%s\n", data.Title))
	b.WriteString(fmt.Sprintf("message:%s\n", data.Message))

	fmt.Println(b.String())

	return nil
}
