package shared

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
