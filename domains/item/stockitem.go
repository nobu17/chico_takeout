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
	StockItemMaxRemain            = 999
)

func NewStockItem(name, description string, priority, maxOrder, price int, kindId string, enabled bool, imageUrl string) (*StockItem, error) {
	// stock first remain is 0
	remain, err := NewStockRemain(0, StockItemMaxRemain)
	if err != nil {
		return nil, err
	}
	common, err := newCommonItem(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	if err != nil {
		return nil, err
	}
	return &StockItem{commonItem: *common, remain: *remain}, nil
}

// only for orm
func NewStockItemForOrm(id, name, description string, priority, maxOrder, price, remain int, kindId string, enabled bool, imageUrl string) (*StockItem, error) {		
	item, err := NewStockItem(name, description, priority, maxOrder, price, kindId, enabled, imageUrl)
	if err != nil {
		return nil, err
	}
	item.id = id
	
	item.SetRemain(remain);
	if err != nil {
		return nil, err
	}
	return item, nil
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
	// try to consume stock
	remain, err := s.remain.Consume(value)
	if err != nil {
		return err
	}
	s.remain = *remain
	return nil
}

func (s *StockItem) IncreaseRemain(value int) error {
	remain, err := s.remain.Increase(value)
	if err != nil {
		return err
	}
	s.remain = *remain
	return nil
}
