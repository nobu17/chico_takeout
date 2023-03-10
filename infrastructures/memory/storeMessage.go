package memory

import (
	"fmt"

	domains "chico/takeout/domains/message"

	"github.com/jinzhu/copier"
)

var storeMessageMemory map[string]*domains.StoreMessage

type StoreMessageRepository struct {
	inMemory map[string]*domains.StoreMessage
}

func NewStoreMessageRepository() *StoreMessageRepository {
	if storeMessageMemory == nil {
		resetStoreMessageMemory()
	}

	return &StoreMessageRepository{storeMessageMemory}
}

func resetStoreMessageMemory() {
	storeMessageMemory = map[string]*domains.StoreMessage{}
}

func (s *StoreMessageRepository) Reset() {
	resetStoreMessageMemory()
}

func (s *StoreMessageRepository) Find(id string) (*domains.StoreMessage, error) {
	if val, ok := s.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.StoreMessage{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (s *StoreMessageRepository) Create(message *domains.StoreMessage) (string, error) {
	s.inMemory[message.GetId()] = message
	return message.GetId(), nil
}

func (s *StoreMessageRepository) Update(message *domains.StoreMessage) error {
	if _, ok := s.inMemory[message.GetId()]; ok {
		s.inMemory[message.GetId()] = message
		return nil
	}
	return fmt.Errorf("update target not exists")
}
