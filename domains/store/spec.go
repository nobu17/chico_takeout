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
