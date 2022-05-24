package order

import (
	"fmt"
	"testing"

	"chico/takeout/common"
	"chico/takeout/tests"

	"github.com/stretchr/testify/assert"
)

type commonItemInfoArgs struct {
	itemId   string
	name     string
	price    int
	quantity int
}
type commonItemInfoInput struct {
	name             string
	args             commonItemInfoArgs
	want             commonItemInfoArgs
	hasValidationErr bool
	hasNotFoundErr   bool
}

func assertCommonItemInfoRoot(t *testing.T, tt commonItemInfoInput, got *commonItemInfo, err error) {
	if tt.hasValidationErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewValidationError("", ""), err)
		return
	}
	if tt.hasNotFoundErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewNotFoundError(""), err)
		return
	}
	assert.NoError(t, err, "no error should be")
	assertCommonItemInfo(t, tt.want, got)
}

func assertCommonItemInfo(t *testing.T, want commonItemInfoArgs, got *commonItemInfo) {
	assert.Equal(t, want.name, got.GetName())
	assert.Equal(t, want.itemId, got.GetItemId())
	assert.Equal(t, want.quantity, got.GetQuantity())
	assert.Equal(t, want.price, got.GetPrice())
}

func TestNewCommonItemInfo(t *testing.T) {
	inputs := []commonItemInfoInput{
		{name: "normal check",
			args: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10,
			},
			want: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10,
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(edge)",
			args: commonItemInfoArgs{
				name: "123456789012345", itemId: "12", price: 100, quantity: 10,
			},
			want: commonItemInfoArgs{
				name: "123456789012345", itemId: "12", price: 100, quantity: 10,
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "error empty name",
			args: commonItemInfoArgs{
				name: "", itemId: "12", price: 100, quantity: 10,
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error overlimit name(16)",
			args: commonItemInfoArgs{
				name: "1234567890123456", itemId: "12", price: 100, quantity: 10,
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error empty itemId",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "", price: 100, quantity: 10,
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error price(0)",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "12", price: 0, quantity: 10,
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error quantity(0)",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "12", price: 10, quantity: 0,
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)
		got, err := newCommonItemInfo(tt.args.itemId, tt.args.name, tt.args.price, tt.args.quantity)
		assertCommonItemInfoRoot(t, tt, got, err)
	}
}

type orderInfoArgs struct {
	id             string
	userId         string
	memo           string
	orderDateTime  string
	pickupDateTime string
	stockItems     []commonItemInfoArgs
	foodItems      []commonItemInfoArgs
	canceled       bool
}
type orderInfoInput struct {
	name             string
	args             orderInfoArgs
	want             orderInfoArgs
	hasValidationErr bool
	hasNotFoundErr   bool
}

func assertOderInfoRoot(t *testing.T, tt orderInfoInput, got *OrderInfo, err error) {
	if tt.hasValidationErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewValidationError("", ""), err)
		return
	}
	if tt.hasNotFoundErr {
		fmt.Println("err:", err)
		assert.Error(t, err, "should have error")
		assert.IsType(t, common.NewNotFoundError(""), err)
		return
	}
	assert.NoError(t, err, "no error should be")
	assertOrderInfo(t, tt.want, got)
}

func assertOrderInfo(t *testing.T, want orderInfoArgs, got *OrderInfo) {
	assert.Equal(t, want.userId, got.userId)
	assert.Equal(t, want.memo, got.GetMemo())
	assert.Equal(t, want.pickupDateTime, got.GetPickupDateTime())
	// assert.Equal(t, want.orderDateTime, got.GetOrderDateTime())
	assert.Equal(t, want.canceled, got.GetCanceled())
	assert.Equal(t, len(want.foodItems), len(got.GetFoodItems()))
	assert.Equal(t, len(want.stockItems), len(got.GetStockItems()))

	for index, got := range got.GetFoodItems() {
		assertCommonItemInfo(t, want.foodItems[index], &got.commonItemInfo)
	}

	for index, got := range got.GetStockItems() {
		assertCommonItemInfo(t, want.stockItems[index], &got.commonItemInfo)
	}
}

func TestNewOrderInfo(t *testing.T) {
	var maxMemo = tests.MakeRandomStr(500)
	inputs := []orderInfoInput{
		{name: "normal check",
			args: orderInfoArgs{
				userId: "abc", memo: "12", pickupDateTime: "2160/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			want: orderInfoArgs{
				userId: "abc", memo: "12", pickupDateTime: "2160/12/10 10:15", canceled: false,
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(min)",
			args: orderInfoArgs{
				userId: "a", memo: "", pickupDateTime: "2160/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			want: orderInfoArgs{
				userId: "a", memo: "", pickupDateTime: "2160/12/10 10:15", canceled: false,
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(max memo)",
			args: orderInfoArgs{
				userId: "a", memo: maxMemo, pickupDateTime: "2160/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			want: orderInfoArgs{
				userId: "a", memo: maxMemo, pickupDateTime: "2160/12/10 10:15", canceled: false,
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "empty userId",
			args: orderInfoArgs{
				userId: "", memo: "12", pickupDateTime: "2160/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "over limit memo(501)",
			args: orderInfoArgs{
				userId: "123", memo: tests.MakeRandomStr(501), pickupDateTime: "2160/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "pick up is before than now",
			args: orderInfoArgs{
				userId: "123", memo: "123", pickupDateTime: "2010/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "pick up is incorrect format",
			args: orderInfoArgs{
				userId: "123", memo: "123", pickupDateTime: "21001210 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "empty items",
			args: orderInfoArgs{
				userId: "123", memo: "123", pickupDateTime: "2020/12/10 10:15",
				stockItems: []commonItemInfoArgs{},
				foodItems:  []commonItemInfoArgs{},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		fOrders := []OrderFoodItem{}
		for _, food := range tt.args.foodItems {
			forder, err := NewOrderFoodItem(food.itemId, food.name, food.price, food.quantity)
			if err != nil {
				assert.Fail(t, "failed to init OrderFoodItem")
				return
			}
			fOrders = append(fOrders, *forder)
		}
		sOrders := []OrderStockItem{}
		for _, stock := range tt.args.stockItems {
			order, err := NewOrderStockItem(stock.itemId, stock.name, stock.price, stock.quantity)
			if err != nil {
				assert.Fail(t, "failed to init OrderStockItem")
				return
			}
			sOrders = append(sOrders, *order)
		}
		got, err := NewOrderInfo(tt.args.userId, tt.args.memo, tt.args.pickupDateTime, sOrders, fOrders)
		assertOderInfoRoot(t, tt, got, err)
	}
}

func TestOrderInfoSetCancel(t *testing.T) {
	inputs := []orderInfoInput{
		{name: "normal check",
			args: orderInfoArgs{
				userId: "abc", memo: "12", pickupDateTime: "2120/12/10 10:15",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			want: orderInfoArgs{
				userId: "abc", memo: "12", pickupDateTime: "2120/12/10 10:15", canceled: true,
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
					},
				},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		fOrders := []OrderFoodItem{}
		for _, food := range tt.args.foodItems {
			forder, err := NewOrderFoodItem(food.itemId, food.name, food.price, food.quantity)
			if err != nil {
				assert.Fail(t, "failed to init OrderFoodItem")
				return
			}
			fOrders = append(fOrders, *forder)
		}
		sOrders := []OrderStockItem{}
		for _, stock := range tt.args.stockItems {
			order, err := NewOrderStockItem(stock.itemId, stock.name, stock.price, stock.quantity)
			if err != nil {
				assert.Fail(t, "failed to init OrderStockItem")
				return
			}
			sOrders = append(sOrders, *order)
		}
		got, err := NewOrderInfo(tt.args.userId, tt.args.memo, tt.args.pickupDateTime, sOrders, fOrders)
		got.SetCancel()
		assertOderInfoRoot(t, tt, got, err)
	}
}
