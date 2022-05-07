package item

import (
	"fmt"

	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

type StockItemModel struct {
	CommonItemModel
	Remain int
}

type StockItemCreateModel struct {
	CommonItemCreateModel
}

type StockItemUpdateModel struct {
	CommonItemUpdateModel
}

type StockItemRemainUpdateModel struct {
	Id     string
	Remain int
}

func newStockItemModel(item *domains.StockItem, kind *domains.ItemKind) *StockItemModel {
	return &StockItemModel{
		CommonItemModel: CommonItemModel{
			Id:   item.GetId(),
			Kind: *newItemKindModel(kind),
			CommonItemBaseModel: CommonItemBaseModel{
				Name:        item.GetName(),
				Priority:    item.GetPriority(),
				MaxOrder:    item.GetMaxOrder(),
				Price:       item.GetPrice(),
				Description: item.GetDescription(),
				Enabled:     item.GetEnabled(),
			},
		},
		Remain: item.GetRemain(),
	}
}

type StockItemUseCase struct {
	stockItemRepository domains.StockItemRepository
	itemKindRepository  domains.ItemKindRepository
	commonItemUseCase
}

func NewStockItemUseCase(stockItemRepository domains.StockItemRepository, itemKindRepository domains.ItemKindRepository) *StockItemUseCase {
	return &StockItemUseCase{
		stockItemRepository: stockItemRepository,
		itemKindRepository:  itemKindRepository,
		commonItemUseCase:   *newCommonItemUseCase(itemKindRepository),
	}
}

func (i *StockItemUseCase) Find(id string) (*StockItemModel, error) {
	item, err := i.stockItemRepository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(fmt.Sprintf("item not found:%s", id))
	}

	kind, err := i.FindKind(item)
	if err != nil {
		return nil, err
	}

	return newStockItemModel(item, kind), nil
}

func (i *StockItemUseCase) FindAll() ([]StockItemModel, error) {
	items, err := i.stockItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	kinds, err := i.itemKindRepository.FindAll()
	if err != nil {
		return nil, err
	}
	models := []StockItemModel{}
	for _, item := range items {
		for _, kind := range kinds {
			if item.GetKindId() == kind.GetId() {
				model := newStockItemModel(&item, &kind)
				models = append(models, *model)
				break
			}
		}
	}

	return models, nil
}

func (i *StockItemUseCase) Create(model StockItemCreateModel) (string, error) {
	item, err := domains.NewStockItem(model.Name, model.Description, model.Priority, model.MaxOrder, model.Price, model.KindId, model.Enabled)
	if err != nil {
		return "", err
	}

	err = i.ExistsKind(item)
	if err != nil {
		return "", err
	}

	return i.stockItemRepository.Create(*item)
}

func (i *StockItemUseCase) Update(model StockItemUpdateModel) error {
	item, err := i.stockItemRepository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = i.ExistsKind(item)
	if err != nil {
		return err
	}

	err = item.Set(model.Name, model.Description, model.Priority, model.MaxOrder, model.Price, model.KindId, model.Enabled)
	if err != nil {
		return err
	}
	return i.stockItemRepository.Update(*item)
}

func (i *StockItemUseCase) Delete(id string) error {
	item, err := i.stockItemRepository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.stockItemRepository.Delete(id)
}

func (i *StockItemUseCase) UpdateRemain(model StockItemRemainUpdateModel) error {
	item, err := i.stockItemRepository.Find(model.Id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = item.SetRemain(model.Remain)
	if err != nil {
		return err
	}
	return i.stockItemRepository.Update(*item)
}
