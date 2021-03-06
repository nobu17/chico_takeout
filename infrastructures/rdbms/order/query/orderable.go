package query

import (
	"fmt"
	"time"

	"chico/takeout/common"
	"chico/takeout/infrastructures/rdbms/items"
	"chico/takeout/infrastructures/rdbms/store"

	order "chico/takeout/usecase/order/query"

	"gorm.io/gorm"
)

type OrderableInfoRdbmsQueryService struct {
	db *gorm.DB
}

func NewOrderableInfoRdbmsQueryService(db *gorm.DB) *OrderableInfoRdbmsQueryService {
	return &OrderableInfoRdbmsQueryService{
		db: db,
	}
}

func (o *OrderableInfoRdbmsQueryService) FetchByDate(startDate, endDate time.Time) (*order.OrderableInfo, error) {
	// get holidays
	holidays := []store.SpecialHolidayModel{}
	// end is need escape
	err := o.db.Where("start <= ? and \"end\" >= ?", startDate, endDate).Find(&holidays).Error
	if err != nil {
		return nil, err
	}
	// get special business hour
	specialHours := []store.SpecialBusinessHourModel{}
	err = o.db.Preload("BusinessHourModel").Where("date <= ? and date >= ?", startDate, endDate).Find(&specialHours).Error
	if err != nil {
		return nil, err
	}
	// get business hour
	hours := []store.BusinessHourModel{}
	err = o.db.Preload("Weekdays").Find(&hours).Error
	if err != nil {
		return nil, err
	}

	availableDates := []time.Time{}
	dates, err := common.ListUpDates(startDate, endDate)
	if err != nil {
		return nil, err
	}
	// at first check holiday or not
	for _, date := range dates {
		for _, holiday := range holidays {
			if common.IsInRangeTime(*holiday.Start, *holiday.End, date) {
				continue
			}
		}
		availableDates = append(availableDates, date)
	}

	stocks, err := o.getStockItems()
	if err != nil {
		return nil, err
	}

	foods := []items.FoodItemModel{}
	err = o.db.Preload("BusinessHours").Find(&foods).Error
	if err != nil {
		return nil, err
	}

	foodConsumption, err := o.getPerDateFoodOrder(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// then check each hour
	infoLists := []order.PerDayOrderableInfo{}
	for _, date := range availableDates {
		hasSpecialHour := false
		// check special
		for _, specialHour := range specialHours {
			// special hour
			if common.DateEqual(date, *specialHour.Date) {
				foodItems := o.getFoodItems(specialHour.BusinessHourModelID, foods)
				o.reduceFoodRemain(foodItems, foodConsumption, date)
				allItems := append(foodItems, stocks...)
				info := order.PerDayOrderableInfo{
					Date:       common.ConvertTimeToDateStr(date),
					HourTypeId: specialHour.BusinessHourModelID,
					StartTime:  common.ConvertTimeToTimeStr(*specialHour.BusinessHourModel.Start),
					EndTime:    common.ConvertTimeToTimeStr(*specialHour.BusinessHourModel.End),
					Items:      allItems,
				}
				infoLists = append(infoLists, info)
				hasSpecialHour = true
			}
		}
		// if not match, check normal
		if !hasSpecialHour {
			weekday := date.Weekday()
			for _, hour := range hours {
				if hour.HasWeekDay(int(weekday)) {
					foodItems := o.getFoodItems(hour.ID, foods)
					o.reduceFoodRemain(foodItems, foodConsumption, date)
					allItems := append(foodItems, stocks...)
					info := order.PerDayOrderableInfo{
						Date:       common.ConvertTimeToDateStr(date),
						HourTypeId: hour.ID,
						StartTime:  common.ConvertTimeToTimeStr(*hour.Start),
						EndTime:    common.ConvertTimeToTimeStr(*hour.End),
						Items:      allItems,
					}
					infoLists = append(infoLists, info)
				}
			}
		}
	}

	data := order.OrderableInfo{}
	data.StartDate = common.ConvertTimeToDateStr(startDate)
	data.EndDate = common.ConvertTimeToDateStr(endDate)
	data.PerDayInfo = infoLists
	return &data, nil
}

func (o *OrderableInfoRdbmsQueryService) getStockItems() ([]order.OrderableItemInfo, error) {
	models := []items.StockItemModel{}

	err := o.db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	infoList := []order.OrderableItemInfo{}

	for _, item := range models {
		info := order.OrderableItemInfo{}
		info.Id = item.ID
		info.ItemType = "stock"
		info.Remain = item.Remain

		infoList = append(infoList, info)
	}

	return infoList, nil
}

func (o *OrderableInfoRdbmsQueryService) getFoodItems(hourTypeId string, foods []items.FoodItemModel) []order.OrderableItemInfo {
	infoList := []order.OrderableItemInfo{}

	for _, item := range foods {
		belongs := false
		for _, hour := range item.BusinessHours {
			if hour.ID == hourTypeId {
				belongs = true
				break
			}
		}
		if !belongs {
			continue
		}
		info := order.OrderableItemInfo{}
		info.Id = item.ID
		info.ItemType = "food"
		info.Remain = item.MaxOrderPerDay

		infoList = append(infoList, info)
	}

	return infoList
}

func (o *OrderableInfoRdbmsQueryService) reduceFoodRemain(items []order.OrderableItemInfo, perDayOrder []foodOrderPerDayOrderedData, date time.Time) {
	for _, order := range perDayOrder {
		if order.PickUpDate == date {
			for _, item := range items {
				// reduce the remain
				if item.Id == order.Id {
					item.Remain = item.Remain - order.Quantity
				}
			}
		}
	}
}

func (o *OrderableInfoRdbmsQueryService) getPerDateFoodOrder(startDate, endDate time.Time) ([]foodOrderPerDayOrderedData, error) {
	models := []foodOrderPerDayOrderedData{}
	o.db.Raw(`select pick_up_date, food_item_model_id as id, food_order_quantity.quantity as quantity 
	from (select order_info.pick_up_date, food_order.food_item_model_id , SUM(food_order.quantity) as quantity from (select *, CAST(pickup_date_time AS DATE) as pick_up_date  from order_info_models where pickup_date_time >= ? and pickup_date_time  <= ?) as order_info
	 inner join ordered_food_item_models as food_order on order_info.id = food_order.order_info_model_id
	 group by food_order.food_item_model_id, order_info.pick_up_date) as food_order_quantity
	 left join food_item_models as food_models on food_order_quantity.food_item_model_id = food_models.id
	 order by pick_up_date`, startDate, endDate).Scan(&models)

	for _, ite := range models {
		fmt.Println("model", ite.Id, ite.PickUpDate, ite.Quantity)
	}

	return models, nil
}

type foodOrderPerDayOrderedData struct {
	PickUpDate time.Time
	Id         string
	Quantity   int
}
