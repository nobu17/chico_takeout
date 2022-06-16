package items

import (
	"chico/takeout/infrastructures/rdbms"
	"sort"

	domains "chico/takeout/domains/item"

	"gorm.io/gorm"
)

type ItemKindRepository struct {
	db *gorm.DB
}

func NewItemKindRepository(db *gorm.DB) *ItemKindRepository {
	return &ItemKindRepository{
		db: db,
	}
}

type ItemKindModel struct {
	rdbms.BaseModel
	Name     string
	Priority int
}

func (i *ItemKindModel) toDomain() (*domains.ItemKind, error) {
	model, err := domains.NewItemKindForOrm(i.ID, i.Name, i.Priority)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func newItemModel(i *domains.ItemKind) *ItemKindModel {
	model := ItemKindModel{}
	model.ID = i.GetId()
	model.Name = i.GetName()
	model.Priority = i.GetPriority()

	return &model
}

func (i *ItemKindRepository) Find(id string) (*domains.ItemKind, error) {
	model := ItemKindModel{}

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

func (i *ItemKindRepository) FindAll() ([]domains.ItemKind, error) {
	models := []ItemKindModel{}

	err := i.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.ItemKind{}
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

func (i *ItemKindRepository) Create(item *domains.ItemKind) (string, error) {
	model := newItemModel(item)
	err := i.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (i *ItemKindRepository) Update(item *domains.ItemKind) error {
	model := newItemModel(item)
	err := i.db.Save(&model).Error
	return err
}

func (i *ItemKindRepository) Delete(id string) error {
	model := ItemKindModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	err := i.db.Delete(&model).Error
	return err
}
