package item

import (
	"fmt"

	"chico/takeout/common"
	domains "chico/takeout/domains/item"
	storeDomains "chico/takeout/domains/store"
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
				ImageUrl:    item.GetImageUrl(),
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

type FoodItemUseCase interface {
	Find(id string) (*FoodItemModel, error)
	FindAll() ([]FoodItemModel, error)
	Create(model *FoodItemCreateModel) (string, error)
	Update(model *FoodItemUpdateModel) error
	Delete(id string) error
}

type foodItemUseCase struct {
	foodItemRepository   domains.FoodItemRepository
	itemKindRepository   domains.ItemKindRepository
	businessHoursService storeDomains.BusinessHoursService
	commonItemUseCase
}

func NewFoodItemUseCase(foodRepos domains.FoodItemRepository,
	itemKindRepos domains.ItemKindRepository,
	businessHoursRepos storeDomains.BusinessHoursRepository) FoodItemUseCase {
	return &foodItemUseCase{
		foodItemRepository:   foodRepos,
		itemKindRepository:   itemKindRepos,
		businessHoursService: *storeDomains.NewBusinessHoursService(businessHoursRepos),
		commonItemUseCase:    *newCommonItemUseCase(itemKindRepos),
	}
}

func (f *foodItemUseCase) Find(id string) (*FoodItemModel, error) {
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

func (f *foodItemUseCase) FindAll() ([]FoodItemModel, error) {
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

func (f *foodItemUseCase) Create(model *FoodItemCreateModel) (string, error) {
	item, err := domains.NewFoodItem(model.Name, model.Description, model.Priority, model.MaxOrder, model.MaxOrderPerDay, model.Price, model.KindId, model.ScheduleIds, model.Enabled, model.ImageUrl)
	if err != nil {
		return "", err
	}
	err = f.ExistsKind(item)
	if err != nil {
		return "", err
	}
	// check schedule is exists
	for _, id := range item.GetScheduleIds() {
		ok, err := f.businessHoursService.ExistsBusinessHour(id)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", common.NewUpdateTargetRelatedNotFoundError(id)
		}
	}

	return f.foodItemRepository.Create(item)
}

func (i *foodItemUseCase) Update(model *FoodItemUpdateModel) error {
	item, err := i.foodItemRepository.Find(model.Id)
	if err != nil {
		return err
	}

	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = item.Set(model.Name, model.Description, model.Priority, model.MaxOrder, model.MaxOrderPerDay, model.Price, model.KindId, model.ScheduleIds, model.Enabled, model.ImageUrl)
	if err != nil {
		return err
	}

	// check schedule is exists
	for _, id := range item.GetScheduleIds() {
		ok, err := i.businessHoursService.ExistsBusinessHour(id)
		if err != nil {
			return err
		}
		if !ok {
			return common.NewUpdateTargetRelatedNotFoundError(id)
		}
	}

	err = i.ExistsKind(item)
	if err != nil {
		return err
	}

	return i.foodItemRepository.Update(item)
}

func (f *foodItemUseCase) Delete(id string) error {
	item, err := f.foodItemRepository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return f.foodItemRepository.Delete(id)
}
