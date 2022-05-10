package memory

import (
	"fmt"
	"sort"

	domains "chico/takeout/domains/item"
	"chico/takeout/domains/store"

	"github.com/jinzhu/copier"
)

var foodMemory map[string]*domains.FoodItem

type FoodItemMemoryRepository struct {
	inMemory map[string]*domains.FoodItem
	allHours *store.BusinessHours
}

func NewFoodItemMemoryRepository() *FoodItemMemoryRepository {
	if foodMemory == nil {
		resetFoodItemMemory()
	}

	return &FoodItemMemoryRepository{foodMemory, businessHoursMemory}
}

func resetFoodItemMemory() {
	kindRepos := NewItemKindMemoryRepository()
	hourRepos := NewBusinessHoursMemoryRepository()

	allKinds, _ := kindRepos.FindAll()
	allHours, _ := hourRepos.Fetch()
	schedules := allHours.GetSchedules()

	foodMemory = map[string]*domains.FoodItem{}
	scheduleIds1 := []string{schedules[0].GetId(), schedules[1].GetId()}
	item1, _ := domains.NewFoodItem("food1", "item1", 1, 4, 10, 100, allKinds[0].GetId(), scheduleIds1, true)
	foodMemory[item1.GetId()] = item1

	scheduleIds2 := []string{schedules[1].GetId(), schedules[2].GetId()}
	item2, _ := domains.NewFoodItem("food2", "item2", 2, 5, 18, 200, allKinds[1].GetId(), scheduleIds2, true)
	foodMemory[item2.GetId()] = item2
}

func (s *FoodItemMemoryRepository) GetMemory() map[string]*domains.FoodItem {
	return s.inMemory
}

func (s *FoodItemMemoryRepository) GetBusinsHoursMemory() *store.BusinessHours {
	return s.allHours
}

func (s *FoodItemMemoryRepository) Reset() {
	resetFoodItemMemory()
}

func (s *FoodItemMemoryRepository) Find(id string) (*domains.FoodItem, error) {
	if val, ok := s.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.FoodItem{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (s *FoodItemMemoryRepository) FindAll() ([]domains.FoodItem, error) {
	items := []domains.FoodItem{}
	for _, item := range s.inMemory {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].GetPriority() < items[j].GetPriority() })
	return items, nil
}

func (s *FoodItemMemoryRepository) Create(item *domains.FoodItem) (string, error) {
	s.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (s *FoodItemMemoryRepository) Update(item *domains.FoodItem) error {
	if _, ok := s.inMemory[item.GetId()]; ok {
		s.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (s *FoodItemMemoryRepository) Delete(id string) error {
	if _, ok := s.inMemory[id]; ok {
		delete(s.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
