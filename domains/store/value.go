package store

import (
	"fmt"

	"chico/takeout/common"
	"chico/takeout/domains/shared"
)

const (
	offsetMinutes = 59
)

type TimeRange struct {
	start string
	end   string
}

func NewTimeRange(start, end string) (*TimeRange, error) {
	// check start <= end with duration
	startTime, err := common.ConvertStrToTime(start)
	if err != nil {
		return nil, common.NewValidationError("start", fmt.Sprintf("can not convert time:%s", start))
	}
	endTime, err := common.ConvertStrToTime(end)
	if err != nil {
		return nil, common.NewValidationError("end", fmt.Sprintf("can not convert time:%s", start))
	}
	if !common.StartIsBeforeEnd(*startTime, *endTime, offsetMinutes) {
		return nil, common.NewValidationError("start end time", fmt.Sprintf("start time(%s) should be greater than end time(%s) with offset(%d)", start, end, offsetMinutes))
	}
	return &TimeRange{start: start, end: end}, nil
}

func (t *TimeRange) IsOverlap(other TimeRange) bool {
	tStart, _ := common.ConvertStrToTime(t.start)
	tEnd, _ := common.ConvertStrToTime(t.end)

	oStart, _ := common.ConvertStrToTime(other.start)
	oEnd, _ := common.ConvertStrToTime(other.end)

	return common.IsOverlap(*tStart, *tEnd, *oStart, *oEnd)
}

func (t *TimeRange) GetStart() string {
	return t.start
}

func (t *TimeRange) GetEnd() string {
	return t.end
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// date format is yyyy/MM/dd
type DateRange struct {
	start string
	end   string
}

func NewDateRange(start, end string) (*DateRange, error) {
	// check start <= end with duration
	startDate, err := common.ConvertStrToDate(start)
	if err != nil {
		return nil, common.NewValidationError("start", fmt.Sprintf("can not convert date:%s", start))
	}
	endDate, err := common.ConvertStrToDate(end)
	if err != nil {
		return nil, common.NewValidationError("end", fmt.Sprintf("can not convert date:%s", end))
	}
	// allow start == end (-1)
	if !common.StartIsBeforeEnd(*startDate, *endDate, -1) {
		return nil, common.NewValidationError("start end end", fmt.Sprintf("start date(%s) should be greater than end date(%s)", start, end))
	}
	return &DateRange{start: start, end: end}, nil
}

func (d *DateRange) GetStart() string {
	return d.start
}

func (d *DateRange) GetEnd() string {
	return d.end
}

func (d *DateRange) IsOverlap(other DateRange) bool {
	tStart, _ := common.ConvertStrToDate(d.start)
	tEnd, _ := common.ConvertStrToDate(d.end)

	oStart, _ := common.ConvertStrToDate(other.start)
	oEnd, _ := common.ConvertStrToDate(other.end)

	return common.IsOverlap(*tStart, *tEnd, *oStart, *oEnd)
}

type Date struct {
	shared.StringValue
}

func NewDate(value string) (*Date, error) {
	_, err := common.ConvertStrToDate(value)
	if err != nil {
		return nil, common.NewValidationError("date", fmt.Sprintf("can not convert date:%s", value))
	}
	return &Date{StringValue: shared.NewStringValue(value)}, nil
}
