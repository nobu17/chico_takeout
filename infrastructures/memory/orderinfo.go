package memory

import (
	"fmt"

	"chico/takeout/domains/item"
	domains "chico/takeout/domains/order"
)

var orderMemory map[string]*domains.OrderInfo
var orderStockItems []item.StockItem
var orderFoodItems []item.FoodItem

type OrderInfoMemoryRepository struct {
	inMemory   map[string]*domains.OrderInfo
	stockItems []item.StockItem
	foodItems  []item.FoodItem
}

func NewOrderInfoMemoryRepository() *OrderInfoMemoryRepository {
	if orderMemory == nil {
		resetOrderInfoMemory()
	}
	return &OrderInfoMemoryRepository{
		inMemory: orderMemory,
	}
}

func resetOrderInfoMemory() {
	NewItemKindMemoryRepository()
	NewBusinessHoursMemoryRepository()

	stockItemRepos := NewStockItemMemoryRepository()
	allStocks, _ := stockItemRepos.FindAll()

	foodItemRepos := NewFoodItemMemoryRepository()
	allFoods, _ := foodItemRepos.FindAll()

	orderMemory = map[string]*domains.OrderInfo{}

	// order1
	foodOrders1 := []domains.OrderFoodItem{}
	foodOrder1, err := domains.NewOrderFoodItem(allFoods[0].GetId(), allFoods[0].GetName(), allFoods[0].GetPrice(), 3)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders1 = append(foodOrders1, *foodOrder1)
	foodOrder2, err := domains.NewOrderFoodItem(allFoods[1].GetId(), allFoods[1].GetName(), allFoods[1].GetPrice(), 1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders1 = append(foodOrders1, *foodOrder2)

	stockOrders1 := []domains.OrderStockItem{}
	order1, err := domains.NewOrderInfo("user1", "user1@hoge.com", "123456789", "memo1", "2050/12/10 12:00", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orderMemory[order1.GetId()] = order1

	// order2
	foodOrders2 := []domains.OrderFoodItem{}
	foodOrder3, err := domains.NewOrderFoodItem(allFoods[0].GetId(), allFoods[0].GetName(), allFoods[0].GetPrice(), 1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders2 = append(foodOrders2, *foodOrder3)

	stockOrders2 := []domains.OrderStockItem{}
	stockOrder1, err := domains.NewOrderStockItem(allStocks[0].GetId(), allStocks[0].GetName(), allStocks[0].GetPrice(), 2)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	stockOrders2 = append(stockOrders2, *stockOrder1)
	order2, err := domains.NewOrderInfo("user2", "user2@hoge.com", "987654321", "memo2", "2050/12/14 12:00", stockOrders2, foodOrders2)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	orderMemory[order2.GetId()] = order2
}

func (o *OrderInfoMemoryRepository) GetMemory() map[string]*domains.OrderInfo {
	return o.inMemory
}

func (o *OrderInfoMemoryRepository) Reset() {
	resetOrderInfoMemory()
}

func (o *OrderInfoMemoryRepository) FindAll() ([]domains.OrderInfo, error) {
	items := []domains.OrderInfo{}
	for _, item := range o.inMemory {
		items = append(items, *item)
	}
	return items, nil
}

func (o *OrderInfoMemoryRepository) Find(id string) (*domains.OrderInfo, error) {
	if val, ok := o.inMemory[id]; ok {
		// need copy to protect
		duplicated := *val
		return &duplicated, nil
	}
	return nil, nil
}

func (o *OrderInfoMemoryRepository) FindByPickupDate(date string) ([]domains.OrderInfo, error) {
	items := []domains.OrderInfo{}
	for _, item := range o.inMemory {
		if item.GetPickupDate() == date {
			items = append(items, *item)
		}
	}
	return items, nil
}

func (o *OrderInfoMemoryRepository) Create(item *domains.OrderInfo) (string, error) {
	o.inMemory[item.GetId()] = item
	return item.GetId(), nil
}

func (o *OrderInfoMemoryRepository) UpdateOrderStatus(item *domains.OrderInfo) error {
	if _, ok := o.inMemory[item.GetId()]; ok {
		o.inMemory[item.GetId()] = item
		return nil
	}
	return fmt.Errorf("update target not exists")
}

func (o *OrderInfoMemoryRepository) Transact(fc func() error) error {
	return fc()
}
