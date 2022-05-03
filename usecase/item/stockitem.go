package item

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/item"
	"fmt"
)

type StockItemModel struct {
	Id          string
	Name        string
	Priority    int
	MaxOrder    int
	Price       int
	Description string
	Kind        ItemKindModel
	Remain      int
}

type StockItemCreateModel struct {
	Name        string
	Priority    int
	MaxOrder    int
	Price       int
	Description string
	KindId      string
}

type StockItemUpdateModel struct {
	Id          string
	Name        string
	Priority    int
	MaxOrder    int
	Price       int
	Description string
	KindId      string
}

type StockItemRemainUpdateModel struct {
	Id     string
	Remain int
}

func newStockItemModel(item *domains.StockItem, kind *domains.ItemKind) *StockItemModel {
	return &StockItemModel{
		Id:          item.GetId(),
		Name:        item.GetName(),
		Priority:    item.GetPriority(),
		MaxOrder:    item.GetMaxOrder(),
		Price:       item.GetPrice(),
		Description: item.GetDescription(),
		Kind:        *newItemKindModel(kind),
		Remain:      item.GetRemain(),
	}
}

type StockItemUseCase struct {
	StockItemRepository domains.StockItemRepository
	itemKindRepository  domains.ItemKindRepository
	itemService         domains.ItemService
}

func NewStockItemUseCase(StockItemRepository domains.StockItemRepository, itemKindRepository domains.ItemKindRepository) *StockItemUseCase {
	return &StockItemUseCase{
		StockItemRepository: StockItemRepository,
		itemKindRepository:  itemKindRepository,
		itemService:         *domains.NewItemService(StockItemRepository, itemKindRepository),
	}
}

func (i *StockItemUseCase) Find(id string) (*StockItemModel, error) {
	item, err := i.StockItemRepository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(fmt.Sprintf("user not found:%s", id))
	}

	kind, err := i.itemService.FindKind(*item)
	if err != nil {
		return nil, err
	}

	return newStockItemModel(item, kind), nil
}

func (i *StockItemUseCase) FindAll() ([]StockItemModel, error) {
	items, err := i.StockItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	kinds, err := i.itemKindRepository.FindAll()
	if err != nil {
		return nil, err
	}
	models := []StockItemModel{}
	for _, item := range items {
		fmt.Println("item", item.GetKindId())
		for _, kind := range kinds {
			fmt.Println("kind", kind.GetId())
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
	item, err := domains.NewStockItem(model.Name, model.Description, model.Priority, model.MaxOrder, model.Price, model.KindId)
	if err != nil {
		return "", err
	}
	exists, err := i.itemService.ExistsKind(*item)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", common.NewUpdateTargetNotFoundError(fmt.Sprintf("kind id is not exists.kind id:%s", model.KindId))
	}

	return i.StockItemRepository.Create(*item)
}

func (i *StockItemUseCase) Update(model StockItemUpdateModel) error {
	item, err := i.StockItemRepository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	exists, err := i.itemService.ExistsKind(*item)
	if err != nil {
		return err
	}
	if !exists {
		return common.NewUpdateTargetNotFoundError(fmt.Sprintf("kind id is not exists.kind id:%s", model.KindId))
	}

	err = item.Set(model.Name, model.Description, model.Priority, model.MaxOrder, model.Price, model.KindId)
	if err != nil {
		return err
	}
	return i.StockItemRepository.Update(*item)
}

func (i *StockItemUseCase) Delete(id string) error {
	item, err := i.StockItemRepository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.StockItemRepository.Delete(id)
}

func (i *StockItemUseCase) UpdateRemain(model StockItemRemainUpdateModel) error {
	item, err := i.StockItemRepository.Find(model.Id)
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
	return i.StockItemRepository.Update(*item)
}

// todo1: kind and itemkind の存在チェックはserv or queryにする
