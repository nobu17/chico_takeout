package item

import (
	"chico/takeout/common"
	"fmt"
)

type ItemService struct {
	itemKindRepository  ItemKindRepository
}

func NewItemService(itemKindRepository ItemKindRepository) *ItemService {
	return &ItemService{
		itemKindRepository:  itemKindRepository,
	}
}

func (i *ItemService) ExistsKind(item CommonItemImpl) (bool, error) {
	kind, err := i.itemKindRepository.Find(item.GetKindId())
	if err != nil {
		return false, err
	}
	if kind == nil {
		// return false, common.NewNotFoundError(fmt.Sprintf("item kind not found.StockItemId:%s, ItemKindId:%s", item.GetId(), item.GetKindId()))
		return false, nil
	}
	return true, nil
}

func (i *ItemService) FindKind(item CommonItemImpl) (*ItemKind, error) {
	kind, err := i.itemKindRepository.Find(item.GetKindId())
	if err != nil {
		return nil, err
	}
	if kind == nil {
		return nil, common.NewNotFoundError(fmt.Sprintf("item kind not found.StockItemId:%s, ItemKindId:%s", item.GetId(), item.GetKindId()))
	}
	return kind, nil
}
