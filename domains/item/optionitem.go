package item

import (
	"github.com/google/uuid"
)

type OptionItemRepository interface {
	Find(id string) (*OptionItem, error)
	FindAll() ([]OptionItem, error)
	Create(item *OptionItem) (string, error)
	Update(item *OptionItem) error
	Delete(id string) error
}

type OptionItem struct {
	id          string
	name        Name
	priority    Priority
	price       Price
	description Description
	enabled     bool
}

const (
	OptionItemMaxPrice             = 20000
	OptionItemDescriptionMaxLength = 150
	OptionItemNameMaxLength        = 25
)

func NewOptionItem(name, description string, priority, price int, enabled bool) (*OptionItem, error) {
	item := OptionItem{id: uuid.NewString()}
	err := item.Set(name, description, priority, price, enabled)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// only for orm
func NewOptionItemForOrm(id, name, description string, priority, price int, enabled bool) (*OptionItem, error) {
	item := &OptionItem{id: id}
	err := item.Set(name, description, priority, price, enabled)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (i *OptionItem) GetId() string {
	return i.id
}

func (i *OptionItem) GetName() string {
	return i.name.GetValue()
}

func (i *OptionItem) GetPriority() int {
	return i.priority.GetValue()
}

func (s *OptionItem) GetPrice() int {
	return s.price.GetValue()
}

func (s *OptionItem) GetDescription() string {
	return s.description.GetValue()
}

func (s *OptionItem) GetEnabled() bool {
	return s.enabled
}

func (o *OptionItem) Set(name, description string, priority, price int, enabled bool) error {
	nameV, err := NewName(name, OptionItemNameMaxLength)
	if err != nil {
		return err
	}

	priV, err := NewPriority(priority)
	if err != nil {
		return err
	}

	priceV, err := NewPrice(price, OptionItemMaxPrice)
	if err != nil {
		return err
	}

	desc, err := NewDescription(description, OptionItemNameMaxLength)
	if err != nil {
		return err
	}

	o.name = *nameV
	o.priority = *priV
	o.price = *priceV
	o.description = *desc
	o.enabled = enabled

	return nil
}

func (s *OptionItem) HasSameId(id string) bool {
	return s.id == id
}