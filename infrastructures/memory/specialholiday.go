package memory

import (
	"fmt"

	domains "chico/takeout/domains/store"
)

var specialHolidayMemory map[string]*domains.SpecialHoliday

type SpecialHolidayMemoryRepository struct {
	inMemory map[string]*domains.SpecialHoliday
}

func NewSpecialHolidayMemoryRepository() *SpecialHolidayMemoryRepository {
	if specialHolidayMemory == nil {
		specialHolidayMemory = map[string]*domains.SpecialHoliday{}
		item1, err := domains.NewSpecialHoliday("おやすみ１", "2022/05/06", "2022/06/03")
		if err != nil {
			fmt.Println(err)
			panic("failed to create special holiday")
		}
		specialHolidayMemory[item1.GetId()] = item1
		item2, err := domains.NewSpecialHoliday("おやすみ2", "2022/07/06", "2022/08/01")
		if err != nil {
			panic("failed to create special holiday")
		}
		specialHolidayMemory[item2.GetId()] = item2
	}
	return &SpecialHolidayMemoryRepository{specialHolidayMemory}
}

func (i *SpecialHolidayMemoryRepository) Find(id string) (*domains.SpecialHoliday, error) {
	if val, ok := i.inMemory[id]; ok {
		return val, nil
	}
	return nil, nil
}

func (i *SpecialHolidayMemoryRepository) FindAll() ([]domains.SpecialHoliday, error) {
	items := []domains.SpecialHoliday{}
	for _, item := range i.inMemory {
		items = append(items, *item)
	}
	return items, nil
}

func (i *SpecialHolidayMemoryRepository) Create(item domains.SpecialHoliday) (string, error) {
	i.inMemory[item.GetId()] = &item
	return item.GetId(), nil
}

func (i *SpecialHolidayMemoryRepository) Update(item domains.SpecialHoliday) error {
	if _, ok := i.inMemory[item.GetId()]; ok {
		i.inMemory[item.GetId()] = &item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (b *SpecialHolidayMemoryRepository) Delete(id string) error {
	if _, ok := b.inMemory[id]; ok {
		delete(b.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
