package message

import (
	"chico/takeout/domains/shared"
	"chico/takeout/domains/shared/validator"
)

const (
	MaxContentLength = 1000
)

type Content struct {
	shared.StringValue
}

func NewContent(value string) (*Content, error) {
	validator := validator.NewStingLength("Content", MaxContentLength)
	if err := validator.Validate(value); err != nil {
		return nil, err
	}

	return &Content{StringValue: shared.NewStringValue(value)}, nil
}