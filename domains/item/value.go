package item

import (
	"fmt"

	"chico/takeout/common"
)

type Priority struct {
	value int
}

func NewPriority(value int) (*Priority, error) {
	if err := validatePriorityValue(value); err != nil {
		return nil, err
	}

	return &Priority{value: value}, nil
}

func validatePriorityValue(priority int) error {
	if priority < 1 {
		return common.NewValidationError("priority", "Need to be greater than 1")
	}
	return nil
}

func (p *Priority) GetValue() int {
	return p.value
}

type MaxOrder struct {
	value int
}

const (
	MaxOrderMaxValue = 30
)

func NewMaxOrder(value int) (*MaxOrder, error) {
	if err := validateMaxOrder(value); err != nil {
		return nil, err
	}

	return &MaxOrder{value: value}, nil
}

func validateMaxOrder(maxOrder int) error {
	if maxOrder < 1 {
		return common.NewValidationError("maxOrder", "Need to be greater than 1")
	}
	if maxOrder > MaxOrderMaxValue {
		return common.NewValidationError("maxOrder", fmt.Sprintf("Need to be less than %d", MaxOrderMaxValue))
	}
	return nil
}

func (m *MaxOrder) GetValue() int {
	return m.value
}

func (m *MaxOrder) WithinLimit(request int) error {
	if m.value < request {
		return common.NewValidationError("maxOrder", fmt.Sprintf("Need to be less than. max:%d, request:%d", m.value, request))
	}
	return nil
}

type Price struct {
	maxValue int
	value    int
}

func NewPrice(value, maxValue int) (*Price, error) {
	if err := validatePrice(value, maxValue); err != nil {
		return nil, err
	}

	return &Price{value: value, maxValue: maxValue}, nil
}

func validatePrice(price, maxValue int) error {
	if price < 1 {
		return common.NewValidationError("price", "Need to be greater than 1")
	}
	if price > maxValue {
		return common.NewValidationError("price", fmt.Sprintf("Need to be less than %d", maxValue))
	}
	return nil
}

func (p *Price) GetValue() int {
	return p.value
}

type StockRemain struct {
	maxValue int
	value    int
}

func NewStockRemain(value, maxValue int) (*StockRemain, error) {
	if err := validateStockRemain(value, maxValue); err != nil {
		return nil, err
	}

	return &StockRemain{value: value, maxValue: maxValue}, nil
}

func validateStockRemain(value, maxValue int) error {
	if value < 0 {
		return common.NewValidationError("stock remain", "Need to be greater than 1")
	}
	if value > maxValue {
		return common.NewValidationError("stock remain", fmt.Sprintf("Need to be less than %d", maxValue))
	}
	return nil
}

func (p *StockRemain) GetValue() int {
	return p.value
}

func (p *StockRemain) Consume(request int) (*StockRemain, error) {
	if request < 1 {
		return nil, common.NewValidationError("stock remain", fmt.Sprintf("request is needed more than 1. request:%d", request))
	}
	remain := p.value - request
	if remain < 0 {
		return nil, common.NewValidationError("stock remain", fmt.Sprintf("remain count is insufficient. remain:%d, request:%d", p.value, request))
	}
	return &StockRemain{remain, p.maxValue}, nil
}
