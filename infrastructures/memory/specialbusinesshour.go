package memory

import (
	"fmt"

	domains "chico/takeout/domains/store"

	"github.com/jinzhu/copier"
)

var specialBusinessHourMemory map[string]*domains.SpecialBusinessHour

type SpecialBusinessHourMemoryRepository struct {
	inMemory map[string]*domains.SpecialBusinessHour
}

func NewSpecialBusinessHourMemoryRepository() *SpecialBusinessHourMemoryRepository {
	if specialBusinessHourMemory == nil {
		resetSpecialBusinessHour()
	}
	return &SpecialBusinessHourMemoryRepository{specialBusinessHourMemory}
}

func resetSpecialBusinessHour() {
	businessHour, _ := NewBusinessHoursMemoryRepository().Fetch()
	schedules := businessHour.GetSchedules()
	specialBusinessHourMemory = map[string]*domains.SpecialBusinessHour{}
	item1, err := domains.NewSpecialBusinessHour("特別日程1", "2022/05/06", "08:00", "12:00", schedules[0].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	specialBusinessHourMemory[item1.GetId()] = item1
	item2, err := domains.NewSpecialBusinessHour("特別日程2", "2022/05/08", "11:00", "14:00", schedules[1].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	specialBusinessHourMemory[item2.GetId()] = item2
}

func (i *SpecialBusinessHourMemoryRepository) Reset() {
	resetSpecialBusinessHour()
}

func (i *SpecialBusinessHourMemoryRepository) GetMemory() map[string]*domains.SpecialBusinessHour {
	return i.inMemory
}

func (i *SpecialBusinessHourMemoryRepository) Find(id string) (*domains.SpecialBusinessHour, error) {
	if val, ok := i.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.SpecialBusinessHour{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (i *SpecialBusinessHourMemoryRepository) FindAll() ([]domains.SpecialBusinessHour, error) {
	items := []domains.SpecialBusinessHour{}
	for _, item := range i.inMemory {
		items = append(items, *item)
	}
	return items, nil
}

func (i *SpecialBusinessHourMemoryRepository) Create(item *domains.SpecialBusinessHour) (string, error) {
	i.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (i *SpecialBusinessHourMemoryRepository) Update(item *domains.SpecialBusinessHour) error {
	if _, ok := i.inMemory[item.GetId()]; ok {
		i.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (b *SpecialBusinessHourMemoryRepository) Delete(id string) error {
	if _, ok := b.inMemory[id]; ok {
		delete(b.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
