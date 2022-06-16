package store

import (
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/store"
	"chico/takeout/infrastructures/rdbms"

	"gorm.io/gorm"
)

type SpecialBusinessHoursRepository struct {
	db *gorm.DB
}

func NewSpecialBusinessHoursRepository(db *gorm.DB) *SpecialBusinessHoursRepository {
	return &SpecialBusinessHoursRepository{
		db: db,
	}
}

// `SpecialBusinessHourModel` belongs to `BusinessHourModel`, `BusinessHourModelID` is the foreign key
type SpecialBusinessHourModel struct {
	rdbms.BaseModel
	Name                string
	Date                *time.Time
	Start               *time.Time
	End                 *time.Time
	BusinessHourModelID string
	BusinessHourModel   BusinessHourModel
}

func (s *SpecialBusinessHourModel) toDomain() (*domains.SpecialBusinessHour, error) {
	startStr := common.ConvertTimeToTimeStr(*s.Start)
	endStr := common.ConvertTimeToTimeStr(*s.End)
	dateStr := common.ConvertTimeToDateStr(*s.Date)
	model, err := domains.NewSpecialBusinessHourForOrm(s.ID, s.Name, dateStr, startStr, endStr, s.BusinessHourModelID)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func newSpecialBusinessHourModel(s *domains.SpecialBusinessHour) (*SpecialBusinessHourModel, error) {
	model := SpecialBusinessHourModel{}
	model.ID = s.GetId()
	model.Name = s.GetName()

	date, err := common.ConvertStrToDate(s.GetDate())
	if err != nil {
		return nil, err
	}
	model.Date = date

	start, err := common.ConvertStrToTime(s.GetStart())
	if err != nil {
		return nil, err
	}
	model.Start = start

	end, err := common.ConvertStrToTime(s.GetEnd())
	if err != nil {
		return nil, err
	}
	model.End = end

	model.BusinessHourModelID = s.GetBusinessHourId()
	return &model, nil
}

func (s *SpecialBusinessHoursRepository) Find(id string) (*domains.SpecialBusinessHour, error) {
	model := SpecialBusinessHourModel{}

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

func (s *SpecialBusinessHoursRepository) FindAll() ([]domains.SpecialBusinessHour, error) {
	models := []SpecialBusinessHourModel{}

	err := s.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := []domains.SpecialBusinessHour{}
	for _, model := range models {
		item, err := model.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}

func (s *SpecialBusinessHoursRepository) Create(item *domains.SpecialBusinessHour) (string, error) {
	model, err := newSpecialBusinessHourModel(item)
	if err != nil {
		return "", err
	}
	err = s.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return item.GetId(), nil
}

func (s *SpecialBusinessHoursRepository) Update(item *domains.SpecialBusinessHour) error {
	model, err := newSpecialBusinessHourModel(item)
	if err != nil {
		return err
	}
	err = s.db.Save(&model).Error
	return err
}

func (s *SpecialBusinessHoursRepository) Delete(id string) error {
	model := SpecialBusinessHourModel{
		BaseModel: rdbms.BaseModel{ID: id},
	}
	err := s.db.Delete(&model).Error
	return err
}
