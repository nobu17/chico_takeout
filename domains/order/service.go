package order

import (
	"chico/takeout/common"
	"chico/takeout/domains/item"
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
				err = s.stockRepo.Update(stock)
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
				err = stock.IncreseRemain(order.GetQuantity())
				if err != nil {
					return err
				}
				// update stock db
				err = s.stockRepo.Update(stock)
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
	samedateOrders, err := f.orderRepo.FindByPickupDate(pickupDateTime)
	if err != nil {
		return err
	}
	spec := newFoodItemRemainQuantitySpecification(samedateOrders)

	// step2: check each order remain
	foods, err := f.foodRepo.FindAll()
	if err != nil {
		return err
	}
	for _, foodOrder := range foodOrders {
		for _, food := range foods {
			if foodOrder.HasSameId(food.GetId()) {
				if spec.IsOverRemain(food.GetId(), foodOrder.GetQuantity(), food.GetMaxOrderPerDay()) {
					return common.NewValidationError("foodOrders", "Food items remain count perday is over limit.")
				}
				break
			}
		}
	}
	return nil
}
