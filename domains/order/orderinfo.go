package order

import (
	"fmt"
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
	UpdateUserInfo(item *OrderInfo) error
	Transact(fc func() error) error
}

const (
	OrderInfoMaxMemoLength = 500
	UserNameMaxLength      = 10
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
	userNameV, _ := NewUserName(userName, UserNameMaxLength)
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
	pD, _ := NewDateTime(pickupDateTime)
	order.pickupDateTime.DateTime = *pD

	oD, _ := NewDateTime(orderDateTime)
	order.orderDateTime.DateTime = *oD

	return order, nil
}

func (o *OrderInfo) UpdateUserInfo(userName, userEmail, userTelNo, memo string) error {
	memoV, err := NewMemo(memo, OrderInfoMaxMemoLength)
	if err != nil {
		return err
	}
	userNameV, err := NewUserName(userName, UserNameMaxLength)
	if err != nil {
		return err
	}
	userEmailV, err := NewEmail(userEmail)
	if err != nil {
		return err
	}
	userTelNoV, err := NewTelNo(userTelNo)
	if err != nil {
		return err
	}

	o.memo = *memoV
	o.userName = *userNameV
	o.userEmail = *userEmailV
	o.userTelNo = *userTelNoV
	return nil
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

func (o *OrderInfo) GetTotalCost() int {
	total := 0

	for _, food := range o.foodItems {
		total += food.GetTotalCost()
	}
	for _, stock := range o.stockItems {
		total += stock.GetTotalCost()
	}

	return total
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

func NewOrderStockItem(itemId, name string, price, quantity int, options []OptionItemInfo) (*OrderStockItem, error) {
	item, err := newCommonItemInfo(itemId, name, price, quantity, options)
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

func NewOrderFoodItem(itemId, name string, price, quantity int, options []OptionItemInfo) (*OrderFoodItem, error) {
	item, err := newCommonItemInfo(itemId, name, price, quantity, options)
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
	options  []OptionItemInfo
}

type OptionItemInfo struct {
	itemId string
	name   item.Name
	price  Price
}

func (o *OptionItemInfo) GetId() string {
	return o.itemId
}

func (o *OptionItemInfo) GetName() string {
	return o.name.GetValue()
}

func (o *OptionItemInfo) GetPrice() int {
	return o.price.value
}

func NewOptionItemInfo(itemId, name string, price int) (*OptionItemInfo, error) {
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
	return &OptionItemInfo{
		itemId: itemId,
		name:   *nameV,
		price:  *priceV,
	}, nil
}

func (c *commonItemInfo) HasSameId(id string) bool {
	return c.itemId == id
}

func (c *commonItemInfo) GetOptionItems() []OptionItemInfo {
	return c.options
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

func (c *commonItemInfo) GetTotalCost() int {
	return c.price.value * c.quantity.value
}

func newCommonItemInfo(itemId, name string, price, quantity int, options []OptionItemInfo) (*commonItemInfo, error) {
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

	err = checkOptionItem(options)
	if err != nil {
		return nil, err
	}
	
	return &commonItemInfo{
		itemId:   itemId,
		name:     *nameV,
		price:    *priceV,
		quantity: *quantityV,
		options:  options,
	}, nil
}

func checkOptionItem(options []OptionItemInfo) error {
	ids := []string{}
	for _, opt := range options {
		opt.GetId()
		for _, id := range ids {
			if id == opt.GetId() {
				return common.NewValidationError("OptionItemInfo.Id", fmt.Sprintf("duplicate Id:%s", id))			
			}
		}
		ids = append(ids, opt.GetId())
	}
	return nil
}