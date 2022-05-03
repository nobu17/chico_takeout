package item

import (
	"chico/takeout/common"
	"fmt"
)

type ItemService struct {
	StockItemRepository StockItemRepository
	itemKindRepository  ItemKindRepository
}

func NewItemService(StockItemRepository StockItemRepository, itemKindRepository ItemKindRepository) *ItemService {
	return &ItemService{
		StockItemRepository: StockItemRepository,
		itemKindRepository:  itemKindRepository,
	}
}

func (i *ItemService) ExistsKind(stockItem StockItem) (bool, error) {
	kind, err := i.itemKindRepository.Find(stockItem.kindId)
	if err != nil {
		return false, err
	}
	if kind == nil {
		return false, common.NewNotFoundError(fmt.Sprintf("item kind not found.StockItemId:%s, ItemKindId:%s", stockItem.id, stockItem.kindId))
	}
	return true, nil
}

func (i *ItemService) FindKind(stockItem StockItem) (*ItemKind, error) {
	kind, err := i.itemKindRepository.Find(stockItem.kindId)
	if err != nil {
		return nil, err
	}
	if kind == nil {
		return nil, common.NewNotFoundError(fmt.Sprintf("item kind not found.StockItemId:%s, ItemKindId:%s", stockItem.id, stockItem.kindId))
	}
	return kind, nil
}
