package order

import (
	"chico/takeout/common"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type OrderInfoRepository interface {
	Find(id string) (*OrderInfo, error)
	FindByPickupDate(date string) ([]OrderInfo, error)
	FindAll() ([]OrderInfo, error)
	Create(item OrderInfo) (string, error)
	UpdateOrderStatus(item OrderInfo) error
}

type OrderInfo struct {
	id             string
	userId         string
	memo           string
	orderDateTime  OrderDateTime
	pickupDateTime PickupDateTime
	stockItems     []OrderStockItem
	foodItems      []OrderFoodItem
	canceled       bool
}

func NewOrderInfo(userId, memo, pickupDateTime string, stockItems []OrderStockItem, foodItems []OrderFoodItem) (*OrderInfo, error) {
	order := &OrderInfo{id: uuid.NewString(), canceled: false}
	if err := order.validateUserId(userId); err != nil {
		return nil, err
	}
	if err := order.validateMemo(memo); err != nil {
		return nil, err
	}
	// order date is current time
	orderDate, err := NewOrderDateTime()
	if err != nil {
		return nil, err
	}
	pickupDate, err := NewPickupDateTime(pickupDateTime)
	if err != nil {
		return nil, err
	}
	// if both items are empty, it is error
	if len(stockItems) == 0 && len(foodItems) == 0 {
		return nil, common.NewValidationError("stockItems and foodItems", "both items are empty")
	}

	order.userId = userId
	order.memo = memo
	order.orderDateTime = *orderDate
	order.pickupDateTime = *pickupDate
	order.stockItems = stockItems
	order.foodItems = foodItems
	return order, nil
}

func (o *OrderInfo) GetId() string {
	return o.id
}

func (o *OrderInfo) GetCanceled() bool {
	return o.canceled
}

func (o *OrderInfo) GetPickupDateTune() string {
	return o.pickupDateTime.value
}

func (o *OrderInfo) GetPickupDate() string {
	return o.pickupDateTime.GetAsDate()
}

func (o *OrderInfo) GetFoodItems() []OrderFoodItem {
	return o.foodItems
}

func (o *OrderInfo) GetStockItems() []OrderStockItem {
	return o.stockItems
}

func (o *OrderInfo) FindStockItemQuantity(itemId string) int {
	for _, item := range o.foodItems {
		if item.HasSameId(itemId) {
			return item.GetQuantity()
		}
	}
	return 0
}

func (o *OrderInfo) SetCancel() {
	o.canceled = true
}

func (o *OrderInfo) validateUserId(userId string) error {
	if strings.TrimSpace(userId) == "" {
		return common.NewValidationError("userId", "required")
	}
	return nil
}

const (
	OrderInfoMaxMemo = 580
)

func (o *OrderInfo) validateMemo(memo string) error {
	if len(memo) > OrderInfoMaxMemo {
		return common.NewValidationError("memo", fmt.Sprintf("MaxLength:%d", OrderInfoMaxMemo))
	}
	return nil
}

type OrderStockItem struct {
	commonItemInfo
}

func newOrderStockItem(itemId, name string, price, quantity int) (*OrderStockItem, error) {
	item, err := newCommonItemInfo(itemId, name, price, quantity)
	if err != nil {
		return nil, err
	}
	return &OrderStockItem{
		commonItemInfo: *item,
	}, nil
}

type OrderFoodItem struct {
	commonItemInfo
}

func newOrderFoodItem(itemId, name string, price, quantity int) (*OrderFoodItem, error) {
	item, err := newCommonItemInfo(itemId, name, price, quantity)
	if err != nil {
		return nil, err
	}
	return &OrderFoodItem{
		commonItemInfo: *item,
	}, nil
}

type commonItemInfo struct {
	itemId   string
	name     string
	price    Price
	quantity Quantity
}

func (c *commonItemInfo) HasSameId(id string) bool {
	if c.itemId == id {
		return true
	}
	return false
}

func (c *commonItemInfo) GetItemId() string {
	return c.itemId
}

func (c *commonItemInfo) GetQuantity() int {
	return c.quantity.value
}

func newCommonItemInfo(itemId, name string, price, quantity int) (*commonItemInfo, error) {
	if strings.TrimSpace(itemId) == "" {
		return nil, common.NewValidationError("itemId", "required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, common.NewValidationError("name", "required")
	}
	priceV, err := NewPrice(price)
	if err != nil {
		return nil, err
	}
	quantityV, err := NewQuantity(quantity)
	if err != nil {
		return nil, err
	}
	return &commonItemInfo{
		itemId:   itemId,
		price:    *priceV,
		quantity: *quantityV,
	}, nil
}
