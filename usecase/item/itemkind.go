package item

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

type ItemKindModel struct {
	Id       string
	Name     string
	Priority int
	OptionItemIds []string
}

func newItemKindModel(item *domains.ItemKind) *ItemKindModel {
	return &ItemKindModel{
		Id:       item.GetId(),
		Name:     item.GetName(),
		Priority: item.GetPriority(),
		OptionItemIds: item.GetOptionItemIds(),
	}
}

type ItemKindCreateModel struct {
	Name          string
	Priority      int
	OptionItemIds []string
}

type ItemKindUpdateModel struct {
	Id            string
	Name          string
	Priority      int
	OptionItemIds []string
}

type ItemKindUseCase interface {
	Find(id string) (*ItemKindModel, error)
	FindAll() ([]ItemKindModel, error)
	Create(model *ItemKindCreateModel) (string, error)
	Update(model *ItemKindUpdateModel) error
	Delete(id string) error
}

type itemKindUseCase struct {
	repository domains.ItemKindRepository
	optionItemRepository domains.OptionItemRepository
}

func NewItemKindUseCase(repository domains.ItemKindRepository, optionItemRepository domains.OptionItemRepository) ItemKindUseCase {
	return &itemKindUseCase{
		repository: repository,
		optionItemRepository: optionItemRepository,
	}
}

func (i *itemKindUseCase) Find(id string) (*ItemKindModel, error) {
	item, err := i.repository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(id)
	}

	return newItemKindModel(item), nil
}

func (i *itemKindUseCase) FindAll() ([]ItemKindModel, error) {
	items, err := i.repository.FindAll()
	if err != nil {
		return nil, err
	}

	models := []ItemKindModel{}
	for _, item := range items {
		model := newItemKindModel(&item)
		models = append(models, *model)
	}

	return models, nil
}

func (i *itemKindUseCase) Create(model *ItemKindCreateModel) (string, error) {
	err := i.checkOptionItemExists(model.OptionItemIds)
	if err != nil {
		return "", err
	}

	item, err := domains.NewItemKind(model.Name, model.Priority, model.OptionItemIds)
	if err != nil {
		return "", err
	}
	return i.repository.Create(item)
}

func (i *itemKindUseCase) Update(model *ItemKindUpdateModel) error {
	item, err := i.repository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = i.checkOptionItemExists(model.OptionItemIds)
	if err != nil {
		return err
	}

	err = item.Set(model.Name, model.Priority, model.OptionItemIds)
	if err != nil {
		return err
	}
	return i.repository.Update(item)
}

func (i *itemKindUseCase) Delete(id string) error {
	item, err := i.repository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.repository.Delete(id)
}

func (i *itemKindUseCase) checkOptionItemExists(ids []string) error {
	for _, id := range ids {
		item, err := i.optionItemRepository.Find(id)
		if err != nil {
			return err
		}
		if item == nil {
			return common.NewRelatedItemNotFoundError(id)
		}
	}
	return nil
}
