package order

import (
	"fmt"
	"time"

	"chico/takeout/common"
)

type OrderStatisticQueryService interface {
	FetchMonthlyStatistic(startMonth, endMonth time.Time) (*MonthlyStatisticData, error)
}

type OrderStatisticUseCase struct {
	service OrderStatisticQueryService
}

func NewOrderStatisticUseCase(queryService OrderStatisticQueryService) *OrderStatisticUseCase {
	return &OrderStatisticUseCase{
		service: queryService,
	}
}

type MonthlyStatisticRequestModel struct {
	Start string
	End string
}

type MonthlyStatisticData struct {
	Data []MonthlyData
}

type MonthlyData struct {
	Month         string
	OrderTotal    int
	QuantityTotal int
	MoneyTotal    int
}

func (o *OrderStatisticUseCase) FetchMonthlyData(req MonthlyStatisticRequestModel) (*MonthlyStatisticData, error) {
	start, err := common.ConvertStrToMonth(req.Start)
	if err != nil {
		return nil, common.NewValidationError("Start", fmt.Sprintf("failed to convert month data:%s", err))
	}
	
	end, err := common.ConvertStrToMonth(req.End)
	if err != nil {
		return nil, common.NewValidationError("End", fmt.Sprintf("failed to convert month data:%s", err))
	}

	// check duration
	if start.After(*end) || start.Equal(*end) {
		return nil, common.NewValidationError("Start, End", fmt.Sprintf("start should greater than end. start%s, end:%s", *start, *end))
	}
	// allow within 1 year
	limit := start.AddDate(1, 0, 0)
	if end.After(limit) {
		return nil, common.NewValidationError("Start, End", fmt.Sprintf("start should within 1year from end. start%s, end:%s", *start, *end))
	}
	return o.service.FetchMonthlyStatistic(*start, *end)
}
