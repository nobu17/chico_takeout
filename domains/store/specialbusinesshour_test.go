package store_test

import (
	"fmt"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/store"

	"github.com/stretchr/testify/assert"
)

type specialBusinessHourArgs struct {
	name           string
	date           string
	start          string
	end            string
	businessHourId string
	hourOffset     uint
}

type specialBusinessHourInput struct {
	name             string
	args             specialBusinessHourArgs
	want             specialBusinessHourArgs
	hasValidationErr bool
	hasNotFoundErr   bool
}

func assertSpecialHourRoot(t *testing.T, tt specialBusinessHourInput, got *store.SpecialBusinessHour, err error) {
	if tt.hasValidationErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewValidationError("", ""), err)
		return
	}
	if tt.hasNotFoundErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewNotFoundError(""), err)
		return
	}
	assert.NoError(t, err, "no error should be")
	assertSpecialHour(t, tt.want, got)
}

func assertSpecialHour(t *testing.T, want specialBusinessHourArgs, got *store.SpecialBusinessHour) {
	assert.Equal(t, want.name, got.GetName())
	assert.Equal(t, want.date, got.GetDate())
	assert.Equal(t, want.start, got.GetStart())
	assert.Equal(t, want.end, got.GetEnd())
	assert.Equal(t, want.businessHourId, got.GetBusinessHourId())
	assert.Equal(t, want.hourOffset, got.GetHourOffset())
}

func getSpecialBusinessHourInput() *[]specialBusinessHourInput {
	inputs := []specialBusinessHourInput{
		{name: "normal case",
			args: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     1,
			},
			want: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     1,
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal case(12 offset hour)",
			args: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     12,
			},
			want: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     12,
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "error empty name",
			args: specialBusinessHourArgs{
				name:           "",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error max name(31)",
			args: specialBusinessHourArgs{
				name:           "1234567890123456789012345678901",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error date format",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "20220104",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error start error",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "2022/01/04",
				start:          "0800",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error end error",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "0900",
				businessHourId: "123",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error businessHourId error(empty)",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error start > end",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "2022/01/04",
				start:          "09:00",
				end:            "08:00",
				businessHourId: "124",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error start > end + 60",
			args: specialBusinessHourArgs{
				name:           "1234",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "08:59",
				businessHourId: "124",
				hourOffset:     3,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error zero offset",
			args: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     0,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error 13 offset",
			args: specialBusinessHourArgs{
				name:           "special morning",
				date:           "2022/01/04",
				start:          "08:00",
				end:            "09:00",
				businessHourId: "123",
				hourOffset:     13,
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}
	return &inputs
}

func TestNewSpecialHourInput(t *testing.T) {
	inputs := getSpecialBusinessHourInput()
	for _, tt := range *inputs {
		fmt.Println("name:", tt.name)
		got, err := store.NewSpecialBusinessHour(tt.args.name, tt.args.date, tt.args.start, tt.args.end, tt.args.businessHourId, tt.args.hourOffset)
		assertSpecialHourRoot(t, tt, got, err)
	}
}

func TestSpecialHourSet(t *testing.T) {
	inputs := getSpecialBusinessHourInput()
	for _, tt := range *inputs {
		fmt.Println("name:", tt.name)
		got, err := store.NewSpecialBusinessHour("init", "2022/10/10", "04:00", "05:00", "123", 4)
		if err != nil {
			assert.Fail(t, "init is failed")
		}
		err = got.Set(tt.args.name, tt.args.date, tt.args.start, tt.args.end, tt.args.businessHourId, tt.args.hourOffset)
		assertSpecialHourRoot(t, tt, got, err)
	}
}
