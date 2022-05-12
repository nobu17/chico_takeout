package store_test

import (
	"fmt"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/store"

	"github.com/stretchr/testify/assert"
)

type specialHoliDaysArgs struct {
	name  string
	start string
	end   string
}
type specialHoliDaysInpus struct {
	name             string
	args             specialHoliDaysArgs
	want             specialHoliDaysArgs
	hasValidationErr bool
	hasNotFoundErr   bool
}

func assertSpecialHolidayRoot(t *testing.T, tt specialHoliDaysInpus, got *store.SpecialHoliday, err error) {
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
	assert.NoError(t, err, "Update() failed")
	assertSpecialHoliday(t, tt.want, got)
}

func assertSpecialHoliday(t *testing.T, want specialHoliDaysArgs, got *store.SpecialHoliday) {
	assert.Equal(t, want.name, got.GetName())
	assert.Equal(t, want.start, got.GetStart())
	assert.Equal(t, want.end, got.GetEnd())
}

func getSpecialHoliDaysInpus() *[]specialHoliDaysInpus {
	inputs := []specialHoliDaysInpus{
		{name: "normal case",
			args: specialHoliDaysArgs{
				name:  "special1",
				start: "2022/01/04",
				end:   "2022/02/10",
			},
			want: specialHoliDaysArgs{
				name:  "special1",
				start: "2022/01/04",
				end:   "2022/02/10",
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "error name(empty)",
			args: specialHoliDaysArgs{
				name:  "",
				start: "2022/01/04",
				end:   "2022/02/10",
			},
			want:             specialHoliDaysArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error name(over 20)",
			args: specialHoliDaysArgs{
				name:  "123456789012345678901",
				start: "2022/01/04",
				end:   "2022/02/10",
			},
			want:             specialHoliDaysArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error start(incorrect format)",
			args: specialHoliDaysArgs{
				name:  "1234",
				start: "20220104",
				end:   "2022/02/10",
			},
			want:             specialHoliDaysArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error end(incorrect format)",
			args: specialHoliDaysArgs{
				name:  "1234",
				start: "2022/01/04",
				end:   "20220210",
			},
			want:             specialHoliDaysArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error start > end(incorrect format)",
			args: specialHoliDaysArgs{
				name:  "1234",
				start: "2022/02/11",
				end:   "2022/02/10",
			},
			want:             specialHoliDaysArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "normal start == end",
			args: specialHoliDaysArgs{
				name:  "1234",
				start: "2022/02/11",
				end:   "2022/02/11",
			},
			want: specialHoliDaysArgs{
				name:  "1234",
				start: "2022/02/11",
				end:   "2022/02/11",
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
	}
	return &inputs
}

func TestNewSpecialHoliday(t *testing.T) {
	inputs := getSpecialHoliDaysInpus()
	for _, tt := range *inputs {
		fmt.Println("name:", tt.name)
		got, err := store.NewSpecialHoliday(tt.args.name, tt.args.start, tt.args.end)
		assertSpecialHolidayRoot(t, tt, got, err)
	}
}

func TestSetSpecialHoliday(t *testing.T) {
	inputs := getSpecialHoliDaysInpus()
	for _, tt := range *inputs {
		fmt.Println("name:", tt.name)
		got, err := store.NewSpecialHoliday("holidays", "2022/04/01", "2022/05/02")
		assert.NoError(t, err, "initializing is failed.")
		err = got.Set(tt.args.name, tt.args.start, tt.args.end)
		assertSpecialHolidayRoot(t, tt, got, err)
	}
}
