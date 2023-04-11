package common

import (
	"time"
)

const offsetMinutes = 60

type DailySchedularTask struct {
	start time.Time
	last  *time.Time
	task  func()
}

func NewDailySchedularTask(startStr string, task func()) (*DailySchedularTask, error) {
	if task == nil {
		return nil, NewValidationError("task", "should not nil")
	}
	time, err := ConvertStrToTime(startStr)
	if err != nil {
		return nil, err
	}
	return &DailySchedularTask{
		start: *time,
		last:  nil,
		task:  task,
	}, nil
}

func (d *DailySchedularTask) CheckAndExecTask() {
	end := GetDateWithOffset(d.start, offsetMinutes)
	now := GetNowDate()
	if IsInRangeTime(d.start, *end, *now) {
		if d.last == nil || !DateEqual(*d.last, *now) {
			d.task()
			d.last = now
		}
	}
}

type TimerScheduleTask struct {
	intervalMinutes int
	task func(time.Time)
}

func NewTimerScheduleTask(intervalMinutes int, task func(time.Time)) (*TimerScheduleTask, error) {
	if task == nil {
		return nil, NewValidationError("task", "should not nil")
	}
	if intervalMinutes < 1 {
		return nil, NewValidationError("intervalMinutes", "should be greater than 0")
	}
	return &TimerScheduleTask{
		intervalMinutes: intervalMinutes,
		task:  task,
	}, nil
}

func (t *TimerScheduleTask) Start() {
	timer := time.NewTicker(time.Minute * time.Duration(t.intervalMinutes))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			t.task(*GetNowDateUntilMinutes())
		}
	}
}