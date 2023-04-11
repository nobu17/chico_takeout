package order

import (
	"fmt"
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/order"
	storeDomains "chico/takeout/domains/store"
)

type OrderTaskUseCase interface {
	NotifyDailyOrder(start time.Time) error
	NotifyOrderByHour(currentTime time.Time) error
}

type orderTaskUseCase struct {
	filter        domains.OrderFilter
	mailerService SendOrderMailService
	mngService    storeDomains.BusinessHourManagementService
}

func NewOrderTaskUseCase(
	orderRepos domains.OrderInfoRepository,
	mailerService SendOrderMailService,
	businessHoursRepository storeDomains.BusinessHoursRepository,
	specialHolidayRepository storeDomains.SpecialHolidayRepository,
	specialBusinessHourRepository storeDomains.SpecialBusinessHourRepository) OrderTaskUseCase {
	return &orderTaskUseCase{
		filter:        *domains.NewOrderFilter(orderRepos),
		mailerService: mailerService,
		mngService:    *storeDomains.NewBusinessHourManagementService(businessHoursRepository, specialHolidayRepository, specialBusinessHourRepository),
	}
}

func (o *orderTaskUseCase) NotifyDailyOrder(start time.Time) error {
	orders, err := o.filter.GetActiveOrderOfSpecifiedDay(start)
	if err != nil {
		return err
	}
	cfg := common.GetConfig().Mail
	var mailData *ReservationSummaryMailData
	if len(orders) == 0 {
		mailData, err = NewNoReservationSummaryMailData(cfg.From, cfg.Admin, start)
		if err != nil {
			return err
		}
	} else {
		mailData, err = NewReservationSummaryMailData(orders, cfg.From, cfg.Admin, start)
		if err != nil {
			return err
		}
	}

	return o.mailerService.SendDailySummary(*mailData)
}

// 30分間隔チェック => 敷居値 25分
// 5:00 => 5:00 ~ 5:25 ok
// 04:50(ng), 05:20(ok), 05:50(ng)
// 04:40(ng), 05:10(ok), 05:40(ng)
// 04:45(ng), 05:15(ok), 05:40(ng)
// 04:30(ng), 05:00(ok), 05:30(ng)

const (
	thEndMinutes = 25
	startBeforeMinutes = 120 // 2 hour
)

func (o *orderTaskUseCase) NotifyOrderByHour(currentTime time.Time) error {
	// Step1 get today business hours
	data, err := o.mngService.GetSpecificDateHour(currentTime)
	if err != nil {
		return err
	}
	if len(data.Hours) == 0 {
		fmt.Println("no hours:" + data.Date)
		return nil
	}

	// Step2 compare with current time and start time
	notifyTarget := []storeDomains.HourInfo{}
	for _, hour := range data.Hours {
		// start => start - 2 hour
		startTime, err := common.ConvertStrToTime(hour.StartTime)
		if err != nil {
			return err
		}
		thStart := startTime.Add(-(time.Minute * startBeforeMinutes))
		// end => start -2 hour +25 min
		thEnd := thStart.Add(time.Minute * thEndMinutes)

		// ex: (7:00 => 5:00 ~ 5:25 is scope notify)
		if common.IsInRangeTime(thStart, thEnd, currentTime) {
			notifyTarget = append(notifyTarget, hour)
		}
	}

	if len(notifyTarget) == 0 {
		fmt.Println("no notify target.")
		return nil
	}

	// Step3 get all orders and then filter matched time orders and then send mail
	return o.sendMail(currentTime, notifyTarget)
}

func (o *orderTaskUseCase) sendMail(currentTime time.Time, targets []storeDomains.HourInfo) error {
	for _, target := range targets {
		startTime, err := common.ConvertStrToTime(target.StartTime)
		if err != nil {
			return err
		}
		endTime, err := common.ConvertStrToTime(target.EndTime)
		if err != nil {
			return err
		}
		orders, err := o.filter.GetActiveOrderOfSpecifiedDayAndTime(currentTime, *startTime, *endTime)
		if err != nil {
			return err
		}
		cfg := common.GetConfig().Mail
		var mailData *ReservationSummaryMailData
		todayDateStr := common.ConvertTimeToDateStr(currentTime)
		if len(orders) == 0 {
			mailData, err = NewNoHourReservationSummaryMailData(cfg.From, cfg.Admin, todayDateStr, target.StartTime, target.EndTime)
			if err != nil {
				return err
			}
		} else {
			mailData, err = NewHourReservationSummaryMailData(orders, cfg.From, cfg.Admin, todayDateStr, target.StartTime, target.EndTime)
			if err != nil {
				return err
			}
		}
		err = o.mailerService.SendDailySummary(*mailData)
		if err != nil {
			return err
		}
	}
	return nil
}
