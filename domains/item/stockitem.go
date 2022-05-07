package item

type StockItemRepository interface {
	Find(id string) (*StockItem, error)
	FindAll() ([]StockItem, error)
	Create(item *StockItem) (string, error)
	Update(item *StockItem) error
	Delete(id string) error
}

type StockItem struct {
	commonItem
	remain StockRemain
}

const (
	StockItemNameMaxLength        = 15
	StockItemDescriptionMaxLength = 150
	StockItemMaxOrderMaxValue     = 30
	StockItemMaxPrice             = 20000
	StockItemMaxRemain            = 999
)

func NewStockItem(name, description string, priority, maxOrder, price int, kindId string, enabled bool) (*StockItem, error) {
	// stock first remain is 0
	remain, _ := NewStockRemain(0, StockItemMaxRemain)
	common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled)
	if err != nil {
		return nil, err
	}
	item := StockItem{commonItem: *common, remain: *remain}
	return &item, nil
}

func (s *StockItem) GetRemain() int {
	return s.remain.GetValue()
}

func (s *StockItem) SetRemain(value int) error {
	remain, err := NewStockRemain(value, StockItemMaxRemain)
	if err != nil {
		return err
	}
	s.remain = *remain
	return nil
}

func (s *StockItem) ConsumeRemain(value int) error {
	// check max order request at first
	if err := s.maxOrder.WithinLimit(value); err != nil {
		return nil
	}
	// try to consume stock
	remain, err := s.remain.Consume(value)
	if err != nil {
		return err
	}
	s.remain = *remain
	return nil
}

func (s *StockItem) IncreseRemain(value int) error {
	remain, err := s.remain.Increase(value)
	if err != nil {
		return err
	}
	s.remain = *remain
	return nil
}
