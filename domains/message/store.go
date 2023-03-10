package message

import (
	"strings"
	"time"

	"chico/takeout/common"
)

const (
	TopMessageId    = "1"
	MyPageMessageId = "2"
)

type MessageRepository interface {
	Find(id string) (*StoreMessage, error)
	Create(message *StoreMessage) (string, error)
	Update(message *StoreMessage) error
}

type StoreMessage struct {
	id      string
	content Content
	edited  time.Time
	created time.Time
}

func (s *StoreMessage) GetId() string {
	return s.id
}

func (s *StoreMessage) GetContent() string {
	return s.content.GetValue()
}

func (s *StoreMessage) GetEdited() time.Time {
	return s.edited
}

func (s *StoreMessage) GetCreated() time.Time {
	return s.created
}

func NewMessage(id, content string) (*StoreMessage, error) {
	switch id {
	case TopMessageId:
		return NewStoreTopMessage(content)
	case MyPageMessageId:
		return NewMyPageMessage(content)
	}
	return nil, common.NewValidationError("Id", "not supported type message")
}

func NewMessageForOrm(id, content string, edited, created time.Time) (*StoreMessage, error) {
	msg := StoreMessage{id: id, edited: edited, created: created}
	err := msg.Set(content)
	return &msg, err
}

func NewStoreTopMessage(content string) (*StoreMessage, error) {
	return newStoreMessage(TopMessageId, content)
}

func NewMyPageMessage(content string) (*StoreMessage, error) {
	return newStoreMessage(MyPageMessageId, content)
}

func newStoreMessage(id, content string) (*StoreMessage, error) {
	contentVal, err := NewContent(content)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) == "" {
		return nil, common.NewValidationError(id, "required")
	}
	now := *common.GetNowTime()
	msg := &StoreMessage{
		id:      id,
		content: *contentVal,
		edited:  now,
		created: now,
	}
	return msg, nil
}

func (s *StoreMessage) Set(content string) error {
	contentVal, err := NewContent(content)
	if err != nil {
		return err
	}
	now := *common.GetNowTime()

	s.content = *contentVal
	s.edited = now

	return nil
}
