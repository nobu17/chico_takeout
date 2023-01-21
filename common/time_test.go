package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type inRangeTestParam struct {
	name  string
	start time.Time
	end   time.Time
	target time.Time
}

func Test_IsInRange_In(t *testing.T) {
	inputs := []inRangeTestParam {
		{ 
			name: "same start",
			start: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end : time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
		},
		{ 
			name: "after start and before end",
			start: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end : time.Date(2023, 1, 26, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
		},
		{ 
			name: "same before end",
			start: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end : time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
		},
	}

	for _, input := range inputs {
		result := IsInRange(input.start, input.end, input.target);
		assert.Equal(t, true, result, fmt.Sprintf("%v:case should be true", input.name))
	}
}

func Test_IsInRange_Out(t *testing.T) {
	inputs := []inRangeTestParam {
		{ 
			name: "before start",
			start: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end : time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 23, 00, 00, 00, 0, time.Local),
		},
		{ 
			name: "after end",
			start: time.Date(2023, 1, 24, 00, 00, 00, 0, time.Local),
			end : time.Date(2023, 1, 25, 00, 00, 00, 0, time.Local),
			target: time.Date(2023, 1, 26, 00, 00, 00, 0, time.Local),
		},
	}

	for _, input := range inputs {
		result := IsInRange(input.start, input.end, input.target);
		assert.Equal(t, false, result, fmt.Sprintf("%v:case should be false", input.name))
	}
}