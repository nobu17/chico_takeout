package item

import (
	"chico/takeout/common"
	"chico/takeout/domains/shared/validator"
	"fmt"

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
	id            string
	name          string
	priority      Priority
	optionItemIds []string
}

const (
	ItemNameMaxLength = 15
)

var nameValidator = validator.NewStingLength("ItemKind", ItemNameMaxLength)

// only for orm
func NewItemKindForOrm(id string, name string, priority int, optionItemIds []string) (*ItemKind, error) {
	item := &ItemKind{id: id}
	if err := item.Set(name, priority, optionItemIds); err != nil {
		return nil, err
	}

	return item, nil
}

func NewItemKind(name string, priority int, optionItemIds []string) (*ItemKind, error) {
	item := &ItemKind{id: uuid.NewString()}
	if err := item.Set(name, priority, optionItemIds); err != nil {
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

func (i *ItemKind) GetOptionItemIds() []string {
	return i.optionItemIds
}

func (i *ItemKind) Set(name string, priority int, optionItemIds []string) error {
	if err := nameValidator.Validate(name); err != nil {
		return err
	}

	pri, err := NewPriority(priority)
	if err != nil {
		return err
	}

	err = i.validateOptionIds(optionItemIds)
	if err != nil {
		return err
	}

	i.name = name
	i.priority = *pri
	i.optionItemIds = optionItemIds
	return nil
}

func (i *ItemKind) validateOptionIds(optionItemIds []string) error {
	duplicated := false
	duplicatedId := ""
	encountered := map[string]bool{}
	for i := 0; i < len(optionItemIds); i++ {
		if !encountered[optionItemIds[i]] {
			encountered[optionItemIds[i]] = true
		} else {
			duplicatedId = optionItemIds[i]
			duplicated = true
			break
		}
	}
	if duplicated {
		return common.NewValidationError("optionItemIds", fmt.Sprintf("duplicate Id are not allowed:%s", duplicatedId))
	}
	return nil
}