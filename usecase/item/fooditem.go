package item

import (
	"fmt"

	"chico/takeout/common"
	domains "chico/takeout/domains/item"
)

type FoodItemModel struct {
	CommonItemModel
	ScheduleIds    []string
	MaxOrderPerDay int
}

func newFoodItemModel(item *domains.FoodItem, kind *domains.ItemKind) *FoodItemModel {
	return &FoodItemModel{
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
		ScheduleIds:    item.GetScheduleIds(),
		MaxOrderPerDay: item.GetMaxOrderPerDay(),
	}
}

type FoodItemCreateModel struct {
	CommonItemCreateModel
	ScheduleIds    []string
	MaxOrderPerDay int
}

type FoodItemUpdateModel struct {
	CommonItemUpdateModel
	ScheduleIds    []string
	MaxOrderPerDay int
}

type FoodItemUseCase struct {
	foodItemRepository domains.FoodItemRepository
	itemKindRepository domains.ItemKindRepository
	commonItemUseCase
}

func NewFoodItemUseCase(foodItemRepository domains.FoodItemRepository, itemKindRepository domains.ItemKindRepository) *FoodItemUseCase {
	return &FoodItemUseCase{
		foodItemRepository: foodItemRepository,
		itemKindRepository: itemKindRepository,
		commonItemUseCase:  *newCommonItemUseCase(itemKindRepository),
	}
}

func (f *FoodItemUseCase) Find(id string) (*FoodItemModel, error) {
	item, err := f.foodItemRepository.Find(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, common.NewNotFoundError(fmt.Sprintf("item not found:%s", id))
	}

	kind, err := f.FindKind(item)
	if err != nil {
		return nil, err
	}

	return newFoodItemModel(item, kind), nil
}

func (f *FoodItemUseCase) FindAll() ([]FoodItemModel, error) {
	items, err := f.foodItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	kinds, err := f.itemKindRepository.FindAll()
	if err != nil {
		return nil, err
	}
	models := []FoodItemModel{}
	for _, item := range items {
		for _, kind := range kinds {
			if item.GetKindId() == kind.GetId() {
				model := newFoodItemModel(&item, &kind)
				models = append(models, *model)
				break
			}
		}
	}

	return models, nil
}

func (f *FoodItemUseCase) Create(model FoodItemCreateModel) (string, error) {
	item, err := domains.NewFoodItem(model.Name, model.Description, model.Priority, model.MaxOrder, model.MaxOrderPerDay, model.Price, model.KindId, model.ScheduleIds, model.Enabled)
	if err != nil {
		return "", err
	}
	err = f.ExistsKind(item)
	if err != nil {
		return "", err
	}

	return f.foodItemRepository.Create(*item)
}

func (i *FoodItemUseCase) Update(model FoodItemUpdateModel) error {
	item, err := i.foodItemRepository.Find(model.Id)
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
	return i.foodItemRepository.Update(*item)
}

func (f *FoodItemUseCase) Delete(id string) error {
	item, err := f.foodItemRepository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return f.foodItemRepository.Delete(id)
}
