package order

import (
	"os"
	"time"

	domains "chico/takeout/domains/order"
)

type OrderTaskUseCase interface {
	NotifyDailyOrder(start time.Time) error
}

type orderTaskUseCase struct {
	filter        domains.OrderFilter
	mailerService SendOrderMailService
}

func NewOrderTaskUseCase(orderRepos domains.OrderInfoRepository, mailerService SendOrderMailService) OrderTaskUseCase {
	return &orderTaskUseCase{
		filter: *domains.NewOrderFilter(orderRepos),
		mailerService: mailerService,
	}
}

func (o *orderTaskUseCase) NotifyDailyOrder(start time.Time) error {
	orders, err := o.filter.GetActiveOrderOfSpecifiedDay(start)
	if err != nil {
		return err
	}
	var mailData *ReservationSummaryMailData;
	if len(orders) == 0 {
		mailData, err = NewNoReservationSummaryMailData(os.Getenv("MAIL_FROM"), start)
		if err != nil {
			return err
		}
	} else {
		mailData, err = NewReservationSummaryMailData(orders, os.Getenv("MAIL_FROM"), start)
		if err != nil {
			return err
		}
	}

	return o.mailerService.SendDailySummary(*mailData)
}
