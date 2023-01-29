package query

import (
	"testing"
	"time"

	"chico/takeout/common"

	order "chico/takeout/usecase/order/query"

	"github.com/stretchr/testify/assert"
)

func TestModifyTodayInfo_HasPastDate_Filtered(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 10:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// not have today (10/5) info
	// has past (10/3, 4) info
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/03",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/04",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/04",
				HourTypeId: "2",
				StartTime:  "15:00",
				EndTime:    "19:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// expected is only have future date
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_OnlyFutureDates_NotFiltered(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 10:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// not have today (10/5) info
	// only has 10/6 ~ info
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "15:00",
				EndTime:    "19:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/08",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// expected has same input info
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "15:00",
				EndTime:    "19:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/08",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_AlreadyPassed(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 10:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	// but today info is already passed (now * 3 hours is over the end time)
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "13:20",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// today is removed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_NotPassed(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 10:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	// but today info is not started yet (now * 3 hours is before start time)
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// today is existed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_PartiallyPassed(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 10:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	// 1 item is passed, but another is still not passed yet
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// only 1 today item is removed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_AllItemsPassed(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 17:20.30)
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 17, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	// both items are passed
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// all today items are removed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_NowDateIsChangedToNextDate(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 23:20.30) => +3 hours is tomorrow
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 23, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// all today items are removed. but next date is still existed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_ModifiedStartTime(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 11:10.30) => +3hours is 14:10
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 11, 10, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30", // this time is passed. but not end
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05", // this info is not affected
				HourTypeId: "2",
				StartTime:  "19:00",
				EndTime:    "21:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// first item is removed
	// second item is modified start time. (14:30)
	// third. fourth items are not modified
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "14:30", // modified
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "19:00",
				EndTime:    "21:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_ModifiedStartTime_SameEndTime(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 11:10.30) => +3hours is 14:10
	common.MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 11, 10, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		common.ResetNow()
	})

	// have today (10/5) info
	info := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30", // this time is passed. but not end
				EndTime:    "14:30", // start will be 14:30. (same start and end)
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}

	// act
	err := o.modifyTodayInfo(&info)
	assert.NoError(t, err, "no error should be")

	// start and end is same case is removed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func AssertOrderableInfo(t *testing.T, expected, actual order.OrderableInfo) {
	assert.Equal(t, expected.StartDate, actual.StartDate)
	assert.Equal(t, expected.EndDate, actual.EndDate)
	assert.Equal(t, len(expected.PerDayInfo), len(actual.PerDayInfo), "PerDayInfo length is not match.")

	for i, item := range expected.PerDayInfo {
		AssertPerDayOrderableInfo(t, item, actual.PerDayInfo[i])
	}
}

func AssertPerDayOrderableInfo(t *testing.T, expected, actual order.PerDayOrderableInfo) {
	assert.Equal(t, expected.Date, actual.Date)
	assert.Equal(t, expected.HourTypeId, actual.HourTypeId)
	assert.Equal(t, expected.StartTime, actual.StartTime)
	assert.Equal(t, expected.EndTime, actual.EndTime)
	// item is not check target. only check both are nil.
	assert.Nil(t, expected.Items, "This text context Items should not set without nil.")
	assert.Nil(t, actual.Items, "This text context Items should not set without nil.")
}