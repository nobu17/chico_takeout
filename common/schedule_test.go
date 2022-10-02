package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDailySchedularTask_DoTaskOnSchedule(t *testing.T) {

	// mock now (2022/10/5 10:20.30)
	MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		ResetNow()
	})

	counter := 0
	task := func() {
		counter++
	}

	sch, err := NewDailySchedularTask("10:20", task)
	assert.NoError(t, err)

	// act
	sch.CheckAndExecTask()

	assert.Equal(t, 1, counter, "Should be called task.")
}

func TestDailySchedularTask_NotTaskAgain_UntilNextDay(t *testing.T) {

	// mock now (2022/10/5 10:20.30)
	MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		ResetNow()
	})

	counter := 0
	task := func() {
		counter++
	}

	sch, err := NewDailySchedularTask("10:20", task)
	assert.NoError(t, err)

	// call first time
	sch.CheckAndExecTask()
	assert.Equal(t, 1, counter, "Should be called task.")

	// call again
	sch.CheckAndExecTask()
	assert.Equal(t, 1, counter, "Should not be called task again.")

	// increment a day
	MockNow(func() time.Time {
		return time.Date(2022, 10, 6, 10, 20, 30, 0, time.Local)
	})
	// call
	sch.CheckAndExecTask()
	assert.Equal(t, 2, counter, "Should be called task.")
}

func TestDailySchedularTask_DoNotTaskBeforeSchedule(t *testing.T) {

	// mock now (2022/10/5 10:20.30)
	MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		ResetNow()
	})

	counter := 0
	task := func() {
		counter++
	}

	sch, err := NewDailySchedularTask("10:23", task)
	assert.NoError(t, err)

	// act
	sch.CheckAndExecTask()

	assert.Equal(t, 0, counter, "Should not be called task.")
}

func TestDailySchedularTask_DoNotTaskAfterSchedule(t *testing.T) {

	// mock now (2022/10/5 10:20.30)
	MockNow(func() time.Time {
		return time.Date(2022, 10, 5, 10, 20, 30, 0, time.Local)
	})
	t.Cleanup(func() {
		ResetNow()
	})

	counter := 0
	task := func() {
		counter++
	}

	// 1 hour later
	sch, err := NewDailySchedularTask("9:19", task)
	assert.NoError(t, err)

	// act
	sch.CheckAndExecTask()

	assert.Equal(t, 0, counter, "Should not be called task.")
}
