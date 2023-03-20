package store

import (
	"chico/takeout/common"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BusinessHoursMaxSchedules = 3
)

type BusinessHoursRepository interface {
	Fetch() (*BusinessHours, error)
	Update(target *BusinessHours) error
	Create(target *BusinessHours) error
}

type BusinessHours struct {
	schedules []BusinessHour
}

func NewDefaultBusinessHours() (*BusinessHours, error) {
	// create default
	morning, _ := NewBusinessHour("morning", "07:00", "09:30", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	lunch, _ := NewBusinessHour("lunch", "11:30", "15:00", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	dinner, _ := NewBusinessHour("dinner", "18:00", "21:00", []Weekday{Wednesday, Saturday})

	schedules := []BusinessHour{}
	schedules = append(schedules, *morning)
	schedules = append(schedules, *lunch)
	schedules = append(schedules, *dinner)
	return NewBusinessHours(schedules)
}

func NewBusinessHours(schedules []BusinessHour) (*BusinessHours, error) {
	businessHours := &BusinessHours{
		schedules: schedules,
	}
	if err := businessHours.validateSchedules(); err != nil {
		return nil, err
	}
	return businessHours, nil
}

// func (b *BusinessHours) GetSchedule(id string) (BusinessHour, error) {
// 	_, item := b.findSchedule(id)
// 	if item == nil {
// 		return BusinessHour{}, common.NewNotFoundError("not found item")
// 	}
// 	return *item, nil
// }

func (b *BusinessHours) GetSchedules() []BusinessHour {
	// copy for immutable
	tmp := append([]BusinessHour{}, b.schedules...)
	return tmp
}

func (b *BusinessHours) FindById(id string) *BusinessHour {
	_, item := b.findSchedule(id)
	if item == nil {
		return nil
	}
	// return copy for immutable
	newItem := item
	return newItem
}

func (b *BusinessHours) IsInBusiness(targetDateTime time.Time) bool {
	// check each schedule
	for _, bs := range b.schedules {
		if !bs.enabled {
			continue
		}
		isIn := bs.IsInSchedule(targetDateTime)
		if isIn {
			return true
		}
	}
	return false
}

// currently add is not needed.
// func (b *BusinessHours) Add(name, start, end string, weekdays []Weekday) error {
// 	hour, err := NewBusinessHour(name, start, end, weekdays)
// 	if err != nil {
// 		return err
// 	}
// 	tmp := append([]BusinessHour{}, b.schedules...)
// 	tmp = append(tmp, *hour)
// 	// check duplicate
// 	_, err = NewBusinessHours(tmp)
// 	if err != nil {
// 		return err
// 	}

// 	b.schedules = append(b.schedules, *hour)
// 	return nil
// }

func (b *BusinessHours) Update(id, name, start, end string, weekdays []Weekday) (*BusinessHours, error) {
	selfCopy, err := b.Copy()
	if err != nil {
		return nil, fmt.Errorf("unexpected error at copy:%s", err)
	}

	_, target := selfCopy.findSchedule(id)
	if target == nil {
		return nil, common.NewNotFoundError("id")
	}
	err = target.Set(name, start, end, weekdays)
	if err != nil {
		return nil, err
	}
	// check overWrap
	err = selfCopy.validateSchedules()
	if err != nil {
		return nil, err
	}
	return selfCopy, nil
}

func (b *BusinessHours) UpdateEnabled(id string, enabled bool) (*BusinessHours, error) {
	selfCopy, err := b.Copy()
	if err != nil {
		return nil, fmt.Errorf("unexpected error at copy:%s", err)
	}

	_, target := selfCopy.findSchedule(id)
	if target == nil {
		return nil, common.NewNotFoundError("id")
	}
	target.SetEnabled(enabled)
	return selfCopy, nil
}

func (b *BusinessHours) findSchedule(id string) (int, *BusinessHour) {
	// return pointer, so it is not immutable
	for i := 0; i < len(b.schedules); i++ {
		if b.schedules[i].id == id {
			return i, &b.schedules[i]
		}
	}
	return -1, nil
}

func (b *BusinessHours) validateSchedules() error {
	if len(b.schedules) == 0 {
		return common.NewValidationError("schedules", "is empty or nil.")
	}

	if len(b.schedules) > BusinessHoursMaxSchedules {
		return common.NewValidationError("schedules", fmt.Sprintf("MaxLength:%d", BusinessHoursMaxSchedules))
	}
	// check duplicate
	return b.validateDuplicate()
}

func (b *BusinessHours) validateDuplicate() error {
	// check duplicate business hour
	// this case not check enabled or disabled
	// even if disabled, duplication is not allowed
	for i, schedule := range b.schedules {
		for j := i + 1; j < len(b.schedules); j++ {
			target := b.schedules[j]
			if schedule.IsOverlap(target) {
				return common.NewValidationError("business hour", fmt.Sprintf("%s and %s time is overlap", schedule.name, target.name))
			}
		}
	}
	return nil
}

func (b *BusinessHours) Copy() (*BusinessHours, error) {
	hours := []BusinessHour{}
	for _, sc := range b.schedules {
		hours = append(hours, *sc.Copy())
	}
	return NewBusinessHours(hours)
}

const (
	BusinessHourNameMaxLength = 10
)

type BusinessHour struct {
	id       string
	name     Name
	shift    TimeRange
	weekdays []Weekday
	enabled     bool
}

func NewBusinessHour(name, start, end string, weekdays []Weekday) (*BusinessHour, error) {
	businessHour := &BusinessHour{id: uuid.NewString()}

	err := businessHour.Set(name, start, end, weekdays)
	if err != nil {
		return nil, err
	}
	// default enabled
	businessHour.SetEnabled(true)
	return businessHour, nil
}

func NewBusinessHourForOrm(id, name, start, end string, weekdays []Weekday, enabled bool) (*BusinessHour, error) {
	businessHour := &BusinessHour{id: id}

	err := businessHour.Set(name, start, end, weekdays)
	if err != nil {
		return nil, err
	}
	businessHour.SetEnabled(enabled)
	return businessHour, nil
}

func (b *BusinessHour) Copy() *BusinessHour {
	businessHour, _ := NewBusinessHour(b.name.GetValue(), b.shift.start, b.shift.end, b.weekdays)
	businessHour.id = b.id
	return businessHour
}

func (b *BusinessHour) Equals(other BusinessHour) bool {
	return b.id == other.id
}

func (b *BusinessHour) Set(name, start, end string, weekdays []Weekday) error {
	nameVal, err := NewName(name, BusinessHourNameMaxLength)
	if err != nil {
		return err
	}

	shift, err := NewTimeRange(start, end)
	if err != nil {
		return err
	}

	// empty is allowed
	if weekdays == nil {
		return common.NewValidationError("weekdays", "is nil")
	}
	if validateDuplicatedWeekdays(weekdays) {
		return common.NewValidationError("weekdays", "duplicated value exists")
	}

	b.name = *nameVal
	b.shift = *shift
	b.weekdays = weekdays

	return nil
}

func (b *BusinessHour) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func validateBusinessHourName(name string) error {
	if strings.TrimSpace(name) == "" {
		return common.NewValidationError("name", "required")
	}

	if len(name) > BusinessHourNameMaxLength {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", BusinessHourNameMaxLength))
	}
	return nil
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

func (b *BusinessHour) IsInSchedule(targetDateTime time.Time) bool {
	targetWk := targetDateTime.Weekday()
	// check week day
	for _, wk := range b.weekdays {
		if int(wk) == int(targetWk) {
			// check time shift is in range
			return b.shift.IsInRange(targetDateTime)
		}
	}
	return false
}

func (b *BusinessHour) GetId() string {
	return b.id
}

func (b *BusinessHour) GetName() string {
	return b.name.GetValue()
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

func (b *BusinessHour) GetEnabled() bool {
	return b.enabled
}