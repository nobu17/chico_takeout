package store

import (
	domains "chico/takeout/domains/store"
)

type BusinessHoursModel struct {
	Schedules []BusinessHourModel
}

func newBusinessHoursModel(item *domains.BusinessHours) *BusinessHoursModel {
	schedules := []BusinessHourModel{}
	for _, schedule := range item.GetSchedules() {
		schedules = append(schedules, *newBusinessHourModel(schedule))
	}
	return &BusinessHoursModel{
		Schedules: schedules,
	}
}

type BusinessHoursUpdateModel struct {
	Id         string
	Name       string
	Start      string
	End        string
	Weekdays   []Weekday
	OffsetHour uint
}

type BusinessHoursEnabledUpdateModel struct {
	Id      string
	Enabled bool
}

type BusinessHourModel struct {
	Id         string
	Name       string
	Start      string
	End        string
	Weekdays   []Weekday
	Enabled    bool
	OffsetHour uint
}

func newBusinessHourModel(item domains.BusinessHour) *BusinessHourModel {
	weekdays := []Weekday{}
	for _, week := range item.GetWeekdays() {
		weekdays = append(weekdays, newWeekday(week))
	}

	return &BusinessHourModel{
		Id:         item.GetId(),
		Name:       item.GetName(),
		Start:      item.GetStart(),
		End:        item.GetEnd(),
		Weekdays:   weekdays,
		Enabled:    item.GetEnabled(),
		OffsetHour: item.GetHourOffset(),
	}
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func newWeekday(weekday domains.Weekday) Weekday {
	switch weekday {
	case domains.Sunday:
		return Sunday
	case domains.Monday:
		return Monday
	case domains.Tuesday:
		return Tuesday
	case domains.Wednesday:
		return Wednesday
	case domains.Thursday:
		return Thursday
	case domains.Friday:
		return Friday
	case domains.Saturday:
		return Saturday
	default:
		panic("can not convert")
	}
}

func toDomainWeekday(weekdays []Weekday) []domains.Weekday {
	converted := []domains.Weekday{}
	for _, weekday := range weekdays {
		switch weekday {
		case Sunday:
			converted = append(converted, domains.Sunday)
		case Monday:
			converted = append(converted, domains.Monday)
		case Tuesday:
			converted = append(converted, domains.Tuesday)
		case Wednesday:
			converted = append(converted, domains.Wednesday)
		case Thursday:
			converted = append(converted, domains.Thursday)
		case Friday:
			converted = append(converted, domains.Friday)
		case Saturday:
			converted = append(converted, domains.Saturday)
		default:
			panic("can not convert")
		}
	}
	return converted
}

type BusinessHoursUseCase interface {
	GetAll() (*BusinessHoursModel, error)
	Update(model *BusinessHoursUpdateModel) error
	UpdateEnabled(model *BusinessHoursEnabledUpdateModel) error
	InitIfNotExists() error
}

type businessHoursUseCase struct {
	businessHoursRepository domains.BusinessHoursRepository
	businessHoursService    domains.BusinessHoursService
}

func NewBusinessHoursUseCase(
	businessHoursRepository domains.BusinessHoursRepository,
	specialBusinessHourRepository domains.SpecialBusinessHourRepository) BusinessHoursUseCase {
	return &businessHoursUseCase{
		businessHoursRepository: businessHoursRepository,
		businessHoursService:    *domains.NewBusinessHoursService(businessHoursRepository),
	}
}

func (b *businessHoursUseCase) GetAll() (*BusinessHoursModel, error) {
	data, err := b.businessHoursService.FetchBusinessHours()
	if err != nil {
		return nil, err
	}
	return newBusinessHoursModel(data), nil
}

func (b *businessHoursUseCase) InitIfNotExists() error {
	err := b.businessHoursService.InitDataIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

func (b *businessHoursUseCase) Update(model *BusinessHoursUpdateModel) error {
	businessHours, err := b.businessHoursService.FetchBusinessHours()
	if err != nil {
		return err
	}

	new, err := businessHours.Update(model.Id, model.Name, model.Start, model.End, toDomainWeekday(model.Weekdays), model.OffsetHour)
	if err != nil {
		return err
	}

	err = b.businessHoursRepository.Update(new)
	if err != nil {
		return err
	}

	return nil
}

func (b *businessHoursUseCase) UpdateEnabled(model *BusinessHoursEnabledUpdateModel) error {
	businessHours, err := b.businessHoursService.FetchBusinessHours()
	if err != nil {
		return err
	}

	new, err := businessHours.UpdateEnabled(model.Id, model.Enabled)
	if err != nil {
		return err
	}

	err = b.businessHoursRepository.Update(new)
	if err != nil {
		return err
	}

	return nil
}
