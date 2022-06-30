package item

import (
	"strings"

	"chico/takeout/common"

	"github.com/google/uuid"
)

const (
	CommonItemMaxPrice             = 20000
	CommonItemDescriptionMaxLength = 150
	CommonItemNameMaxLength        = 15
)

type commonItem struct {
	id          string
	name        Name
	priority    Priority
	maxOrder    MaxOrder
	price       Price
	description Description
	kindId      string
	enabled     bool
	imageUrl    ImageUrl
}

type CommonItemImpl interface {
	GetId() string
	GetKindId() string
}

func newCommonItem(name, description string, priority, maxOrder, price int, kindId string, enabled bool, imageUrl string) (*commonItem, error) {
	item := commonItem{id: uuid.NewString()}
	err := item.Set(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *commonItem) Set(name, description string, priority, maxOrder, price int, kindId string, enabled bool, imageUrl string) error {
	nameV, err := NewName(name, CommonItemNameMaxLength)
	if err != nil {
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

	priceV, err := NewPrice(price, CommonItemMaxPrice)
	if err != nil {
		return err
	}

	desc, err := NewDescription(description, CommonItemDescriptionMaxLength)
	if err != nil {
		return err
	}

	if err := s.validateKindId(kindId); err != nil {
		return err
	}

	image, err := NewImageUrl(imageUrl)
	if err != nil {
		return err
	}

	s.name = *nameV
	s.priority = *priV
	s.maxOrder = *maxO
	s.price = *priceV
	s.description = *desc
	s.kindId = kindId
	s.enabled = enabled
	s.imageUrl = *image
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
	return s.name.GetValue()
}

func (s *commonItem) GetPriority() int {
	return s.priority.GetValue()
}

func (s *commonItem) GetMaxOrder() int {
	return s.maxOrder.GetValue()
}

func (s *commonItem) GetPrice() int {
	return s.price.GetValue()
}

func (s *commonItem) GetDescription() string {
	return s.description.GetValue()
}

func (s *commonItem) GetKindId() string {
	return s.kindId
}

func (s *commonItem) GetEnabled() bool {
	return s.enabled
}

func (s *commonItem) GetImageUrl() string {
	return s.imageUrl.GetValue()
}

func (s *commonItem) HasKind(kind ItemKind) bool {
	return s.kindId == kind.GetId()
}
