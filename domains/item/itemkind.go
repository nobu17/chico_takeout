package item

import (
	"fmt"
	"strings"

	"chico/takeout/common"

	"github.com/google/uuid"
)

type ItemKindRepository interface {
	Find(id string) (*ItemKind, error)
	FindAll() ([]ItemKind, error)
	Create(item ItemKind) (string, error)
	Update(item ItemKind) error
	Delete(id string) error
}

type ItemKind struct {
	id       string
	name     string
	priority Priority
}

const (
	ItemNameMaxLength = 10
)

func NewItemKind(name string, priority int) (*ItemKind, error) {

	if err := validateItemKindName(name); err != nil {
		return nil, err
	}

	pri, err := NewPriority(priority)
	if err != nil {
		return nil, err
	}

	return &ItemKind{id: uuid.NewString(), name: name, priority: *pri}, nil
}

func validateItemKindName(name string) error {
	if strings.TrimSpace(name) == "" {
		return common.NewValidationError("name", "required")
	}

	if len(name) > ItemNameMaxLength {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", ItemNameMaxLength))
	}
	return nil
}

func (i *ItemKind) GetId() string {
	return i.id
}

func (i *ItemKind) GetName() string {
	return i.name
}

func (i *ItemKind) GetPriority() int {
	return i.priority.value
}

func (i *ItemKind) Set(name string, priority int) error {
	if err := validateItemKindName(name); err != nil {
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
