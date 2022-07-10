package order

import (
	"chico/takeout/common"
	"chico/takeout/domains/item"
)

type ItemOrder struct {
	id       string
	quantity int
}

func NewItemOrder(id string, quantity int) *ItemOrder {
	return &ItemOrder{
		id:       id,
		quantity: quantity,
	}
}

type OrderInfoFactory struct {
	stockRepo item.StockItemRepository
	foodRepo  item.FoodItemRepository
}

func NewOrderInfoFactory(stockRepo item.StockItemRepository, foodRepo item.FoodItemRepository) *OrderInfoFactory {
	return &OrderInfoFactory{
		stockRepo: stockRepo,
		foodRepo:  foodRepo,
	}
}

func (o *OrderInfoFactory) Create(userId, userName, userEmail, userTelNo, memo, pickupDateTime string, stockOrders, foodOrders []ItemOrder) (*OrderInfo, error) {
	stocks, err := o.createOrderStockItems(stockOrders)
	if err != nil {
		return nil, err
	}

	foods, err := o.createOrderFoodItems(foodOrders)
	if err != nil {
		return nil, err
	}
	return NewOrderInfo(userId, userName, userEmail, userTelNo, memo, pickupDateTime, stocks, foods)
}

func (o *OrderInfoFactory) createOrderStockItems(stockOrders []ItemOrder) ([]OrderStockItem, error) {
	stocks, err := o.stockRepo.FindAll()
	if err != nil {
		return nil, err
	}
	stockItems := []OrderStockItem{}
	for _, stockOrder := range stockOrders {
		for _, stock := range stocks {
			if stock.HasSameId(stockOrder.id) {
				// check order max at first
				err = stock.WithInMaxOrder(stockOrder.quantity)
				if err != nil {
					return nil, err
				}
				item, err := NewOrderStockItem(stock.GetId(), stock.GetName(), stock.GetPrice(), stockOrder.quantity)
				if err != nil {
					return nil, err
				}
				if err != nil {
					return nil, err
				}
				// stock item remain check and update will be done at next step of usecase (consumer)
				stockItems = append(stockItems, *item)
				break
			}
		}
	}
	// check all items are existed
	if len(stockItems) != len(stockOrders) {
		return nil, common.NewValidationError("stockOrders", "there is not match stock item from id.")
	}
	return stockItems, nil
}

func (o *OrderInfoFactory) createOrderFoodItems(foodOrders []ItemOrder) ([]OrderFoodItem, error) {
	foods, err := o.foodRepo.FindAll()
	if err != nil {
		return nil, err
	}
	foodItems := []OrderFoodItem{}
	for _, foodOrder := range foodOrders {
		for _, food := range foods {
			if food.HasSameId(foodOrder.id) {
				// check order max at first
				err = food.WithInMaxOrder(foodOrder.quantity)
				if err != nil {
					return nil, err
				}
				item, err := NewOrderFoodItem(food.GetId(), food.GetName(), food.GetPrice(), foodOrder.quantity)
				if err != nil {
					return nil, err
				}
				foodItems = append(foodItems, *item)
				break
			}
		}
	}
	// check all items are existed
	if len(foodItems) != len(foodOrders) {
		return nil, common.NewValidationError("foodOrders", "there is not match food item from id.")
	}
	return foodItems, nil
}
