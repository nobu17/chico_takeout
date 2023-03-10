package message

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/message"
	"time"
)

type StoreMessageModel struct {
	Id      string
	Content string
	Edited  time.Time
	Created time.Time
}

func newStoreMessageModel(item *domains.StoreMessage) *StoreMessageModel {
	return &StoreMessageModel{
		Id:      item.GetId(),
		Content: item.GetContent(),
		Edited:  item.GetEdited(),
		Created: item.GetCreated(),
	}
}

type StoreMessageCreateModel struct {
	Id      string
	Content string
}

type StoreMessageUpdateModel struct {
	Id      string
	Content string
}

type StoreMessageUseCase interface {
	Find(id string) (*StoreMessageModel, error)
	Create(model *StoreMessageCreateModel) (string, error)
	CreateInitialMessage() error
	Update(model *StoreMessageUpdateModel) error
}

type storeMessageUseCase struct {
	repository domains.MessageRepository
}

func NewStoreMessageUseCase(repository domains.MessageRepository) StoreMessageUseCase {
	return &storeMessageUseCase{
		repository: repository,
	}
}

// CreateInitialMessage implements MessageUseCase
func (m *storeMessageUseCase) CreateInitialMessage() error {
	topMessage, err := m.repository.Find("1")
	if err != nil {
		return nil
	}
	if topMessage == nil {
		item1, err := domains.NewStoreTopMessage("トップメッセージです。")
		if err != nil {
			return err
		}
		_, err = m.repository.Create(item1)
		if err != nil {
			return err
		}
	}

	myMessage, err := m.repository.Find("2")
	if err != nil {
		return nil
	}
	if myMessage == nil {
		item1, err := domains.NewMyPageMessage("マイページメッセージです。")
		if err != nil {
			return err
		}
		_, err = m.repository.Create(item1)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create implements MessageUseCase
func (m *storeMessageUseCase) Create(model *StoreMessageCreateModel) (string, error) {
	item, err := domains.NewMessage(model.Id, model.Content)
	if err != nil {
		return "", err
	}
	return m.repository.Create(item)
}

// Find implements MessageUseCase
func (m *storeMessageUseCase) Find(id string) (*StoreMessageModel, error) {
	item, err := m.repository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(id)
	}

	return newStoreMessageModel(item), nil
}

// Update implements MessageUseCase
func (m *storeMessageUseCase) Update(model *StoreMessageUpdateModel) error {
	item, err := m.repository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = item.Set(model.Content)
	if err != nil {
		return err
	}
	return m.repository.Update(item)
}
