package common

import (
	"fmt"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

const dateLayout = "2006/01/02"

func ConvertStrToTime(timeStr string) (*time.Time, error) {
	timeLayout := "2006/01/02T15:04"
	currentTime := time.Now()
	currentDate := currentTime.Format(dateLayout)

	startDateStr := currentDate + "T" + timeStr
	actualTime, err := time.ParseInLocation(timeLayout, startDateStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func ConvertStrToDate(dateStr string) (*time.Time, error) {
	actualTime, err := time.ParseInLocation(dateLayout, dateStr, jst)
	if err != nil {
		return nil, err
	}

	return &actualTime, nil
}

func StartIsBeforeEnd(start, end time.Time, offsetMinutes float64) bool {
	diff := end.Sub(start)
	fmt.Println("diff", start, end, diff)
	return diff.Minutes() > offsetMinutes
}

func IsOverlap(start1, end1, start2, end2 time.Time) bool {
	// B.start(start2) < A.end(end1) && A.start(start1) < B.end(end2)
	return (start2.Before(end1) && start1.Before(end2))
}
