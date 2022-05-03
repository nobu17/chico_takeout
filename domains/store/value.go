package store

import (
	"fmt"

	"chico/takeout/common"
)

const (
	offsetMinutes = 60
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
	if !common.CompareTime(*startTime, *endTime, offsetMinutes) {
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