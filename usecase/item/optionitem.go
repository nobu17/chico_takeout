package item

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

type OptionItemModel struct {
	Id          string
	Name        string
	Priority    int
	Price       int
	Description string
	Enabled     bool
}

type OptionItemCreateModel struct {
	Name        string
	Priority    int
	Price       int
	Description string
	Enabled     bool
}

type OptionItemUpdateModel struct {
	Id          string
	Name        string
	Priority    int
	Price       int
	Description string
	Enabled     bool
}

func newOptionItemModel(item *domains.OptionItem) *OptionItemModel {
	return &OptionItemModel{
		Id:          item.GetId(),
		Name:        item.GetName(),
		Priority:    item.GetPriority(),
		Price:       item.GetPrice(),
		Description: item.GetDescription(),
		Enabled:     item.GetEnabled(),
	}
}

type OptionItemUseCase interface {
	Find(id string) (*OptionItemModel, error)
	FindAll() ([]OptionItemModel, error)
	Create(model *OptionItemCreateModel) (string, error)
	Update(model *OptionItemUpdateModel) error
	Delete(id string) error
}

type optionItemUseCase struct {
	repository domains.OptionItemRepository
}

func NewOptionItemUseCase(repository domains.OptionItemRepository) OptionItemUseCase {
	return &optionItemUseCase{
		repository: repository,
	}
}

func (i *optionItemUseCase) Find(id string) (*OptionItemModel, error) {
	item, err := i.repository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(id)
	}

	return newOptionItemModel(item), nil
}

func (i *optionItemUseCase) FindAll() ([]OptionItemModel, error) {
	items, err := i.repository.FindAll()
	if err != nil {
		return nil, err
	}

	models := []OptionItemModel{}
	for _, item := range items {
		model := newOptionItemModel(&item)
		models = append(models, *model)
	}

	return models, nil
}

func (i *optionItemUseCase) Create(model *OptionItemCreateModel) (string, error) {
	item, err := domains.NewOptionItem(model.Name, model.Description, model.Priority, model.Price, model.Enabled)
	if err != nil {
		return "", err
	}
	return i.repository.Create(item)
}

func (i *optionItemUseCase) Update(model *OptionItemUpdateModel) error {
	item, err := i.repository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = item.Set(model.Name, model.Description, model.Priority, model.Price, model.Enabled)
	if err != nil {
		return err
	}
	return i.repository.Update(item)
}

func (i *optionItemUseCase) Delete(id string) error {
	item, err := i.repository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.repository.Delete(id)
}
