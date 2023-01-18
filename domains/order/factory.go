package order

import (
	"chico/takeout/common"
	"chico/takeout/domains/item"
	"fmt"
)

type ItemOrder struct {
	id       string
	quantity int
	options  []OptionOrder
}

func (i *ItemOrder) GetOptionOderIds() []string {
	ids := []string{}
	for _, op := range i.options {
		ids = append(ids, op.id)
	}
	return ids
}

type OptionOrder struct {
	id string
}

func NewItemOrder(id string, quantity int, optionIds []string) *ItemOrder {
	options := []OptionOrder{}
	for _, opt := range optionIds {
		options = append(options, OptionOrder{id: opt})
	}
	return &ItemOrder{
		id:       id,
		quantity: quantity,
		options:  options,
	}
}

type OrderInfoFactory struct {
	stockRepo item.StockItemRepository
	foodRepo  item.FoodItemRepository
	kindRepo  item.ItemKindRepository
	optionRepo item.OptionItemRepository
}

func NewOrderInfoFactory(stockRepo item.StockItemRepository, foodRepo item.FoodItemRepository, kindRepo item.ItemKindRepository, optionRepo item.OptionItemRepository) *OrderInfoFactory {
	return &OrderInfoFactory{
		stockRepo: stockRepo,
		foodRepo:  foodRepo,
		kindRepo:  kindRepo,
		optionRepo: optionRepo,
	}
}

func (o *OrderInfoFactory) Create(userId, userName, userEmail, userTelNo, memo, pickupDateTime string, stockOrders, foodOrders []ItemOrder) (*OrderInfo, error) {
	optionItems, err := o.optionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	stocks, err := o.createOrderStockItems(stockOrders, optionItems)
	if err != nil {
		return nil, err
	}

	foods, err := o.createOrderFoodItems(foodOrders, optionItems)
	if err != nil {
		return nil, err
	}
	return NewOrderInfo(userId, userName, userEmail, userTelNo, memo, pickupDateTime, stocks, foods)
}

func (o *OrderInfoFactory) createOrderStockItems(stockOrders []ItemOrder, optionItems []item.OptionItem) ([]OrderStockItem, error) {
	stocks, err := o.stockRepo.FindAll()
	if err != nil {
		return nil, err
	}
	stockItems := []OrderStockItem{}
	for _, stockOrder := range stockOrders {
		for _, stock := range stocks {
			if stock.HasSameId(stockOrder.id) {
				// check option is exists
				err := o.checkOptionItemExists(stockOrder, stock.GetKindId())
				if err != nil {
					return nil, err
				}

				// check order max at first
				err = stock.WithInMaxOrder(stockOrder.quantity)
				if err != nil {
					return nil, err
				}

				opts, err := o.creteOptionItems(stockOrder.GetOptionOderIds(), optionItems)
				if err != nil {
					return nil, err
				}

				item, err := NewOrderStockItem(stock.GetId(), stock.GetName(), stock.GetPrice(), stockOrder.quantity, opts)
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

func (o *OrderInfoFactory) createOrderFoodItems(foodOrders []ItemOrder, optionItems []item.OptionItem) ([]OrderFoodItem, error) {
	foods, err := o.foodRepo.FindAll()
	if err != nil {
		return nil, err
	}
	foodItems := []OrderFoodItem{}
	for _, foodOrder := range foodOrders {
		for _, food := range foods {
			if food.HasSameId(foodOrder.id) {
				// check option is exists
				err := o.checkOptionItemExists(foodOrder, food.GetKindId())
				if err != nil {
					return nil, err
				}

				// check order max at first
				err = food.WithInMaxOrder(foodOrder.quantity)
				if err != nil {
					return nil, err
				}

				opts, err := o.creteOptionItems(foodOrder.GetOptionOderIds(), optionItems)
				if err != nil {
					return nil, err
				}
				item, err := NewOrderFoodItem(food.GetId(), food.GetName(), food.GetPrice(), foodOrder.quantity, opts)
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

func (o *OrderInfoFactory) checkOptionItemExists(item ItemOrder, kindId string) error {
	if len(item.options) == 0 {
		return nil
	}
	kind, err := o.kindRepo.Find(kindId)
	if err != nil {
		return err
	}
	for _, opt := range item.options {
		exists := false
		for _, id := range kind.GetOptionItemIds() {
			if opt.id == id {
				exists = true
				break
			}
		}
		if !exists {
			return common.NewValidationError("optionItemId", fmt.Sprintf("Not exists from kind:%s ", opt.id))
		}
	}
	return nil
}

func (o *OrderInfoFactory) creteOptionItems(optionIds []string, options []item.OptionItem) ([]OptionItemInfo, error) {
	results := []OptionItemInfo{}
	for _, id := range optionIds {
		for _, opt := range options {
			if opt.HasSameId(id) {
				op, err := NewOptionItemInfo(opt.GetId(), opt.GetName(), opt.GetPrice())
				if err != nil {
					return nil, err
				}
				results = append(results, *op)
			}
		}
	}
	return results, nil
}