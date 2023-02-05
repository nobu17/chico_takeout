package shared

import (
	"chico/takeout/common"
	"fmt"
	"time"
)

type IntValue struct {
	value int
}

func NewIntValue(value int) IntValue {
	return IntValue{value: value}
}

func (i *IntValue) GetValue() int {
	return i.value
}

type StringValue struct {
	value string
}

func NewStringValue(value string) StringValue {
	return StringValue{value: value}
}

func (i *StringValue) GetValue() string {
	return i.value
}

type Date struct {
	StringValue
}

func NewDate(value string) (*Date, error) {
	_, err := common.ConvertStrToDate(value)
	if err != nil {
		return nil, common.NewValidationError("date", fmt.Sprintf("can not convert date:%s", value))
	}
	return &Date{StringValue: NewStringValue(value)}, nil
}

func (d *Date) IsSameDate(datetime time.Time) bool {
    dateStr := common.ConvertTimeToDateStr(datetime)
	return d.StringValue.GetValue() == dateStr
}

func (d *Date) GetAsDate() time.Time {
    v, err := common.ConvertStrToDate(d.value);
	if err != nil {
		fmt.Println(err)
		panic("should not be allowed failed convert");
	}
	return *v
}
