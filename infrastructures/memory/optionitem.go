package memory

import (
	"fmt"
	"sort"

	domains "chico/takeout/domains/item"

	"github.com/jinzhu/copier"
)

var optionItemMemory map[string]*domains.OptionItem

type OptionItemMemoryRepository struct {
	inMemory map[string]*domains.OptionItem
}

func (i *OptionItemMemoryRepository) GetMemory() map[string]*domains.OptionItem {
	return i.inMemory
}

func resetOptionItemMemory() {
	optionItemMemory = map[string]*domains.OptionItem{}
	item1, _ := domains.NewOptionItemForOrm("1", "item1", "memo1", 1, 100, true)
	optionItemMemory[item1.GetId()] = item1
	item2, _ := domains.NewOptionItemForOrm("2", "item2", "memo2", 2, 200, true)
	optionItemMemory[item2.GetId()] = item2
	item3, _ := domains.NewOptionItemForOrm("3", "item3", "memo3", 3, 300, false)
	optionItemMemory[item3.GetId()] = item3
}

func NewOptionItemMemoryRepository() *OptionItemMemoryRepository {
	if optionItemMemory == nil {
		resetOptionItemMemory()
	}
	return &OptionItemMemoryRepository{optionItemMemory}
}

func NewOptionItemMemoryRepositoryWithParam(param map[string]*domains.OptionItem) *OptionItemMemoryRepository {
	optionItemMemory = param
	return &OptionItemMemoryRepository{optionItemMemory}
}

func (i *OptionItemMemoryRepository) Reset() {
	resetOptionItemMemory()
}

func (i *OptionItemMemoryRepository) Find(id string) (*domains.OptionItem, error) {
	if val, ok := i.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.OptionItem{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (i *OptionItemMemoryRepository) FindAll() ([]domains.OptionItem, error) {
	items := []domains.OptionItem{}
	for _, item := range i.inMemory {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].GetPriority() < items[j].GetPriority() })
	return items, nil
}

func (i *OptionItemMemoryRepository) Create(item *domains.OptionItem) (string, error) {
	i.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (i *OptionItemMemoryRepository) Update(item *domains.OptionItem) error {
	if _, ok := i.inMemory[item.GetId()]; ok {
		i.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (b *OptionItemMemoryRepository) Delete(id string) error {
	if _, ok := b.inMemory[id]; ok {
		delete(b.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
