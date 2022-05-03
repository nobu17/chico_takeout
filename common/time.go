package common

import (
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func ConvertStrToTime(timeStr string) (*time.Time, error) {
	timeLayout := "2006-01-02T15:04"
	dateOnlyLayout := "2006-01-02"
	currentTime := time.Now()
	currentDate := currentTime.Format(dateOnlyLayout)

	startDateStr := currentDate + "T" + timeStr
	actualTime, err := time.ParseInLocation(timeLayout, startDateStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func CompareTime(start, end time.Time, offsetMinutes float64) bool {
	diff := start.Sub(end)
	return diff.Minutes() < offsetMinutes
}

func IsOverlap(start1, end1, start2, end2 time.Time) bool {
	// B.start(start2) < A.end(end1) && A.start(start1) < B.end(end2)
	return (start2.Before(end1) && start1.Before(end2))
}