package store

import "chico/takeout/common"

type BusinessHoursRepository interface {
	Fetch() (*BusinessHours, error)
	Update(target BusinessHours) error
}

type BusinessHours struct {
	morning BusinessHour
	lunch   BusinessHour
	dinner  BusinessHour
}

func NewDefaultBusinessHours() (*BusinessHours, error) {
	// create default
	morning, _ := NewBusinessHour("07:00", "10:00", []Weekday{Monday})
	lunch, _ := NewBusinessHour("11:00", "15:00", []Weekday{Monday})
	dinner, _ := NewBusinessHour("18:00", "21:00", []Weekday{Friday, Saturday})
	return NewBusinessHours(*morning, *lunch, *dinner)
}

func NewBusinessHours(morning, lunch, dinner BusinessHour) (*BusinessHours, error) {
	businessHours := &BusinessHours{
		morning: morning,
		lunch:   lunch,
		dinner:  dinner,
	}
	if err := businessHours.validateDuplicate(); err != nil {
		return nil, err
	}
	return businessHours, nil
}

func (b *BusinessHours) GetMorning() BusinessHour {
	return b.morning
}

func (b *BusinessHours) GetLunch() BusinessHour {
	return b.lunch
}

func (b *BusinessHours) GetDinner() BusinessHour {
	return b.dinner
}

func (b *BusinessHours) SetMorning(start, end string, weekdays []Weekday) error {
	morning, err := NewBusinessHour(start, end, weekdays)
	if err != nil {
		return err
	}
	before := b.morning
	b.morning = *morning
	if err := b.validateDuplicate(); err != nil {
		b.morning = before
		return err
	}
	return nil
}

func (b *BusinessHours) SetLunch(start, end string, weekdays []Weekday) error {
	lunch, err := NewBusinessHour(start, end, weekdays)
	if err != nil {
		return err
	}
	before := b.lunch
	b.lunch = *lunch
	if err := b.validateDuplicate(); err != nil {
		b.lunch = before
		return err
	}
	return nil
}

func (b *BusinessHours) SetDinner(start, end string, weekdays []Weekday) error {
	dinner, err := NewBusinessHour(start, end, weekdays)
	if err != nil {
		return err
	}
	before := b.dinner
	b.dinner = *dinner
	if err := b.validateDuplicate(); err != nil {
		b.dinner = before
		return err
	}
	return nil
}

func (b *BusinessHours) validateDuplicate() error {
	// check duplicate business hour
	if b.morning.IsOverlap(b.lunch) {
		return common.NewValidationError("business hour", "morning and lunch time is overlap")
	}
	if b.morning.IsOverlap(b.dinner) {
		return common.NewValidationError("business hour", "morning and dinner time is overlap")
	}
	if b.lunch.IsOverlap(b.dinner) {
		return common.NewValidationError("business hour", "lunch and dinner time is overlap")
	}
	return nil
}

type BusinessHour struct {
	shift    TimeRange
	weekdays []Weekday
}

func NewBusinessHour(start, end string, weekdays []Weekday) (*BusinessHour, error) {
	shift, err := NewTimeRange(start, end)
	if err != nil {
		return nil, err
	}

	// empty is allowed
	if weekdays == nil {
		return nil, common.NewValidationError("weekdays", "is nil")
	}
	if validateDuplicatedWeekdays(weekdays) {
		return nil, common.NewValidationError("weekdays", "duplicated value exists")
	}

	return &BusinessHour{
		shift:    *shift,
		weekdays: weekdays,
	}, nil
}

func validateDuplicatedWeekdays(weekdays []Weekday) bool {
	duplicated := false
	encountered := map[Weekday]bool{}
	for i := 0; i < len(weekdays); i++ {
		if !encountered[weekdays[i]] {
			encountered[weekdays[i]] = true
		} else {
			duplicated = true
			break
		}
	}
	return duplicated
}

func (b *BusinessHour) IsOverlap(other BusinessHour) bool {
	isOverlap := false
	for _, weekday := range b.weekdays {
		for _, targetWeekday := range other.weekdays {
			if weekday == targetWeekday {
				// check time in overlap
				if b.shift.IsOverlap(other.shift) {
					isOverlap = true
					break
				}
			}
		}
	}
	return isOverlap
}

func (b *BusinessHour) GetStart() string {
	return b.shift.start
}

func (b *BusinessHour) GetEnd() string {
	return b.shift.end
}

func (b *BusinessHour) GetWeekdays() []Weekday {
	return b.weekdays
}
