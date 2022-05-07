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

func NewFoodItem(name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool) (*FoodItem, error) {
	common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled)
	if err != nil {
		return nil, err
	}
	maxOrderPValue, err := NewMaxOrderPerDay(maxOrderPerDay)
	if err != nil {
		return nil, err
	}
	err = validateScheduleIds(scheduleIds)
	if err != nil {
		return nil, err
	}
	item := FoodItem{commonItem: *common, maxOrderPerDay: *maxOrderPValue, scheduleIds: scheduleIds}
	return &item, nil
}

func validateScheduleIds(scheduleIds []string) error {
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
	return s.maxOrderPerDay.value
}

func (s *FoodItem) GetScheduleIds() []string {
	return s.scheduleIds
}

func (s *FoodItem) HasSameId(id string) bool {
	return s.id == id
}