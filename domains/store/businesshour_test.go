package store_test

import (
	"fmt"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/store"

	"github.com/stretchr/testify/assert"
)

type busHoursArgs struct {
	id         string
	idIndex    int
	name       string
	start      string
	end        string
	weekdays   []store.Weekday
	enabled    bool
	hourOffset uint
}
type busHoursInput struct {
	name             string
	args             []busHoursArgs
	want             []busHoursArgs
	hasValidationErr bool
	hasNotFoundErr   bool
}

func assertBusinessHoursRoot(t *testing.T, tt busHoursInput, got *store.BusinessHours, err error) {
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
	assertBusinessHours(t, tt.want, *got)
}

func assertBusinessHours(t *testing.T, want []busHoursArgs, got store.BusinessHours) {
	assert.Equal(t, len(want), len(got.GetSchedules()))
	gotsch := got.GetSchedules()
	for index, wantsch := range want {
		assertBusinessHour(t, wantsch, gotsch[index])
	}
}

func assertBusinessHour(t *testing.T, want busHoursArgs, got store.BusinessHour) {
	assert.Equal(t, want.name, got.GetName())
	assert.Equal(t, want.start, got.GetStart())
	assert.Equal(t, want.end, got.GetEnd())
	assert.ElementsMatch(t, want.weekdays, got.GetWeekdays())
	assert.Equal(t, want.enabled, got.GetEnabled())
	assert.Equal(t, want.hourOffset, got.GetHourOffset())
}

func TestNewDefaultBusinessHours(t *testing.T) {
	inputs := []busHoursInput{
		{name: "normal check",
			args: []busHoursArgs{},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)
		got, err := store.NewDefaultBusinessHours()
		assertBusinessHoursRoot(t, tt, got, err)
	}
}

func TestBusinessHours_FindById(t *testing.T) {
	inputs := []busHoursInput{
		{name: "normal check",
			args: []busHoursArgs{},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)
		bus, err := store.NewDefaultBusinessHours()
		assert.NoError(t, err, "should not got error at NewDefaultBusinessHours()")

		schedules := bus.GetSchedules()
		for i, sch := range schedules {
			got := bus.FindById(sch.GetId())
			assert.NotNil(t, got, "should not be nil")
			assertBusinessHour(t, tt.want[i], *got)
		}
	}
}

func TestBusinessHours_FindById_NotFound(t *testing.T) {
	inputs := []busHoursInput{
		{name: "not found",
			args: []busHoursArgs{
				{id: "1234"},
			},
			want:             []busHoursArgs{},
			hasValidationErr: false,
			hasNotFoundErr:   true,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)
		bus, err := store.NewDefaultBusinessHours()
		assert.NoError(t, err, "test initialize failed")

		upInfo := tt.args[0]

		got, err := bus.Update(upInfo.id, upInfo.name, upInfo.start, upInfo.end, upInfo.weekdays, upInfo.hourOffset)
		assertBusinessHoursRoot(t, tt, got, err)
	}
}

func TestBusinessHours_Update(t *testing.T) {
	inputs := []busHoursInput{
		{name: "morning update",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning2", start: "08:00", end: "09:00", weekdays: []store.Weekday{store.Tuesday}, hourOffset: 2},
			},
			want: []busHoursArgs{
				{name: "morning2", start: "08:00", end: "09:00", enabled: true, weekdays: []store.Weekday{store.Tuesday}, hourOffset: 2},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "lunch update",
			args: []busHoursArgs{
				{idIndex: 1, name: "lunch2", start: "12:00", end: "14:00", weekdays: []store.Weekday{store.Wednesday, store.Friday}, hourOffset: 12},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch2", start: "12:00", end: "14:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Friday}, hourOffset: 12},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "dinner update",
			args: []busHoursArgs{
				{idIndex: 2, name: "dinner2", start: "19:00", end: "23:00", weekdays: []store.Weekday{store.Wednesday, store.Friday}, hourOffset: 1},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner2", start: "19:00", end: "23:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Friday}, hourOffset: 1},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "duplicate weekdays",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Tuesday, store.Sunday}, hourOffset: 2},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "irregular time(start < end)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "07:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Tuesday, store.Sunday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "irregular time(start < end + 59)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "08:59", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Tuesday, store.Sunday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "edge time(start < end + 60)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "09:00", weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 2},
			},
			want: []busHoursArgs{
				{name: "morning", start: "08:00", end: "09:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 2},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "overlap time1(morning end is overlap lunch)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "10:00", end: "12:00", weekdays: []store.Weekday{store.Tuesday}, hourOffset: 4},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "overlap time2(morning is include lunch)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "10:00", end: "16:00", weekdays: []store.Weekday{store.Tuesday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "overlap time3(morning is inside lunch",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "12:00", end: "14:00", weekdays: []store.Weekday{store.Tuesday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "overlap time4(morning end is overlap lunch",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "10:00", end: "17:00", weekdays: []store.Weekday{store.Tuesday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "time is overlap but weekday is not match",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "13:00", weekdays: []store.Weekday{store.Monday}, hourOffset: 4},
			},
			want: []busHoursArgs{
				{name: "morning", start: "08:00", end: "13:00", enabled: true, weekdays: []store.Weekday{store.Monday}, hourOffset: 4},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "irregular time format(start)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "0800", end: "13:00", weekdays: []store.Weekday{store.Monday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "irregular time format(end)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "13:a0", weekdays: []store.Weekday{store.Monday}, hourOffset: 3},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "no weekdays (allowed)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "13:00", weekdays: []store.Weekday{}, hourOffset: 4},
			},
			want: []busHoursArgs{
				{name: "morning", start: "08:00", end: "13:00", enabled: true, weekdays: []store.Weekday{}, hourOffset: 4},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "0 offset hour",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Sunday}, hourOffset: 0},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "13 offset hour(over limit)",
			args: []busHoursArgs{
				{idIndex: 0, name: "morning", start: "08:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Sunday}, hourOffset: 13},
			},
			want:             []busHoursArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {

		fmt.Println("name:", tt.name)
		bus, err := store.NewDefaultBusinessHours()
		assert.NoError(t, err, "test initialize failed")

		schedules := bus.GetSchedules()
		upInfo := tt.args[0]

		got, err := bus.Update(schedules[upInfo.idIndex].GetId(), upInfo.name, upInfo.start, upInfo.end, upInfo.weekdays, upInfo.hourOffset)
		assertBusinessHoursRoot(t, tt, got, err)
	}
}

func TestBusinessHours_Update_Enabled(t *testing.T) {
	inputs := []busHoursInput{
		{name: "morning update (true to false)",
			args: []busHoursArgs{
				{idIndex: 0, enabled: false},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: false, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "morning true to true update",
			args: []busHoursArgs{
				{idIndex: 0, enabled: true},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "lunch false to false update",
			args: []busHoursArgs{
				{idIndex: 1, enabled: false},
				{idIndex: 1, enabled: false},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: false, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "lunch false to true update",
			args: []busHoursArgs{
				{idIndex: 1, enabled: false},
				{idIndex: 1, enabled: true},
			},
			want: []busHoursArgs{
				{name: "morning", start: "07:00", end: "09:30", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "lunch", start: "11:30", end: "15:00", enabled: true, weekdays: []store.Weekday{store.Tuesday, store.Wednesday, store.Friday, store.Saturday, store.Sunday}, hourOffset: 3},
				{name: "dinner", start: "18:00", end: "21:00", enabled: true, weekdays: []store.Weekday{store.Wednesday, store.Saturday}, hourOffset: 3},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {

		fmt.Println("name:", tt.name)
		bus, err := store.NewDefaultBusinessHours()
		assert.NoError(t, err, "test initialize failed")

		schedules := bus.GetSchedules()
		upInfo := tt.args[0]
		got, err := bus.UpdateEnabled(schedules[upInfo.idIndex].GetId(), upInfo.enabled)
		if len(tt.args) > 1 {
			upInfo := tt.args[1]
			got, err := bus.UpdateEnabled(schedules[upInfo.idIndex].GetId(), upInfo.enabled)
			assertBusinessHoursRoot(t, tt, got, err)
		} else {
			assertBusinessHoursRoot(t, tt, got, err)
		}
	}
}
