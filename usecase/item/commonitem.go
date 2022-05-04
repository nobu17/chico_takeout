package item

import (
	"fmt"

	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

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
		return common.NewUpdateTargetNotFoundError(fmt.Sprintf("kind id is not exists.kind id:%s", item.GetKindId()))
	}
	return nil
}

func (c *commonItemUseCase) FindKind(item domains.CommonItemImpl) (*domains.ItemKind, error) {
	return c.itemService.FindKind(item)
}
