package store

import (
	"chico/takeout/common"
	"fmt"
)

type SpecialBusinessHourSpecification struct {
	specialHours []SpecialBusinessHour
}

func NewSpecialBusinessHourSpecification(specialHours []SpecialBusinessHour) *SpecialBusinessHourSpecification {
	return &SpecialBusinessHourSpecification{
		specialHours: specialHours,
	}
}

func (s *SpecialBusinessHourSpecification) Validate(item *SpecialBusinessHour) error {
	// check same hour id is assigned at samedate
	exists, err := s.businessHourIdIsAssigned(item)
	if err != nil {
		return err
	}
	if exists {
		return common.NewValidationError("businessHourId", fmt.Sprintf("business hour id at same date is already assigned. only 1 assign is allowed. id:%s, date=%s", item.GetBusinessHourId(), item.GetDate()))
	}

	// check time is overlapped at samedate
	overLapped, err := s.dateAndTimeIsOverLapped(item)
	if err != nil {
		return err
	}
	if overLapped {
		return common.NewValidationError("Date", fmt.Sprintf("time is overwrapped:date=%s, start=%s, end=%s", item.GetDate(), item.GetStart(), item.GetEnd()))
	}

	return nil
}

func (s *SpecialBusinessHourSpecification) businessHourIdIsAssigned(item *SpecialBusinessHour) (bool, error) {
	for _, hour := range s.specialHours {
		// skip self
		if item.Equals(hour) {
			continue
		}
		if item.IsSameDateAndHourId(hour) {
			return true, nil
		}
	}
	return false, nil
}

func (s *SpecialBusinessHourSpecification) dateAndTimeIsOverLapped(item *SpecialBusinessHour) (bool, error) {
	for _, hour := range s.specialHours {
		// skip self
		if item.Equals(hour) {
			continue
		}
		if item.IsOverlap(hour) {
			return true, nil
		}
	}
	return false, nil
}

type HolidaySpecification struct {
	normalSchedules  BusinessHours
	specialSchedules []SpecialBusinessHour
	specialHolidays  []SpecialHoliday
}

func NewHolidaySpecification(normalSchedules BusinessHours,
	specialSchedules []SpecialBusinessHour,
	specialHolidays []SpecialHoliday) *HolidaySpecification {
	return &HolidaySpecification{
		normalSchedules:  normalSchedules,
		specialSchedules: specialSchedules,
		specialHolidays:  specialHolidays,
	}
}

func (h *HolidaySpecification) IsStoreInBusiness(datetime string) (bool, error) {
	time, err := common.ConvertStrToDateTime(datetime)
	if err != nil {
		return false, err
	}
	// step1 :specific date is available

	// check special holiday (most high priority)
	for _, sh := range h.specialHolidays {
		if sh.shift.InRangeDate(*time) {
			// store is holiday. can not reserve
			return false, nil
		}
	}

	// get specific date schedule by day of week
	// first: special schedule
	// second: if not exists, check normal

	hasSameDate := false
	for _, ss := range h.specialSchedules {
		if ss.IsSameDate(*time) {
			hasSameDate = true
			// check time is overlap
			if ss.IsInRange(*time) {
				// can reserve
				return true, nil
			}
		}
	}
	// if same date's special schedule has, other normal schedule is ignored (can not reserve)
	if hasSameDate {
		return false, nil
	}

	// get normal schedules by day of week
	return h.normalSchedules.IsInBusiness(*time), nil
}

type BusinessHoursManagementSpecification struct {
	normalSchedules  BusinessHours
	specialSchedules []SpecialBusinessHour
	specialHolidays  []SpecialHoliday
}

type BusinessHourInfo struct {
	date   string
	hours  []HourInfo
}
type HourInfo struct {
	hourId string
	name   string
	start  string
	end    string
}

// func (b *BusinessHoursManagementSpecification) GetStoreBusinessHours(date string) (*BusinessHourInfo, error) {
// 	dateTime, err := common.ConvertStrToDateTime(date)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// at first check special holiday
// 	// check special holiday (most high priority)
// 	for _, sh := range b.specialHolidays {
// 		if sh.shift.InRangeDate(*dateTime) {
// 			// store is holiday. can not reserve
// 			return nil, nil
// 		}
// 	}

// 	// check special schedule


// 	// get week of day
// 	// weekDay := dateTime.Weekday()
// 	// get match day
// }
