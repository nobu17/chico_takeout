package item

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

type CommonItemModel struct {
	Id   string
	Kind ItemKindModel
	CommonItemBaseModel
}

type CommonItemCreateModel struct {
	KindId string
	CommonItemBaseModel
}

type CommonItemUpdateModel struct {
	Id     string
	KindId string
	CommonItemBaseModel
}

type CommonItemBaseModel struct {
	Name        string
	Priority    int
	MaxOrder    int
	Price       int
	Description string
	Enabled     bool
	ImageUrl    string
}

type commonItemUseCase struct {
	itemService domains.ItemService
}

func newCommonItemUseCase(itemKindRepository domains.ItemKindRepository) *commonItemUseCase {
	return &commonItemUseCase{
		itemService: *domains.NewItemService(itemKindRepository),
	}
}

func (c *commonItemUseCase) ExistsKind(item domains.CommonItemImpl) error {
	exists, err := c.itemService.ExistsKind(item)
	if err != nil {
		return err
	}
	if !exists {
		return common.NewUpdateTargetRelatedNotFoundError(item.GetKindId())
	}
	return nil
}

func (c *commonItemUseCase) FindKind(item domains.CommonItemImpl) (*domains.ItemKind, error) {
	return c.itemService.FindKind(item)
}
