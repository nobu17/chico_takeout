package item

import (
	"chico/takeout/domains/shared/validator"

	"github.com/google/uuid"
)

type ItemKindRepository interface {
	Find(id string) (*ItemKind, error)
	FindAll() ([]ItemKind, error)
	Create(item *ItemKind) (string, error)
	Update(item *ItemKind) error
	Delete(id string) error
}

type ItemKind struct {
	id       string
	name     string
	priority Priority
}

const (
	ItemNameMaxLength = 15
)

var nameValidator = validator.NewStingLength("ItemKind", ItemNameMaxLength)

func NewItemKind(name string, priority int) (*ItemKind, error) {
	item := &ItemKind{id: uuid.NewString()}
	if err := item.Set(name, priority); err != nil {
		return nil, err
	}

	return item, nil
}

func (i *ItemKind) GetId() string {
	return i.id
}

func (i *ItemKind) GetName() string {
	return i.name
}

func (i *ItemKind) GetPriority() int {
	return i.priority.GetValue()
}

func (i *ItemKind) Set(name string, priority int) error {
	if err := nameValidator.Validate(name); err != nil {
		return err
	}

	pri, err := NewPriority(priority)
	if err != nil {
		return err
	}

	i.name = name
	i.priority = *pri
	return nil
}
