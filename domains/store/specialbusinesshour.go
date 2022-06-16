package store

import (
	"fmt"
	"strings"
	"time"

	"chico/takeout/common"

	"github.com/google/uuid"
)

type SpecialBusinessHourRepository interface {
	Find(id string) (*SpecialBusinessHour, error)
	FindAll() ([]SpecialBusinessHour, error)
	Create(item *SpecialBusinessHour) (string, error)
	Update(item *SpecialBusinessHour) error
	Delete(id string) error
}

type SpecialBusinessHour struct {
	id             string
	name           string
	date           Date
	shift          TimeRange
	businessHourId string
}

const (
	SpecialBusinessHourMaxName = 30
)

func NewSpecialBusinessHour(name, date, start, end, businessHourId string) (*SpecialBusinessHour, error) {
	hour := &SpecialBusinessHour{id: uuid.NewString()}
	if err := hour.Set(name, date, start, end, businessHourId); err != nil {
		return nil, err
	}
	return hour, nil
}

func NewSpecialBusinessHourForOrm(id, name, date, start, end, businessHourId string) (*SpecialBusinessHour, error) {
	hour := &SpecialBusinessHour{id: id}
	if err := hour.Set(name, date, start, end, businessHourId); err != nil {
		return nil, err
	}
	return hour, nil
}

func (s *SpecialBusinessHour) Set(name, date, start, end, businessHourId string) error {
	if err := s.validateName(name); err != nil {
		return err
	}

	if err := s.validateBusinessHourId(businessHourId); err != nil {
		return err
	}

	dateVal, err := NewDate(date)
	if err != nil {
		return err
	}

	shift, err := NewTimeRange(start, end)
	if err != nil {
		return err
	}

	s.name = name
	s.date = *dateVal
	s.shift = *shift
	s.businessHourId = businessHourId

	return nil
}

func (s *SpecialBusinessHour) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return common.NewValidationError("name", "required")
	}

	if len(name) > SpecialBusinessHourMaxName {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", SpecialBusinessHourMaxName))
	}
	return nil
}

func (s *SpecialBusinessHour) validateBusinessHourId(businessHourId string) error {
	if strings.TrimSpace(businessHourId) == "" {
		return common.NewValidationError("businessHourId", "required")
	}
	return nil
}

func (b *SpecialBusinessHour) IsOverlap(other SpecialBusinessHour) bool {
	if b.HaveSameDate(other) {
		if b.shift.IsOverlap(other.shift) {
			return true
		}
	}
	return false
}

func (b *SpecialBusinessHour) Equals(other SpecialBusinessHour) bool {
	return b.id == other.GetId()
}

func (s *SpecialBusinessHour) GetId() string {
	return s.id
}

func (s *SpecialBusinessHour) GetName() string {
	return s.name
}

func (s *SpecialBusinessHour) GetDate() string {
	return s.date.GetValue()
}

func (s *SpecialBusinessHour) GetStart() string {
	return s.shift.start
}

func (s *SpecialBusinessHour) GetEnd() string {
	return s.shift.end
}

func (s *SpecialBusinessHour) GetBusinessHourId() string {
	return s.businessHourId
}

func (s *SpecialBusinessHour) HaveSameBusinessHourId(other SpecialBusinessHour) bool {
	return s.businessHourId == other.businessHourId
}

func (s *SpecialBusinessHour) HaveSameDate(other SpecialBusinessHour) bool {
	return s.date == other.date
}

func (s *SpecialBusinessHour) IsSameDateAndHourId(other SpecialBusinessHour) bool {
	return s.HaveSameBusinessHourId(other) && s.HaveSameDate(other)
}

func (s *SpecialBusinessHour) IsSameDate(datetime time.Time) bool {
	return s.date.IsSameDate(datetime)
}

func (s *SpecialBusinessHour) IsInRange(datetime time.Time) bool {
	// same date and within range time
	if s.date.IsSameDate(datetime) {
		return s.shift.IsInRange(datetime)
	}
	return false
}


