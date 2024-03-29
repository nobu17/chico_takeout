package memory

import (
	"fmt"
	"sort"

	domains "chico/takeout/domains/item"

	"github.com/jinzhu/copier"
)

var memory map[string]*domains.ItemKind

type ItemKindMemoryRepository struct {
	inMemory map[string]*domains.ItemKind
}

func (i *ItemKindMemoryRepository) GetMemory() map[string]*domains.ItemKind {
	return i.inMemory
}

func resetKindMemory() {
	memory = map[string]*domains.ItemKind{}
	ids1 := []string{}
	item1, _ := domains.NewItemKind("item1", 1, ids1)
	memory[item1.GetId()] = item1
	ids2 := []string{}
	item2, _ := domains.NewItemKind("item2", 2, ids2)
	memory[item2.GetId()] = item2
}

func NewItemKindMemoryRepository() *ItemKindMemoryRepository {
	if memory == nil {
		resetKindMemory()
	}
	return &ItemKindMemoryRepository{memory}
}

func NewItemKindMemoryRepositoryWithParam(param map[string]*domains.ItemKind) *ItemKindMemoryRepository {
	memory = param
	return &ItemKindMemoryRepository{memory}
}

func (i *ItemKindMemoryRepository) Reset() {
	resetKindMemory()
}

func (i *ItemKindMemoryRepository) Find(id string) (*domains.ItemKind, error) {
	if val, ok := i.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.ItemKind{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (i *ItemKindMemoryRepository) FindAll() ([]domains.ItemKind, error) {
	items := []domains.ItemKind{}
	for _, item := range i.inMemory {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].GetPriority() < items[j].GetPriority() })
	return items, nil
}

func (i *ItemKindMemoryRepository) Create(item *domains.ItemKind) (string, error) {
	i.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (i *ItemKindMemoryRepository) Update(item *domains.ItemKind) error {
	if _, ok := i.inMemory[item.GetId()]; ok {
		i.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (b *ItemKindMemoryRepository) Delete(id string) error {
	if _, ok := b.inMemory[id]; ok {
		delete(b.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
