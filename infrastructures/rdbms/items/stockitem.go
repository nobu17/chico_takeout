package items

import (
	"chico/takeout/infrastructures/rdbms"
	"sort"

	domains "chico/takeout/domains/item"

	"gorm.io/gorm"
)

type StockItemRepository struct {
	db *gorm.DB
}

func NewStockItemRepository(db *gorm.DB) *StockItemRepository {
	return &StockItemRepository{
		db: db,
	}
}

// `StockItemModel` belongs to `ItemKindModel`, `ItemKindModelID` is the foreign key
type StockItemModel struct {
	rdbms.BaseModel
	Name            string
	Priority        int
	MaxOrder        int
	Price           int
	Description     string
	Enabled         bool
	Remain          int
	ItemKindModelID string
	ItemKindModel   ItemKindModel
}

func newStockItemModel(s *domains.StockItem) *StockItemModel {
	model := StockItemModel{}
	model.ID = s.GetId()
	model.Name = s.GetName()
	model.Priority = s.GetPriority()
	model.MaxOrder = s.GetMaxOrder()
	model.Price = s.GetPrice()
	model.Description = s.GetDescription()
	model.Enabled = s.GetEnabled()
	model.Remain = s.GetRemain()
	model.ItemKindModelID = s.GetKindId()

	return &model
}

func (s *StockItemModel) toDomain() (*domains.StockItem, error) {
	model, err := domains.NewStockItemForOrm(s.ID, s.Name, s.Description, s.Priority, s.MaxOrder, s.Price, s.Remain, s.ItemKindModelID, s.Enabled)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *StockItemRepository) Find(id string) (*domains.StockItem, error) {
	model := StockItemModel{}

	err := s.db.First(&model, "ID=?", id).Error
	if err != nil {
		return nil, err
	}

	dom, err := model.toDomain()
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *StockItemRepository) FindAll() ([]domains.StockItem, error) {
	models := []StockItemModel{}

	err := s.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.StockItem{}
	for _, model := range models {
		item, err := model.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].GetPriority() < items[j].GetPriority() })
	return items, nil
}

func (s *StockItemRepository) Create(item *domains.StockItem) (string, error) {
	model := newStockItemModel(item)
	err := s.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (s *StockItemRepository) Update(item *domains.StockItem) error {
	model := newStockItemModel(item)
	err := s.db.Save(&model).Error
	return err
}

func (s *StockItemRepository) Delete(id string) error {
	model := StockItemModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	err := s.db.Delete(&model).Error
	return err
}
