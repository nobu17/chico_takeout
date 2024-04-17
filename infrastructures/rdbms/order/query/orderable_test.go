package query

import (
	"testing"
	"time"

	"chico/takeout/common"

	"chico/takeout/infrastructures/rdbms"
	"chico/takeout/infrastructures/rdbms/store"
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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
	// but today info is already passed (now * 2 hours is over the end time)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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
	// but today info is not started yet (now * 2 hours is before start time)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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

	// mock now (2022/10/5 23:20.30)
	// morning +7 hours is tomorrow 6:20 => 6:30
	// lunch  +3 hours is tomorrow 2:20 => 2:30
	// dinner +4 hours is tomorrow 3:20 => 3:30
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
				Date:       "2022/10/05", // before day. filtered
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05", // before day. filtered
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06", // 6:30 over filtered
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
		},
	}

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 7, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 4, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
	assert.NoError(t, err, "no error should be")

	// all today items are removed. but next date is still existed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_NowDateIsChangedToNextDate_BySpecialBusinessDay(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 23:20.30)
	// morning +7 hours is tomorrow 6:20 => 6:30
	// lunch  +3 hours is tomorrow 2:20 => 2:30
	// sp_lunch +11 hours is tomorrow 10:20 => 10:30
	// dinner +4 hours is tomorrow 3:20 => 3:30
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
				Date:       "2022/10/05", // before day. filtered
				HourTypeId: "1",
				StartTime:  "09:00",
				EndTime:    "13:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05", // before day. filtered
				HourTypeId: "2",
				StartTime:  "13:30",
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06", // 6:30 over filtered
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06", // 10:30 over filtered
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
		},
	}

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "morning", OffsetHour: 7, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 4, Enabled: true},
	}
	dt := time.Date(2022, 10, 6, 0, 0, 0, 0, time.Local)
	spHours := []store.SpecialBusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1234"}, Name: "sp_launch", OffsetHour: 11, Date: &dt, BusinessHourModelID: "2"},
		{BaseModel: rdbms.BaseModel{ID: "1235"}, Name: "sp_dinner", OffsetHour: 3, Date: &dt, BusinessHourModelID: "3"},
	}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
	assert.NoError(t, err, "no error should be")

	// all today items are removed. but next date is still existed
	expected := order.OrderableInfo{
		StartDate: "2022/10/03",
		EndDate:   "2022/10/20",
		PerDayInfo: []order.PerDayOrderableInfo{
			{
				Date:       "2022/10/06",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "1",
				StartTime:  "6:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "2",
				StartTime:  "10:00",
				EndTime:    "14:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/07",
				HourTypeId: "3",
				StartTime:  "17:00",
				EndTime:    "20:00",
				Items:      nil,
			},
		},
	}
	AssertOrderableInfo(t, expected, info)
}

func TestModifyTodayInfo_HasToday_ModifiedStartTime(t *testing.T) {
	o := OrderableInfoRdbmsQueryService{}

	// mock now (2022/10/5 11:10.30) => +2hours is 13:10 (rounded:13:30)
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
				StartTime:  "13:00", // this time is not allowed. (limit is 13:30)
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:29", // this time is not allowed. (limit is 13:30)
				EndTime:    "17:00",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "13:30", // this time is just allowed. (13:30)
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "1"}, Name: "dinner", OffsetHour: 3, Enabled: true},
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 2, Enabled: true}, // used
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
	assert.NoError(t, err, "no error should be")

	// 1~3 item is removed
	// other is remained
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

	// mock now (2022/10/5 11:10.30) => +2hours is 13:30
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
				StartTime:  "12:30", // this time is passed. but not end
				EndTime:    "13:30",
				Items:      nil,
			},
			{
				Date:       "2022/10/05",
				HourTypeId: "2",
				StartTime:  "12:30", // this time is passed. but not end
				EndTime:    "13:31",
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

	hours := []store.BusinessHourModel{
		{BaseModel: rdbms.BaseModel{ID: "2"}, Name: "lunch", OffsetHour: 2, Enabled: true}, // used
		{BaseModel: rdbms.BaseModel{ID: "3"}, Name: "dinner", OffsetHour: 3, Enabled: true},
	}
	spHours := []store.SpecialBusinessHourModel{}

	// act
	err := o.modifyTodayInfo(&info, hours, spHours)
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
