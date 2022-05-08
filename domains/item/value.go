package item

import (
	"fmt"

	"chico/takeout/common"
	"chico/takeout/domains/shared"
	"chico/takeout/domains/shared/validator"
)

type Priority struct {
	shared.IntValue
}

var priorityValidator = validator.NewPlusInteger("Priority")

func NewPriority(value int) (*Priority, error) {
	if err := priorityValidator.Validate(value); err != nil {
		return nil, err
	}

	return &Priority{IntValue: shared.NewIntValue(value)}, nil
}

var maxOrderValidator = validator.NewRangeInteger("MaxOrder", 1, MaxOrderMaxValue)

type MaxOrder struct {
	shared.IntValue
}

const (
	MaxOrderMaxValue = 30
)

func NewMaxOrder(value int) (*MaxOrder, error) {
	if err := maxOrderValidator.Validate(value); err != nil {
		return nil, err
	}

	return &MaxOrder{IntValue: shared.NewIntValue(value)}, nil
}

func (m *MaxOrder) WithinLimit(request int) error {
	if m.GetValue() < request {
		return common.NewValidationError("maxOrder", fmt.Sprintf("Need to be less than. max:%d, request:%d", m.GetValue(), request))
	}
	return nil
}

type Price struct {
	shared.IntValue
}

func NewPrice(value, maxValue int) (*Price, error) {
	validator := validator.NewRangeInteger("Price", 1, maxValue)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Price{IntValue: shared.NewIntValue(value)}, nil
}

type StockRemain struct {
	maxValue int
	shared.IntValue
}

func NewStockRemain(value, maxValue int) (*StockRemain, error) {
	validator := validator.NewRangeInteger("StockRemain", 0, maxValue)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &StockRemain{IntValue: shared.NewIntValue(value), maxValue: maxValue}, nil
}

func (p *StockRemain) Consume(request int) (*StockRemain, error) {
	if request < 1 {
		return nil, common.NewValidationError("stock remain", fmt.Sprintf("request is needed more than 1. request:%d", request))
	}
	remain := p.GetValue() - request
	if remain < 0 {
		return nil, common.NewValidationError("stock remain", fmt.Sprintf("remain count is insufficient. remain:%d, request:%d", p.GetValue(), request))
	}
	return NewStockRemain(remain, p.maxValue)
}

func (p *StockRemain) Increase(request int) (*StockRemain, error) {
	if request < 1 {
		return nil, common.NewValidationError("stock remain", fmt.Sprintf("request is needed more than 1. request:%d", request))
	}
	remain := p.GetValue() + request
	return NewStockRemain(remain, p.maxValue)
}

const (
	MaxOrderPerDayMaxValue = 30
)

type MaxOrderPerDay struct {
	shared.IntValue
}

var maxOrderPerDayValidator = validator.NewRangeInteger("MaxOrderPerDay", 1, MaxOrderPerDayMaxValue)

func NewMaxOrderPerDay(value int) (*MaxOrderPerDay, error) {
	if err := maxOrderPerDayValidator.Validate(value); err != nil {
		return nil, err
	}

	return &MaxOrderPerDay{IntValue: shared.NewIntValue(value)}, nil
}

type Descritpion struct {
	shared.StringValue
}

func NewDescritpion(value string, maxLength int) (*Descritpion, error) {
	validator := validator.NewStingLength("Descritpion", maxLength)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Descritpion{StringValue: shared.NewStringValue(value)}, nil
}

type Name struct {
	shared.StringValue
}
func NewName(value string, maxLength int) (*Name, error) {
	validator := validator.NewStingLength("Name", maxLength)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Name{StringValue: shared.NewStringValue(value)}, nil
}
