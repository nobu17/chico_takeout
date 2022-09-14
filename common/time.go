package common

import (
	"fmt"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// for mock
var now = time.Now

const dateHyphenLayout = "2006-01-02"
const dateLayout = "2006/01/02"
const dateTimeLayout = "2006/01/02 15:04"
const timeLayout = "15:04"

// for UT
func MockNow(mock func() time.Time) {
	now = mock
}
func ResetNow() {
	now = time.Now
}

func GetNowDate() *time.Time {
	now := now().In(jst)
	t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, jst)
	return &t
}

func ConvertStrToTime(timeStr string) (*time.Time, error) {
	timeLayout := "2006/01/02T15:04"
	currentTime := now().In(jst)
	currentDate := currentTime.Format(dateLayout)

	startDateStr := currentDate + "T" + timeStr
	actualTime, err := time.ParseInLocation(timeLayout, startDateStr, jst)
	if err != nil {
		return nil, err
	}
	
	return &actualTime, nil
}

func ConvertTimeToTimeStr(target time.Time) string {
	return target.Format(timeLayout)
}

func ConvertStrToDate(dateStr string) (*time.Time, error) {
	actualTime, err := time.ParseInLocation(dateLayout, dateStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func ConvertHyphenStrToDate(dateStr string) (*time.Time, error) {
	actualTime, err := time.ParseInLocation(dateHyphenLayout, dateStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func ConvertDateTimeStrToDateStr(dateTimeStr string) (string, error) {
	actualTime, err := ConvertStrToDateTime(dateTimeStr)
	if err != nil {
		return "", err
	}

	return actualTime.Format(dateLayout), nil
}

func ConvertStrToDateTime(dateTimeStr string) (*time.Time, error) {
	actualTime, err := time.ParseInLocation(dateTimeLayout, dateTimeStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func ConvertTimeToDateTimeStr(target time.Time) string {
	return target.Format(dateTimeLayout)
}

func ConvertTimeToDateStr(target time.Time) string {
	return target.Format(dateLayout)
}

func StartIsBeforeEnd(start, end time.Time, offsetMinutes float64) bool {
	diff := end.Sub(start)
	return diff.Minutes() > offsetMinutes
}

func IsOverlap(start1, end1, start2, end2 time.Time) bool {
	// B.start(start2) < A.end(end1) && A.start(start1) < B.end(end2)
	return (start2.Before(end1) && start1.Before(end2))
}

func IsInRange(startDate, endDate, targetDateTime time.Time) bool {
	// start -1 to include start
	actualStart := startDate.AddDate(0, 0, -1)
	// end +1 to include end
	actualEnd := endDate.AddDate(0, 0, 1)

	isAfterStart := targetDateTime.After(actualStart)
	isBeforeEnd := targetDateTime.Before(actualEnd)
	return isAfterStart && isBeforeEnd
}

func IsInRangeTime(startTime, endTime, targetDateTime time.Time) bool {
	// compare as same date
	acStDate := time.Date(2020, time.December, 1, startTime.Hour(), startTime.Minute()-1, 0, 0, time.UTC)
	acEdDate := time.Date(2020, time.December, 1, endTime.Hour(), endTime.Minute()+1, 0, 0, time.UTC)
	comDate := time.Date(2020, time.December, 1, targetDateTime.Hour(), targetDateTime.Minute(), 0, 0, time.UTC)

	isAfterStart := comDate.After(acStDate)
	isBeforeEnd := comDate.Before(acEdDate)
	return isAfterStart && isBeforeEnd
}

func IsAfterFromNow(target time.Time, offsetHour int) bool {
	now := now().In(jst)
	actualTime := now.Add(time.Duration(offsetHour) * time.Hour)
	return target.After(actualTime)
}

func ListUpDates(startDate, endDate time.Time) ([]time.Time, error) {
	actStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	actEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 1, 0, 0, 0, time.UTC) // add 1minutes to compare

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date is after than end date. start:%s end:%s", startDate, endDate)
	}

	dates := []time.Time{}
	currentDate := actStart
	for currentDate.Before(actEnd) {
		date := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, time.UTC)
		dates = append(dates, date)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dates, nil
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
