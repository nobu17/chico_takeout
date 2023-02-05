package item

import (
	"chico/takeout/common"
	"chico/takeout/domains/shared"
	"fmt"
	"time"
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
	allowDates     AllowDates
}

type AllowDates struct {
	values []shared.Date
}

func (a *AllowDates) GetDates() []string {
	values := []string{}
	for _, val := range a.values {
		values = append(values, val.GetValue())
	}
	return values
}

func (a *AllowDates) GetDatesAsTime() []time.Time {
	values := []time.Time{}
	for _, val := range a.values {
		values = append(values, val.GetAsDate())
	}
	return values
}

func NewAllowDates(dates []string) (*AllowDates, error) {
	// allow empty
	dateV := []shared.Date{}
	for _, date := range dates {
		val, err := shared.NewDate(date)
		if err != nil {
			return nil, err
		}
		dateV = append(dateV, *val)
	}
	return &AllowDates{values: dateV}, nil
}

func NewFoodItem(name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string, allowDates []string) (*FoodItem, error) {
	common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	if err != nil {
		return nil, err
	}

	item := FoodItem{commonItem: *common}
	err = item.set(maxOrderPerDay, scheduleIds, allowDates)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func NewFoodItemForOrm(id, name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string, allowDates []string) (*FoodItem, error) {
	item, err := NewFoodItem(name, description, priority, maxOrder, maxOrderPerDay, price, kindId, scheduleIds, enabled, imageUrl, allowDates)
	if err != nil {
		return nil, err
	}
	item.id = id
	return item, nil
}

func (f *FoodItem) Set(name, description string, priority, maxOrder, maxOrderPerDay, price int, kindId string, scheduleIds []string, enabled bool, imageUrl string, allowDates []string) error {
	err := f.commonItem.Set(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	// common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled)
	if err != nil {
		return err
	}
	err = f.set(maxOrderPerDay, scheduleIds, allowDates)
	if err != nil {
		return err
	}
	return nil
}

func (f *FoodItem) set(maxOrderPerDay int, scheduleIds, allowDates []string) error {
	maxOrderPValue, err := NewMaxOrderPerDay(maxOrderPerDay, f.maxOrder)
	if err != nil {
		return err
	}
	err = validateScheduleIds(scheduleIds)
	if err != nil {
		return err
	}
	dates, err := NewAllowDates(allowDates)
	if err != nil {
		return err
	}
	f.maxOrderPerDay = *maxOrderPValue
	f.scheduleIds = scheduleIds
	f.allowDates = *dates
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

func (s *FoodItem) GetAllowDates() []string {
	return s.allowDates.GetDates()
}

func (s *FoodItem) GetAllowDatesAsTime() []time.Time {
	return s.allowDates.GetDatesAsTime()
}

func (s *FoodItem) HasSameId(id string) bool {
	return s.id == id
}
