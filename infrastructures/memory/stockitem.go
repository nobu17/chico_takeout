package memory

import (
	"fmt"
	"sort"

	domains "chico/takeout/domains/item"
	"github.com/jinzhu/copier"
)

var stockMemory map[string]*domains.StockItem

type StockItemMemoryRepository struct {
	inMemory map[string]*domains.StockItem
}

func NewStockItemMemoryRepository() *StockItemMemoryRepository {
	if stockMemory == nil {
		resetStockItemMemory()
	}

	return &StockItemMemoryRepository{stockMemory}
}

func resetStockItemMemory() {
	kindRepos := NewItemKindMemoryRepository()
	allKinds, _ := kindRepos.FindAll()

	stockMemory = map[string]*domains.StockItem{}
	item1, _ := domains.NewStockItem("stock1", "item1", 1, 4, 100, allKinds[0].GetId(), true)
	item1.SetRemain(10)
	stockMemory[item1.GetId()] = item1
	fmt.Println("stock item1:", item1.GetId(), item1.GetKindId())
	item2, _ := domains.NewStockItem("stock2", "item2", 2, 5, 200, allKinds[1].GetId(), true)
	stockMemory[item2.GetId()] = item2
}

func (s *StockItemMemoryRepository) GetMemory() map[string]*domains.StockItem {
	return s.inMemory
}

func (s *StockItemMemoryRepository) Reset() {
	resetStockItemMemory()
}

func (s *StockItemMemoryRepository) Find(id string) (*domains.StockItem, error) {
	if val, ok := s.inMemory[id]; ok {
		// need copy to protect
		duplicated := domains.StockItem{}
		copier.Copy(&duplicated, &val)
		return &duplicated, nil
	}
	return nil, nil
}

func (s *StockItemMemoryRepository) FindAll() ([]domains.StockItem, error) {
	items := []domains.StockItem{}
	for _, item := range s.inMemory {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].GetPriority() < items[j].GetPriority() })
	return items, nil
}

func (s *StockItemMemoryRepository) Create(item *domains.StockItem) (string, error) {
	s.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (s *StockItemMemoryRepository) Update(item *domains.StockItem) error {
	if _, ok := s.inMemory[item.GetId()]; ok {
		s.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (s *StockItemMemoryRepository) Delete(id string) error {
	if _, ok := s.inMemory[id]; ok {
		delete(s.inMemory, id)
		return nil
	}
	return fmt.Errorf("delete target not exists")
}
