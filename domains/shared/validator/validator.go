package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
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
		return common.NewValidationError(s.name, fmt.Sprintf("MaxLength:%d", s.maxLength))
	}
	return nil
}

type AllowEmptyString struct {
	name      string
	maxLength int
}

func NewAllowEmptyStingLength(name string, maxLength int) *AllowEmptyString {
	return &AllowEmptyString{
		name:      name,
		maxLength: maxLength,
	}
}

func (s *AllowEmptyString) Validate(val string) error {
	if utf8.RuneCountInString(val) > s.maxLength {
		return common.NewValidationError(s.name, fmt.Sprintf("MaxLength:%d", s.maxLength))
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

type UrlValidator struct {
	name string
	allowEmpty bool
}

func NewUrlValidator(name string, allowEmpty bool) *UrlValidator {
	return &UrlValidator{
		name: name,
		allowEmpty: allowEmpty,
	}
}

func (u *UrlValidator) Validate(url string) error {
	if !u.allowEmpty && len(strings.TrimSpace(url)) == 0 {
		return common.NewValidationError(u.name, "Not allowed to be empty.") 
	}
	if u.allowEmpty && len(strings.TrimSpace(url)) == 0 {
		return nil
	}
	if !u.isUrl(url) {
		return common.NewValidationError(u.name, fmt.Sprintf("Incorrect url format:%s", url))	
	}

	return nil
}

func (u *UrlValidator) isUrl(str string) bool {
    url, err := url.Parse(str)
    return err == nil && url.Scheme != "" && url.Host != ""
}

type EmailValidator struct {
	name string
}

func NewEmailValidator(name string) *EmailValidator {
	return &EmailValidator{
		name: name,
	}
}

func (e *EmailValidator) Validate(email string) error {
	if len(strings.TrimSpace(email)) == 0 {
		return common.NewValidationError(e.name, "Not allowed to be empty.") 
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return common.NewValidationError(e.name, fmt.Sprintf("Not email style. %s. %s", err, email)) 	
	}

	return nil
}

type TelNoValidator struct {
	name string
}

func NewTelNoValidator(name string) *TelNoValidator {
	return &TelNoValidator{
		name: name,
	}
}

func (t *TelNoValidator) Validate(telNo string) error {
	r := regexp.MustCompile(`^[0-9]+$`)
	if !r.MatchString(telNo) {
		return common.NewValidationError(t.name, fmt.Sprintf("Not only number. %s.", telNo)) 
	}

	return nil
}