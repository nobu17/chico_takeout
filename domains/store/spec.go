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
	// check same hour id is assigned at same date
	exists, err := s.businessHourIdIsAssigned(item)
	if err != nil {
		return err
	}
	if exists {
		return common.NewValidationError("businessHourId", fmt.Sprintf("business hour id at same date is already assigned. only 1 assign is allowed. id:%s, date=%s", item.GetBusinessHourId(), item.GetDate()))
	}

	// check time is overlapped at same date
	overLapped, err := s.dateAndTimeIsOverLapped(item)
	if err != nil {
		return err
	}
	if overLapped {
		return common.NewValidationError("Date", fmt.Sprintf("time is duplicated:date=%s, start=%s, end=%s", item.GetDate(), item.GetStart(), item.GetEnd()))
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
	Date  string
	Hours []HourInfo
}
type HourInfo struct {
	HourTypeId string
	Name       string
	StartTime  string
	EndTime    string
}

func NewBusinessHoursManagementSpecification(normalSchedules BusinessHours,
	specialSchedules []SpecialBusinessHour,
	specialHolidays []SpecialHoliday) *BusinessHoursManagementSpecification {
	return &BusinessHoursManagementSpecification{
		normalSchedules,
		specialSchedules,
		specialHolidays,
	}
}

func (b *BusinessHoursManagementSpecification) GetStoreBusinessHours(dateStr string) (*BusinessHourInfo, error) {
	date, err := common.ConvertStrToDate(dateStr)
	if err != nil {
		return nil, err
	}
	result := BusinessHourInfo{Date: dateStr, Hours: []HourInfo{}}
	// at first check special holiday
	// check special holiday (most high priority)
	for _, sh := range b.specialHolidays {
		if sh.shift.InRangeDate(*date) {
			// store is holiday. no information
			return &result, nil
		}
	}

	// check special schedule
	hasSpecialSchedules := false
	for _, sc := range b.specialSchedules {
		if sc.IsSameDate(*date) {
			hour := HourInfo{
				HourTypeId: sc.GetBusinessHourId(),
				Name:       sc.GetName(),
				StartTime:  sc.GetStart(),
				EndTime:    sc.GetEnd(),
			}
			result.Hours = append(result.Hours, hour)

			hasSpecialSchedules = true
		}
	}

	// if already has, skip
	if hasSpecialSchedules {
		return &result, nil
	}

	// get normal schedule
	weekday := date.Weekday()
	for _, sc := range b.normalSchedules.FindByWeekday(int(weekday)) {
		hour := HourInfo{
			HourTypeId: sc.GetId(),
			Name:       sc.GetName(),
			StartTime:  sc.GetStart(),
			EndTime:    sc.GetEnd(),
		}
		result.Hours = append(result.Hours, hour)
	}
	return &result, nil
}
