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

func (s *SpecialBusinessHourSpecification) Validate(item SpecialBusinessHour) error {
	exists, err := s.businessHourIdIsAssigned(item)
	if err != nil {
		return err
	}
	if !exists {
		return common.NewValidationError("businessHourId", fmt.Sprintf("business hour id is already assigned. only 1 assign is allowed:%s", item.GetBusinessHourId()))
	}

	overWrapped, err := s.dateIsOverwarped(item)
	if err != nil {
		return err
	}
	if !overWrapped {
		return common.NewValidationError("Date", fmt.Sprintf("date is overwrapped:%s", item.GetDate()))
	}

	return nil
}

func (s *SpecialBusinessHourSpecification) businessHourIdIsAssigned(item SpecialBusinessHour) (bool, error) {
	for _, hour := range s.specialHours {
		if item.HaveSameBusinessHourId(hour) {
			return true, nil
		}
	}
	return false, nil
}

func (s *SpecialBusinessHourSpecification) dateIsOverwarped(item SpecialBusinessHour) (bool, error) {
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
