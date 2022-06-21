package order

import (
	"time"

	"chico/takeout/common"
)

type OrderableInfoQueryService interface {
	FetchByDate(startDate, endDate time.Time) (*OrderableInfo, error)
}

type OrderQueryUseCase struct {
	service OrderableInfoQueryService
}

func NewOrderQueryUseCase(queryService OrderableInfoQueryService) *OrderQueryUseCase {
	return &OrderQueryUseCase{
		service: queryService,
	}
}

type OrderableInfo struct {
	StartDate  string
	EndDate    string
	PerDayInfo []PerDayOrderableInfo
}

type PerDayOrderableInfo struct {
	Date       string
	HourTypeId string
	StartTime  string
	EndTime    string
	Items      []OrderableItemInfo
}

type OrderableItemInfo struct {
	Id       string
	ItemType string
	Remain   int
}

func (o *OrderQueryUseCase) FetchOrderableInfo() (*OrderableInfo, error) {
	// start is now
	start := common.GetNowDate()
	// 1week
	end := start.AddDate(0, 0, 7)
	return o.service.FetchByDate(*start, end)
}
