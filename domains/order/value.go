package order

import (
	"fmt"
	"time"

	"chico/takeout/common"
	"chico/takeout/domains/shared"
	"chico/takeout/domains/shared/validator"
)

type Price struct {
	value int
}

func NewPrice(value int) (*Price, error) {
	if err := validatePrice(value); err != nil {
		return nil, err
	}

	return &Price{value: value}, nil
}

func validatePrice(price int) error {
	if price < 1 {
		return common.NewValidationError("price", "Need to be greater than 1")
	}
	return nil
}

type Quantity struct {
	value int
}

func NewQuantity(value int) (*Quantity, error) {
	if err := validateQuantity(value); err != nil {
		return nil, err
	}

	return &Quantity{value: value}, nil
}

func validateQuantity(value int) error {
	if value < 1 {
		return common.NewValidationError("quantity", "Need to be greater than 1")
	}
	return nil
}

type DateTime struct {
	value string
	time  time.Time
}

func NewDateTime(value string) (*DateTime, error) {
	t, err := common.ConvertStrToDateTime(value)
	if err != nil {
		return nil, common.NewValidationError("date", fmt.Sprintf("can not convert dateTime:%s", value))
	}
	return &DateTime{value: value, time: *t}, nil
}

func (d *DateTime) GetAsDate() string {
	val, _ := common.ConvertDateTimeStrToDateStr(d.value)
	return val
}

func (d *DateTime) GetDateTime() time.Time {
	return d.time
}

// for testing
var now = time.Now

func MockNow(mock func() time.Time) {
	now = mock
}
func ResetNow() {
	now = time.Now
}

type OrderDateTime struct {
	DateTime
}

func NewOrderDateTime() (*OrderDateTime, error) {
	str := common.ConvertTimeToDateTimeStr(now())
	item, err := NewDateTime(str)
	if err != nil {
		return nil, err
	}
	return &OrderDateTime{DateTime: *item}, nil
}

type PickupDateTime struct {
	DateTime
}

func NewPickupDateTime(value string) (*PickupDateTime, error) {
	item, err := NewDateTime(value)
	if err != nil {
		return nil, err
	}
	date, _ := common.ConvertStrToDateTime(value)
	// pick up time should be future from now + offset
	now := now()
	if !common.StartIsBeforeEnd(now, *date, common.OffsetMinutesOrderableBefore) {
		return nil, common.NewValidationError("PickupDateTime", fmt.Sprintf("not allowed time(%s). limit: now(%s) + minutes:(%v)", value, now, common.OffsetMinutesOrderableBefore))
	}
	return &PickupDateTime{DateTime: *item}, nil
}

type Memo struct {
	shared.StringValue
}

func NewMemo(value string, maxLength int) (*Memo, error) {
	validator := validator.NewAllowEmptyStingLength("Memo", maxLength)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Memo{StringValue: shared.NewStringValue(value)}, nil
}

type Email struct {
	shared.StringValue
}

func NewEmail(value string) (*Email, error) {
	validator := validator.NewEmailValidator("Email")
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Email{StringValue: shared.NewStringValue(value)}, nil
}

type TelNo struct {
	shared.StringValue
}

func NewTelNo(value string) (*TelNo, error) {
	validator := validator.NewTelNoValidator("TelNo")
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &TelNo{StringValue: shared.NewStringValue(value)}, nil
}

type UserName struct {
	shared.StringValue
}

func NewUserName(value string, maxLength int) (*UserName, error) {
	validator := validator.NewStingLength("UserName", maxLength)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &UserName{StringValue: shared.NewStringValue(value)}, nil
}
