package order

import (
	"strings"

	"chico/takeout/common"
	"chico/takeout/domains/item"

	"github.com/google/uuid"
)

type OrderInfoRepository interface {
	Find(id string) (*OrderInfo, error)
	FindByPickupDate(date string) ([]OrderInfo, error)
	FindByUserId(userId string) ([]OrderInfo, error)
	FindActiveByUserId(userId string) ([]OrderInfo, error)
	FindAll() ([]OrderInfo, error)
	Create(item *OrderInfo) (string, error)
	UpdateOrderStatus(item *OrderInfo) error
	Transact(fc func() error) (error)
}

const (
	OrderInfoMaxMemoLength = 500
	UserNameMaxLength = 10
)

type OrderInfo struct {
	id             string
	userId         string
	userName       UserName
	userEmail      Email
	userTelNo      TelNo
	memo           Memo
	orderDateTime  OrderDateTime
	pickupDateTime PickupDateTime
	stockItems     []OrderStockItem
	foodItems      []OrderFoodItem
	canceled       bool
}

func NewOrderInfo(userId, userName, userEmail, userTelNo, memo, pickupDateTime string, stockItems []OrderStockItem, foodItems []OrderFoodItem) (*OrderInfo, error) {
	order := &OrderInfo{id: uuid.NewString(), canceled: false}
	if err := order.validateUserId(userId); err != nil {
		return nil, err
	}
	memoV, err := NewMemo(memo, OrderInfoMaxMemoLength)
	if err != nil {
		return nil, err
	}
	userNameV, err := NewUserName(userName, UserNameMaxLength)
	if err != nil {
		return nil, err
	}
	userEmailV, err := NewEmail(userEmail)
	if err != nil {
		return nil, err
	}
	userTelNoV, err := NewTelNo(userTelNo)
	if err != nil {
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
	order.userName = *userNameV
	order.memo = *memoV
	order.userEmail = *userEmailV
	order.userTelNo = *userTelNoV
	order.orderDateTime = *orderDate
	order.pickupDateTime = *pickupDate
	order.stockItems = stockItems
	order.foodItems = foodItems
	return order, nil
}

func NewOrderInfoForOrm(id, userId, userName, userEmail, userTelNo, memo, pickupDateTime, orderDateTime string, stockItems []OrderStockItem, foodItems []OrderFoodItem, canceled bool) (*OrderInfo, error) {
	memoVal, _ := NewMemo(memo, OrderInfoMaxMemoLength)
	userNameV, _ := NewUserName(userTelNo, UserNameMaxLength)
	userEmailV, _ := NewEmail(userEmail)
	userTelNoV, _ := NewTelNo(userTelNo)
	order := &OrderInfo{
		id:             id,
		userId:         userId,
		userName:       *userNameV,
		userEmail:      *userEmailV,
		userTelNo:      *userTelNoV,
		memo:           *memoVal,
		pickupDateTime: PickupDateTime{},
		orderDateTime:  OrderDateTime{},
		stockItems:     stockItems,
		foodItems:      foodItems,
		canceled:       canceled,
	}

	order.pickupDateTime.value = pickupDateTime
	order.orderDateTime.value = orderDateTime

	return order, nil
}

func (o *OrderInfo) GetId() string {
	return o.id
}

func (o *OrderInfo) GetUserId() string {
	return o.userId
}

func (o *OrderInfo) GetUserName() string {
	return o.userName.GetValue()
}

func (o *OrderInfo) GetUserEmail() string {
	return o.userEmail.GetValue()
}

func (o *OrderInfo) GetUserTelNo() string {
	return o.userTelNo.GetValue()
}

func (o *OrderInfo) GetMemo() string {
	return o.memo.GetValue()
}

func (o *OrderInfo) GetCanceled() bool {
	return o.canceled
}

func (o *OrderInfo) GetPickupDateTime() string {
	return o.pickupDateTime.value
}

func (o *OrderInfo) GetOrderDateTime() string {
	return o.orderDateTime.value
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

type OrderStockItem struct {
	commonItemInfo
}

func NewOrderStockItem(itemId, name string, price, quantity int) (*OrderStockItem, error) {
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

func NewOrderFoodItem(itemId, name string, price, quantity int) (*OrderFoodItem, error) {
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
	name     item.Name
	price    Price
	quantity Quantity
}

func (c *commonItemInfo) HasSameId(id string) bool {
	return c.itemId == id
}

func (c *commonItemInfo) GetItemId() string {
	return c.itemId
}

func (c *commonItemInfo) GetName() string {
	return c.name.GetValue()
}

func (c *commonItemInfo) GetQuantity() int {
	return c.quantity.value
}

func (c *commonItemInfo) GetPrice() int {
	return c.price.value
}

func newCommonItemInfo(itemId, name string, price, quantity int) (*commonItemInfo, error) {
	if strings.TrimSpace(itemId) == "" {
		return nil, common.NewValidationError("itemId", "required")
	}

	nameV, err := item.NewName(name, item.CommonItemNameMaxLength)
	if err != nil {
		return nil, err
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
		name:     *nameV,
		price:    *priceV,
		quantity: *quantityV,
	}, nil
}
