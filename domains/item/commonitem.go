package item

import (
	"fmt"
	"strings"

	"chico/takeout/common"

	"github.com/google/uuid"
)

type commonItem struct {
	id          string
	name        string
	priority    Priority
	maxOrder    MaxOrder
	price       Price
	description string
	kindId      string
}

type CommonItemImpl interface {
	GetId() string
	GetKindId() string
}

func newCommonItem(name, description string, priority, maxOrder, price int, kindId string) (*commonItem, error) {
	item := commonItem{id: uuid.NewString()}
	err := item.Set(name, description, priority, maxOrder, price, kindId)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *commonItem) Set(name, description string, priority, maxOrder, price int, kindId string) error {
	if err := s.validateName(name); err != nil {
		return err
	}

	priV, err := NewPriority(priority)
	if err != nil {
		return err
	}

	maxO, err := NewMaxOrder(maxOrder)
	if err != nil {
		return err
	}

	priceV, err := NewPrice(price, StockItemMaxPrice)
	if err != nil {
		return err
	}

	if err := s.validateDescription(description); err != nil {
		return err
	}

	if err := s.validateKindId(kindId); err != nil {
		return err
	}

	s.name = name
	s.priority = *priV
	s.maxOrder = *maxO
	s.price = *priceV
	s.description = description
	s.kindId = kindId
	return nil
}

func (s *commonItem) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return common.NewValidationError("name", "required")
	}

	if len(name) > StockItemNameMaxLength {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", StockItemNameMaxLength))
	}
	return nil
}

func (s *commonItem) validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return common.NewValidationError("description", "required")
	}

	if len(description) > StockItemDescriptionMaxLength {
		return common.NewValidationError("description", fmt.Sprintf("MaxLength:%d", StockItemDescriptionMaxLength))
	}
	return nil
}

func (s *commonItem) validateKindId(kindId string) error {
	if strings.TrimSpace(kindId) == "" {
		return common.NewValidationError("kindId", "required")
	}
	return nil
}

func (s *commonItem) HasSameId(id string) bool {
	return s.id == id
}

func (s *commonItem) WithInMaxOrder(quantity int) error {
	return s.maxOrder.WithinLimit(quantity)
}

func (s *commonItem) GetId() string {
	return s.id
}

func (s *commonItem) GetName() string {
	return s.name
}

func (s *commonItem) GetPriority() int {
	return s.priority.value
}

func (s *commonItem) GetMaxOrder() int {
	return s.maxOrder.value
}

func (s *commonItem) GetPrice() int {
	return s.price.value
}

func (s *commonItem) GetDescription() string {
	return s.description
}

func (s *commonItem) GetKindId() string {
	return s.kindId
}
