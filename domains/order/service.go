package order

import (
	"chico/takeout/common"
	"chico/takeout/domains/item"
	"time"
)

type StockItemRemainCheckAndConsumer struct {
	stockRepo item.StockItemRepository
}

func NewStockItemRemainCheckAndConsumer(stockRepo item.StockItemRepository) *StockItemRemainCheckAndConsumer {
	return &StockItemRemainCheckAndConsumer{
		stockRepo: stockRepo,
	}
}

func (s *StockItemRemainCheckAndConsumer) ConsumeRemainStock(stockOrders []ItemOrder) error {
	allStocks, err := s.stockRepo.FindAll()
	if err != nil {
		return err
	}
	for _, order := range stockOrders {
		for _, stock := range allStocks {
			if stock.HasSameId(order.id) {
				err = stock.ConsumeRemain(order.quantity)
				// out of stock
				if err != nil {
					return err
				}
				// update stock db
				err = s.stockRepo.Update(&stock)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func (s *StockItemRemainCheckAndConsumer) IncrementCanceledRemain(stockOrders []OrderStockItem) error {
	allStocks, err := s.stockRepo.FindAll()
	if err != nil {
		return err
	}
	for _, order := range stockOrders {
		for _, stock := range allStocks {
			if stock.HasSameId(order.GetItemId()) {
				err = stock.IncreaseRemain(order.GetQuantity())
				if err != nil {
					return err
				}
				// update stock db
				err = s.stockRepo.Update(&stock)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

type FoodItemRemainChecker struct {
	orderRepo OrderInfoRepository
	foodRepo  item.FoodItemRepository
}

func NewFoodItemRemainChecker(orderRepo OrderInfoRepository, foodRepo item.FoodItemRepository) *FoodItemRemainChecker {
	return &FoodItemRemainChecker{
		orderRepo: orderRepo,
		foodRepo:  foodRepo,
	}
}

func (f *FoodItemRemainChecker) CheckRemain(pickupDateTime string, foodOrders []OrderFoodItem) error {
	// step1 get same days food order and calc each quantity
	sameDateOrders, err := f.orderRepo.FindByPickupDate(pickupDateTime)
	if err != nil {
		return err
	}
	spec := newFoodItemRemainQuantitySpecification(sameDateOrders)

	// step2: check each order remain
	foods, err := f.foodRepo.FindAll()
	if err != nil {
		return err
	}
	for _, foodOrder := range foodOrders {
		for _, food := range foods {
			if foodOrder.HasSameId(food.GetId()) {
				if spec.IsOverRemain(food.GetId(), foodOrder.GetQuantity(), food.GetMaxOrderPerDay()) {
					return common.NewValidationError("foodOrders", "Food items remain count perDay is over limit.")
				}
				break
			}
		}
	}
	return nil
}

type OrderFilter struct {
	orderRepo OrderInfoRepository
}

func NewOrderFilter(orderRepo OrderInfoRepository) *OrderFilter {
	return &OrderFilter{
		orderRepo: orderRepo,
	}
}

func (o *OrderFilter) GetActiveOrderOfSpecifiedDay(startDateTime time.Time) ([]OrderInfo, error) {
	// fetch orders of specified date
	orders, err := o.orderRepo.FindByPickupDate(common.ConvertTimeToDateStr(startDateTime))
	if err != nil {
		return nil, err
	}
	// check active and after time
	target := []OrderInfo{}
	for _, order := range orders {
		if !order.canceled && order.pickupDateTime.GetDateTime().After(startDateTime) {
			target = append(target, order)
		}
	}
	return target, nil
}

func (o *OrderFilter) GetActiveOrderOfSpecifiedDayAndTime(startDate, startTime, endTime time.Time) ([]OrderInfo, error) {
	// fetch orders of specified date
	orders, err := o.orderRepo.FindByPickupDate(common.ConvertTimeToDateStr(startDate))
	if err != nil {
		return nil, err
	}
	// check active and in range time
	target := []OrderInfo{}
	for _, order := range orders {
		if order.canceled {
			continue
		}
		pickUpTime := order.pickupDateTime.GetDateTime()	
		if common.IsInRangeTime(startTime, endTime, pickUpTime) {
			target = append(target, order)
		}
	}
	return target, nil
}

type OrderDuplicateChecker struct {
	orderRepo OrderInfoRepository
}

func NewOrderDuplicateChecker(orderRepo OrderInfoRepository) *OrderDuplicateChecker {
	return &OrderDuplicateChecker{
		orderRepo: orderRepo,
	}
}

func (o *OrderDuplicateChecker) ActiveOrderExists(userId string) (bool, error) {
	order, err := o.orderRepo.FindActiveByUserId(userId)
	if err != nil {
		return false, err
	}
	if len(order) > 0 {
		return true, nil
	}
	return false, nil
}
