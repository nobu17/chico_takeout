package store

import (
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/store"
	"chico/takeout/infrastructures/rdbms"

	"gorm.io/gorm"
)

type SpecialHolidayRepository struct {
	db *gorm.DB
}

func NewSpecialHolidayRepository(db *gorm.DB) *SpecialHolidayRepository {
	return &SpecialHolidayRepository{
		db: db,
	}
}

type SpecialHolidayModel struct {
	rdbms.BaseModel
	Name  string
	Start *time.Time
	End   *time.Time
}

func (s *SpecialHolidayModel) toDomain() (*domains.SpecialHoliday, error) {
	startStr := common.ConvertTimeToDateStr(*s.Start)
	endStr := common.ConvertTimeToDateStr(*s.End)
	model, err := domains.NewSpecialHolidayForOrm(s.ID, s.Name, startStr, endStr)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func newSpecialHolidayModel(s *domains.SpecialHoliday) (*SpecialHolidayModel, error) {
	model := SpecialHolidayModel{}
	model.ID = s.GetId()
	model.Name = s.GetName()

	start, err := common.ConvertStrToDate(s.GetStart())
	if err != nil {
		return nil, err
	}
	model.Start = start

	end, err := common.ConvertStrToDate(s.GetEnd())
	if err != nil {
		return nil, err
	}
	model.End = end

	return &model, nil
}

func (s *SpecialHolidayRepository) Find(id string) (*domains.SpecialHoliday, error) {
	model := SpecialHolidayModel{}

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

func (s *SpecialHolidayRepository) FindAll() ([]domains.SpecialHoliday, error) {
	models := []SpecialHolidayModel{}

	err := s.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.SpecialHoliday{}
	for _, model := range models {
		item, err := model.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}

func (s *SpecialHolidayRepository) Create(item *domains.SpecialHoliday) (string, error) {
	model, err := newSpecialHolidayModel(item)
	if err != nil {
		return "", err
	}
	err = s.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (s *SpecialHolidayRepository) Update(item *domains.SpecialHoliday) error {
	model, err := newSpecialHolidayModel(item)
	if err != nil {
		return err
	}
	err = s.db.Save(&model).Error
	return err
}

func (s *SpecialHolidayRepository) Delete(id string) error {
	model := SpecialHolidayModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	// delete physically (to reduce record)
	err := s.db.Unscoped().Delete(&model).Error
	return err
}
