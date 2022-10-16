package query

import (
	"time"

	"chico/takeout/common"
	order "chico/takeout/usecase/order/query"

	"gorm.io/gorm"
)

type OrderStatisticQueryService struct {
	db *gorm.DB
}

func NewOrderStatisticQueryService(db *gorm.DB) *OrderStatisticQueryService {
	return &OrderStatisticQueryService{
		db: db,
	}
}

func (o *OrderStatisticQueryService) FetchMonthlyStatistic(startMonth, endMonth time.Time) (*order.MonthlyStatisticData, error) {
	start := common.ConvertTimeToMonthStr(startMonth)
	end := common.ConvertTimeToMonthStr(endMonth.AddDate(0, 1, 0))

	monthlyData, err := o.fetchMonthlyData(start, end)
	if err != nil {
		return nil, err
	}

	months, err := common.ListUpMonths(startMonth, endMonth)
	if err != nil {
		return nil, err
	}
	models := []order.MonthlyData{}

	for _, month := range months {
		hasData := false
		monthStr := common.ConvertTimeToMonthStr(month)
		for _, d := range monthlyData {
			if d.Month == monthStr {
				models = append(models, d)
				hasData = true
				break
			}
		}
		if !hasData {
			// if there is no data add empty data
			d := order.MonthlyData{
				Month:         monthStr,
				OrderTotal:    0,
				QuantityTotal: 0,
				MoneyTotal:    0,
			}
			models = append(models, d)
		}
	}

	return &order.MonthlyStatisticData{
		Data: models,
	}, nil
}

func (o *OrderStatisticQueryService) fetchMonthlyData(startMonth, endMonth string) ([]order.MonthlyData, error) {
	models := []order.MonthlyData{}
	err := o.db.Raw(`select to_char(DATE_TRUNC('month', order_info_models.order_date_time), 'YYYY/MM') as month, count(order_info_models.order_date_time) as order_total, sum(sum_price) as money_total, sum(quantity) as quantity_total from
	(select order_info_model_id as order_id, food_item_model_id as item_id, name, price, quantity, (price * quantity) as sum_price from ordered_food_item_models
	union
	select order_info_model_id as order_id, stock_item_model_id as item_id, name, price, quantity, (price * quantity) as sum_price  from ordered_stock_item_models) as items
	inner join order_info_models
	on items.order_id = order_info_models.id
	 where order_info_models.canceled  = false and order_info_models.deleted_at is null
	 and order_info_models.order_date_time >= ? and order_info_models.order_date_time < ?
	 group by DATE_TRUNC('month', order_info_models.order_date_time);`, addDateAndSecond(startMonth), addDateAndSecond(endMonth)).Scan(&models).Error

	if err != nil {
		return nil, err
	}

	return models, nil
}

func addDateAndSecond(monthStr string) string {
	return monthStr + "/01 00:00:00.000"
}
