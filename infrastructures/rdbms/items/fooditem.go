package items

import (
	"sort"

	domains "chico/takeout/domains/item"
	"chico/takeout/infrastructures/rdbms"
	"chico/takeout/infrastructures/rdbms/store"

	"gorm.io/gorm"
)

type FoodItemRepository struct {
	db *gorm.DB
}

func NewFoodItemRepository(db *gorm.DB) *FoodItemRepository {
	return &FoodItemRepository{
		db: db,
	}
}

// `FoodItemModel` belongs to `ItemKindModel`, `ItemKindModelID` is the foreign key
// FoodItemModel has and belongs to many BusinessHours, `foodItem_businessHours` is the join table
type FoodItemModel struct {
	rdbms.BaseModel
	Name            string
	Priority        int
	MaxOrder        int
	Price           int
	Description     string
	Enabled         bool
	MaxOrderPerDay  int
	ItemKindModelID string
	ItemKindModel   ItemKindModel
	BusinessHours   []store.BusinessHourModel `gorm:"many2many:foodItem_businessHours;"`
}

func newFoodItemModel(s *domains.FoodItem) *FoodItemModel {
	model := FoodItemModel{}
	model.ID = s.GetId()
	model.Name = s.GetName()
	model.Priority = s.GetPriority()
	model.MaxOrder = s.GetMaxOrder()
	model.Price = s.GetPrice()
	model.Description = s.GetDescription()
	model.Enabled = s.GetEnabled()
	model.MaxOrderPerDay = s.GetMaxOrderPerDay()
	model.ItemKindModelID = s.GetKindId()

	hours := []store.BusinessHourModel{}
	for _, hourId := range s.GetScheduleIds() {
		hour := store.BusinessHourModel{}
		hour.ID = hourId
		hours = append(hours, hour)
	}

	model.BusinessHours = hours

	return &model
}

func (s *FoodItemModel) toDomain() (*domains.FoodItem, error) {
	ids := []string{}
	for _, hour := range s.BusinessHours {
		ids = append(ids, hour.ID)
	}
	model, err := domains.NewFoodItemForOrm(s.ID, s.Name, s.Description, s.Priority, s.MaxOrder, s.MaxOrderPerDay, s.Price, s.ItemKindModelID, ids, s.Enabled)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (f *FoodItemRepository) Find(id string) (*domains.FoodItem, error) {
	model := FoodItemModel{}

	err := f.db.Preload("BusinessHours").First(&model, "ID=?", id).Error
	if err != nil {
		return nil, err
	}

	dom, err := model.toDomain()
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *FoodItemRepository) FindAll() ([]domains.FoodItem, error) {
	models := []FoodItemModel{}

	err := s.db.Preload("BusinessHours").Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.FoodItem{}
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

func (s *FoodItemRepository) Create(item *domains.FoodItem) (string, error) {
	model := newFoodItemModel(item)
	err := s.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (s *FoodItemRepository) Update(item *domains.FoodItem) error {
	model := newFoodItemModel(item)

	var gError error = nil
	s.db.Transaction(func(tx *gorm.DB) error {
		// at first delete all relation
		// err := s.db.Model(&model).Association("BusinessHours").Clear()
		err := s.db.Model(&model).Association("BusinessHours").Replace(model.BusinessHours)
		if err != nil {
			gError = err
			return err
		}
		err = s.db.Save(&model).Error
		if err != nil {
			gError = err
			return err
		}
		return nil
	})

	return gError
}

func (s *FoodItemRepository) Delete(id string) error {
	model := FoodItemModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	var gError error = nil
	s.db.Transaction(func(tx *gorm.DB) error {
		// at first delete all relation
		err := s.db.Model(&model).Association("BusinessHours").Clear()
		if err != nil {
			gError = err
			return err
		}
		err = s.db.Delete(&model).Error
		if err != nil {
			gError = err
			return err
		}

		return nil
	})

	return gError
}
