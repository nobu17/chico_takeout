package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type inRangeTestParam struct {
	name   string
	start  time.Time
	end    time.Time
	target time.Time
}

func Test_IsInRange_In(t *testing.T) {
	inputs := []inRangeTestParam{
		{
			name:   "same start",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "+1 nsec from start",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 00, 00, 00, 1, time.Local),
		},
		{
			name:   "after start and before end",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 26, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "same before end",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "after end +1nsec",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 00, 00, 00, 1, time.Local),
		},
		{
			name:   "after end limit (-1nsec before next date",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 23, 59, 59, 900, time.Local),
		},
		{
			name:   "same date(start edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "same date(start +1",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 00, 00, 00, 1, time.Local),
		},
		{
			name:   "same date(normal",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 12, 15, 00, 1, time.Local),
		},
		{
			name:   "same date(end edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 23, 59, 59, 999, time.Local),
		},
	}

	for _, input := range inputs {
		result := IsInRange(input.start, input.end, input.target)
		assert.Equal(t, true, result, fmt.Sprintf("%v:case should be true", input.name))
	}
}

func Test_IsInRange_Out(t *testing.T) {
	inputs := []inRangeTestParam{
		{
			name:   "within 1 day before",
			start:  time.Date(2023, 2, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 3, 01, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 2, 23, 11, 00, 00, 0, time.Local),
		},
		{
			name:   "before start(edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 23, 23, 59, 59, 999, time.Local),
		},
		{
			name:   "before start",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 23, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "before start(edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 23, 23, 59, 59, 999, time.Local),
		},
		{
			name:   "after end(edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 26, 00, 00, 00, 0, time.Local),
		},
		{
			name:   "after end",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 26, 11, 00, 00, 0, time.Local),
		},
		{
			name:   "same date(start edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 23, 23, 59, 59, 999, time.Local),
		},
		{
			name:   "same date(end edge",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 0, 0, 0, 0, time.Local),
		},
		{
			name:   "same date(end next day",
			start:  time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end:    time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 12, 30, 0, 0, time.Local),
		},
	}

	for _, input := range inputs {
		result := IsInRange(input.start, input.end, input.target)
		assert.Equal(t, false, result, fmt.Sprintf("%v:case should be false", input.name))
	}
}
