package items

import (
	domains "chico/takeout/domains/item"
	"chico/takeout/infrastructures/rdbms"
	"sort"

	"gorm.io/gorm"
)

type OptionItemRepository struct {
	db *gorm.DB
}

func NewOptionItemRepository(db *gorm.DB) *OptionItemRepository {
	return &OptionItemRepository{
		db: db,
	}
}

type OptionItemModel struct {
	rdbms.BaseModel
	Name        string
	Priority    int
	Price       int
	Description string
	Enabled     bool
}

func (i *OptionItemModel) toDomain() (*domains.OptionItem, error) {
	model, err := domains.NewOptionItemForOrm(i.ID, i.Name, i.Description, i.Priority, i.Price, i.Enabled)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func newOptionItemModel(i *domains.OptionItem) *OptionItemModel {
	model := OptionItemModel{}
	model.ID = i.GetId()
	model.Name = i.GetName()
	model.Description = i.GetDescription()
	model.Priority = i.GetPriority()
	model.Price = i.GetPrice()
	model.Enabled = i.GetEnabled()

	return &model
}

func (i *OptionItemRepository) Find(id string) (*domains.OptionItem, error) {
	model := OptionItemModel{}

	err := i.db.First(&model, "ID=?", id).Error
	if err != nil {
		return nil, err
	}

	dom, err := model.toDomain()
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (i *OptionItemRepository) FindAll() ([]domains.OptionItem, error) {
	models := []OptionItemModel{}

	err := i.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.OptionItem{}
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

func (i *OptionItemRepository) Create(item *domains.OptionItem) (string, error) {
	model := newOptionItemModel(item)
	err := i.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (i *OptionItemRepository) Update(item *domains.OptionItem) error {
	model := newOptionItemModel(item)
	err := i.db.Save(&model).Error
	return err
}

func (i *OptionItemRepository) Delete(id string) error {
	model := OptionItemModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	err := i.db.Delete(&model).Error
	return err
}
