package store

import (
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/store"
	"chico/takeout/infrastructures/rdbms"

	"gorm.io/gorm"
)

type BusinessHoursRepository struct {
	db *gorm.DB
}

func NewBusinessHoursRepository(db *gorm.DB) *BusinessHoursRepository {
	return &BusinessHoursRepository{
		db: db,
	}
}

// BusinessHourModel has many WeekDaysModel, BusinessHourModelID is the foreign key
type BusinessHourModel struct {
	rdbms.BaseModel
	Name     string
	Start    *time.Time
	End      *time.Time
	Weekdays []WeekDaysModel
}

func (b *BusinessHourModel) HasWeekDay(weekday int) bool {
	for _, week := range b.Weekdays {
		if week.Value == weekday {
			return true
		}
	}
	return false
}

func newBusinessHourModel(b *domains.BusinessHour) (*BusinessHourModel, error) {
	model := BusinessHourModel{}
	model.ID = b.GetId()
	model.Name = b.GetName()
	start, err := common.ConvertStrToTime(b.GetStart())
	if err != nil {
		return nil, err
	}
	model.Start = start

	end, err := common.ConvertStrToTime(b.GetEnd())
	if err != nil {
		return nil, err
	}
	model.End = end

	weekdays := []WeekDaysModel{}
	for _, weekday := range b.GetWeekdays() {
		week := WeekDaysModel{BusinessHourModelID: b.GetId(), Value: int(weekday)}
		weekdays = append(weekdays, week)
	}
	model.Weekdays = weekdays

	return &model, nil
}

func (b *BusinessHourModel) toDomain() (*domains.BusinessHour, error) {
	startStr := common.ConvertTimeToTimeStr(*b.Start)
	endStr := common.ConvertTimeToTimeStr(*b.End)
	weekdays := []domains.Weekday{}
	for _, weekday := range b.Weekdays {
		val := domains.Weekday(weekday.Value)
		weekdays = append(weekdays, val)
	}
	model, err := domains.NewBusinessHourForOrm(b.ID, b.Name, startStr, endStr, weekdays)

	if err != nil {
		return nil, err
	}
	return model, nil
}

type WeekDaysModel struct {
	BusinessHourModelID string `gorm:"primaryKey"`
	Value               int    `gorm:"primaryKey"`
}

func (b *BusinessHoursRepository) Fetch() (*domains.BusinessHours, error) {
	models := []BusinessHourModel{}

	err := b.db.Preload("Weekdays").Find(&models).Error
	if err != nil {
		return nil, err
	}

	hours := []domains.BusinessHour{}
	for _, model := range models {
		hour, err := model.toDomain()
		if err != nil {
			return nil, err
		}
		hours = append(hours, *hour)
	}

	// if no record return nil (will create default values)
	if len(models) == 0 {
		return nil, nil
	}

	schedules, err := domains.NewBusinessHours(hours)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (b *BusinessHoursRepository) Update(target *domains.BusinessHours) error {

	models := []BusinessHourModel{}
	for _, schedule := range target.GetSchedules() {
		model, err := newBusinessHourModel(&schedule)
		if err != nil {
			return err
		}
		models = append(models, *model)
	}

	var gError error = nil
	b.db.Transaction(func(tx *gorm.DB) error {
		// at first delete all week days
		err := b.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&WeekDaysModel{}).Error
		if err != nil {
			gError = err
			return err
		}
		for _, model := range models {
			//err := b.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&model).Error
			err = b.db.Save(&model).Error
			if err != nil {
				gError = err
				return err
			}
		}
		return nil
	})

	return gError
}

func (b *BusinessHoursRepository) Create(target *domains.BusinessHours) error {

	models := []BusinessHourModel{}
	for _, schedule := range target.GetSchedules() {
		model, err := newBusinessHourModel(&schedule)
		if err != nil {
			return err
		}
		models = append(models, *model)
	}

	var gError error = nil
	b.db.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			err := b.db.Create(&model).Error
			if err != nil {
				gError = err
				return err
			}
		}
		return nil
	})

	return gError
}
