package item

import (
	"chico/takeout/common"
	"fmt"
)

type FoodItemRepository interface {
	Find(id string) (*FoodItem, error)
	FindAll() ([]FoodItem, error)
	Create(item *FoodItem) (string, error)
	Update(item *FoodItem) error
	Delete(id string) error
}

type FoodItem struct {
	commonItem
	scheduleIds    []string
	maxOrderPerDay MaxOrderPerDay
}

func NewFoodItem(name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string) (*FoodItem, error) {
	common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	if err != nil {
		return nil, err
	}

	item := FoodItem{commonItem: *common}
	err = item.set(maxOrderPerDay, scheduleIds)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func NewFoodItemForOrm(id, name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string) (*FoodItem, error) {
	item, err := NewFoodItem(name, description, priority, maxOrder, maxOrderPerDay, price, kindId, scheduleIds, enabled, imageUrl)
	if err != nil {
		return nil, err
	}
	item.id = id
	return item, nil
}

func (f *FoodItem) Set(name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string) error {
	err := f.commonItem.Set(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	// common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled)
	if err != nil {
		return err
	}
	err = f.set(maxOrderPerDay, scheduleIds)
	if err != nil {
		return err
	}
	return nil
}

func (f *FoodItem) set(maxOrderPerDay int, scheduleIds []string) error {
	maxOrderPValue, err := NewMaxOrderPerDay(maxOrderPerDay, f.maxOrder)
	if err != nil {
		return err
	}
	err = validateScheduleIds(scheduleIds)
	if err != nil {
		return err
	}
	f.maxOrderPerDay = *maxOrderPValue
	f.scheduleIds = scheduleIds
	return nil
}

func validateScheduleIds(scheduleIds []string) error {
	if len(scheduleIds) == 0 {
		return common.NewValidationError("scheduleIds", "empty")
	}

	duplicated := false
	duplicatedId := ""
	encountered := map[string]bool{}
	for i := 0; i < len(scheduleIds); i++ {
		if !encountered[scheduleIds[i]] {
			encountered[scheduleIds[i]] = true
		} else {
			duplicatedId = scheduleIds[i]
			duplicated = true
			break
		}
	}
	if duplicated {
		return common.NewValidationError("scheduleIds", fmt.Sprintf("duplicate Id are not allowed:%s", duplicatedId))
	}
	return nil
}

func (s *FoodItem) GetMaxOrderPerDay() int {
	return s.maxOrderPerDay.GetValue()
}

func (s *FoodItem) GetScheduleIds() []string {
	return s.scheduleIds
}

func (s *FoodItem) HasSameId(id string) bool {
	return s.id == id
}
