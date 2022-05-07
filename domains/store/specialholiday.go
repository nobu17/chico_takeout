package store

import (
	"fmt"
	"strings"

	"chico/takeout/common"

	"github.com/google/uuid"
)

type SpecialHolidayRepository interface {
	Find(id string) (*SpecialHoliday, error)
	FindAll() ([]SpecialHoliday, error)
	Create(item *SpecialHoliday) (string, error)
	Update(item *SpecialHoliday) error
	Delete(id string) error
}

type SpecialHoliday struct {
	id    string
	name  string
	shift DateRange
}

const (
	SpecialHolidayMaxName = 20
)

func NewSpecialHoliday(name, start, end string) (*SpecialHoliday, error) {
	holiday := &SpecialHoliday{
		id: uuid.NewString(),
	}
	err := holiday.Set(name, start, end)
	if err != nil {
		return nil, err
	}
	return holiday, nil
}

func (h *SpecialHoliday) validateHolidayInfoName(name string) error {
	if strings.TrimSpace(name) == "" {
		return common.NewValidationError("name", "required")
	}

	if len(name) > SpecialHolidayMaxName {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", SpecialHolidayMaxName))
	}
	return nil
}

func (h *SpecialHoliday) Set(name, start, end string) error {
	if err := h.validateHolidayInfoName(name); err != nil {
		return err
	}
	shift, err := NewDateRange(start, end)
	if err != nil {
		return err
	}

	h.name = name
	h.shift = *shift
	return nil
}

func (s *SpecialHoliday) GetId() string {
	return s.id
}

func (s *SpecialHoliday) GetName() string {
	return s.name
}

func (s *SpecialHoliday) GetStart() string {
	return s.shift.GetStart()
}

func (s *SpecialHoliday) GetEnd() string {
	return s.shift.GetEnd()
}

func (s *SpecialHoliday) Equals(other SpecialHoliday) bool {
	return s.GetId() == other.GetId()
}

func (s *SpecialHoliday) IsOverlap(other SpecialHoliday) bool {
	return s.shift.IsOverlap(other.shift)
}
