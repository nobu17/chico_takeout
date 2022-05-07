package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"chico/takeout/common"
)

type StingLength struct {
	name      string
	maxLength int
}

func NewStingLength(name string, maxLength int) *StingLength {
	return &StingLength{
		name:      name,
		maxLength: maxLength,
	}
}

func (s *StingLength) Validate(val string) error {
	if strings.TrimSpace(val) == "" {
		return common.NewValidationError(s.name, "required")
	}

	if utf8.RuneCountInString(val) > s.maxLength {
		return common.NewValidationError("name", fmt.Sprintf("MaxLength:%d", s.maxLength))
	}
	return nil
}

type IntValidator interface {
	Validate(value int) error
}

type PlusInteger struct {
	name string
}

func NewPlusInteger(name string) *PlusInteger {
	return &PlusInteger{
		name: name,
	}
}

func (p *PlusInteger) Validate(val int) error {
	if val < 1 {
		return common.NewValidationError(p.name, "Need to be greater than 1")
	}
	return nil
}

type RangeInteger struct {
	name  string
	start int
	end   int
}

func NewRangeInteger(name string, start, end int) *RangeInteger {
	if start > end {
		panic("not allowed range")
	}
	return &RangeInteger{
		name:  name,
		start: start,
		end:   end,
	}
}

func (r *RangeInteger) Validate(val int) error {
	if val < r.start {
		return common.NewValidationError(r.name, fmt.Sprintf("Need to be greater than %d", r.start))
	}
	if val > r.end {
		return common.NewValidationError(r.name, fmt.Sprintf("Need to be less than %d", r.end))
	}
	return nil
}
