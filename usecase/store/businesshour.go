package store

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/store"
)

type BusinessHoursModel struct {
	Morning BusinessHourModel
	Lunch   BusinessHourModel
	Dinner  BusinessHourModel
}

func newBusinessHoursModel(item *domains.BusinessHours) *BusinessHoursModel {
	return &BusinessHoursModel{
		Morning: *newBusinessHourModel(item.GetMorning()),
		Lunch:   *newBusinessHourModel(item.GetLunch()),
		Dinner:  *newBusinessHourModel(item.GetDinner()),
	}
}

type BusinessHoursUpdateModel struct {
	Morning *BusinessHourModel
	Lunch   *BusinessHourModel
	Dinner  *BusinessHourModel
}

type BusinessHourModel struct {
	Start    string
	End      string
	Weekdays []Weekday
}

func newBusinessHourModel(item domains.BusinessHour) *BusinessHourModel {
	weekdays := []Weekday{}
	for _, week := range item.GetWeekdays() {
		weekdays = append(weekdays, newWeekday(week))
	}

	return &BusinessHourModel{
		Start:    item.GetStart(),
		End:      item.GetEnd(),
		Weekdays: weekdays,
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

type BusinessHoursUseCase struct {
	businessHoursRepository domains.BusinessHoursRepository
	storeService            domains.StoreService
}

func NewBusinessHoursUseCase(businessHoursRepository domains.BusinessHoursRepository) *BusinessHoursUseCase {
	return &BusinessHoursUseCase{
		businessHoursRepository: businessHoursRepository,
		storeService:            *domains.NewStoreService(businessHoursRepository),
	}
}

func (b *BusinessHoursUseCase) Fetch() (*BusinessHoursModel, error) {
	data, err := b.storeService.FetchBusinessHours()
	if err != nil {
		return nil, err
	}
	return newBusinessHoursModel(data), nil
}

func (b *BusinessHoursUseCase) Update(model BusinessHoursUpdateModel) error {
	if model.Morning == nil && model.Lunch == nil && model.Dinner == nil {
		return common.NewValidationError("morning or lunch or dinner", "no update target existed")
	}

	businessHours, err := b.storeService.FetchBusinessHours()
	if err != nil {
		return err
	}

	if model.Morning != nil {
		err = businessHours.SetMorning(model.Morning.Start, model.Morning.End, toDomainWeekday(model.Morning.Weekdays))
		if err != nil {
			return err
		}
	}

	if model.Lunch != nil {
		err = businessHours.SetLunch(model.Lunch.Start, model.Lunch.End, toDomainWeekday(model.Lunch.Weekdays))
		if err != nil {
			return err
		}
	}

	if model.Dinner != nil {
		err = businessHours.SetDinner(model.Dinner.Start, model.Dinner.End, toDomainWeekday(model.Dinner.Weekdays))
		if err != nil {
			return err
		}
	}

	return nil
}
