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
	options  []optionItemInfo
}

type optionItemInfo struct {
	itemId string
	name   string
	price  int
}

type optionItemInfoInput struct {
	name             string
	args             optionItemInfo
	want             optionItemInfo
	hasValidationErr bool
}

func (o *commonItemInfoArgs) toOptions() ([]OptionItemInfo, error) {
	opts := []OptionItemInfo{}
	for _, op := range o.options {
		opt, err := op.toOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, *opt)
	}
	return opts, nil
}

func (o *optionItemInfo) toOption() (*OptionItemInfo, error) {
	return NewOptionItemInfo(o.itemId, o.name, o.price)
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
	if len(want.options) != len(got.GetOptionItems()) {
		assert.Fail(t, "option length is not same.(got=%d, want=%d)", len(want.options), len(got.options))
	}
	gotOps := got.GetOptionItems()
	for i, wantOp := range want.options {
		assertOptionItemInfo(t, wantOp, gotOps[i])
	}
}

func assertOptionItemInfo(t *testing.T, want optionItemInfo, got OptionItemInfo) {
	assert.Equal(t, want.name, got.GetName())
	assert.Equal(t, want.itemId, got.GetId())
	assert.Equal(t, want.price, got.GetPrice())
}

func TestNewOptionItemInfo(t *testing.T) {
	inputs := []optionItemInfoInput{
		{
			name:             "normal",
			args:             optionItemInfo{itemId: "1", name: "opt1", price: 10},
			want:             optionItemInfo{itemId: "1", name: "opt1", price: 10},
			hasValidationErr: false,
		},
		{
			name:             "edge limit name(25)",
			args:             optionItemInfo{itemId: "1", name: "1234567890123456789012345", price: 10},
			want:             optionItemInfo{itemId: "1", name: "1234567890123456789012345", price: 10},
			hasValidationErr: false,
		},
		{
			name:             "empty item id",
			args:             optionItemInfo{itemId: "", name: "opt1", price: 10},
			want:             optionItemInfo{},
			hasValidationErr: true,
		},
		{
			name:             "empty name",
			args:             optionItemInfo{itemId: "1", name: "", price: 10},
			want:             optionItemInfo{},
			hasValidationErr: true,
		},
		{
			name:             "zero price",
			args:             optionItemInfo{itemId: "1", name: "opt1", price: 0},
			want:             optionItemInfo{},
			hasValidationErr: true,
		},
		{
			name:             "minus price",
			args:             optionItemInfo{itemId: "1", name: "opt1", price: -1},
			want:             optionItemInfo{},
			hasValidationErr: true,
		},
		{
			name:             "over limit name(26)",
			args:             optionItemInfo{itemId: "1", name: "12345678901234567890123456", price: 10},
			want:             optionItemInfo{},
			hasValidationErr: true,
		},
	}
	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		opts, err := tt.args.toOption()
		if tt.hasValidationErr {
			fmt.Println("err:", err)
			assert.Error(t, err, "should have error")
			assert.IsType(t, common.NewValidationError("", ""), err)
			continue
		}
		assert.NoError(t, err, "no error should be")
		assertOptionItemInfo(t, tt.want, *opts)
	}
}

func TestNewCommonItemInfo(t *testing.T) {
	inputs := []commonItemInfoInput{
		{name: "normal check",
			args: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			want: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(with option item)",
			args: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10,
				options: []optionItemInfo{
					{itemId: "1", name: "opt1", price: 10},
				},
			},
			want: commonItemInfoArgs{
				name: "test", itemId: "12", price: 100, quantity: 10,
				options: []optionItemInfo{
					{itemId: "1", name: "opt1", price: 10},
				},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(edge)",
			args: commonItemInfoArgs{
				name: "1234567890123456789012345", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			want: commonItemInfoArgs{
				name: "1234567890123456789012345", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "error empty name",
			args: commonItemInfoArgs{
				name: "", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error over limit name(26)",
			args: commonItemInfoArgs{
				name: "12345678901234567890123456", itemId: "12", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error empty itemId",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "", price: 100, quantity: 10, options: []optionItemInfo{},
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error price(0)",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "12", price: 0, quantity: 10, options: []optionItemInfo{},
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "error quantity(0)",
			args: commonItemInfoArgs{
				name: "123456789", itemId: "12", price: 10, quantity: 0, options: []optionItemInfo{},
			},
			want:             commonItemInfoArgs{},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		opts, err := tt.args.toOptions()
		assert.NoError(t, err)
		got, err := newCommonItemInfo(tt.args.itemId, tt.args.name, tt.args.price, tt.args.quantity, opts)
		assertCommonItemInfoRoot(t, tt, got, err)
	}
}

type orderInfoArgs struct {
	id             string
	userId         string
	userName       string
	userEmail      string
	userTelNo      string
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
	assert.Equal(t, want.userName, got.GetUserName())
	assert.Equal(t, want.userEmail, got.GetUserEmail())
	assert.Equal(t, want.userTelNo, got.GetUserTelNo())
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
				userId: "abc", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
				userId: "abc", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				canceled: false,
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
		{name: "normal check(with option item)",
			args: orderInfoArgs{
				userId: "abc", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
						options: []optionItemInfo{
							{itemId: "1", name: "opt1", price: 10},
							{itemId: "2", name: "opt2", price: 11},
						},
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
						options: []optionItemInfo{
							{itemId: "3", name: "opt3", price: 12},
							{itemId: "4", name: "opt4", price: 13},
						},
					},
				},
			},
			want: orderInfoArgs{
				userId: "abc", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				canceled: false,
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
						options: []optionItemInfo{
							{itemId: "1", name: "opt1", price: 10},
							{itemId: "2", name: "opt2", price: 11},
						},
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
						options: []optionItemInfo{
							{itemId: "3", name: "opt3", price: 12},
							{itemId: "4", name: "opt4", price: 13},
						},
					},
				},
			},
			hasValidationErr: false,
			hasNotFoundErr:   false,
		},
		{name: "normal check(min)",
			args: orderInfoArgs{
				userId: "a", userName: "u", memo: "", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "1",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			want: orderInfoArgs{
				userId: "a", userName: "u", memo: "", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "1", canceled: false,
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
				userId: "a", userName: "ユーザーABC", memo: maxMemo, pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
					},
				},
				foodItems: []commonItemInfoArgs{},
			},
			want: orderInfoArgs{
				userId: "a", userName: "ユーザーABC", memo: maxMemo, pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789", canceled: false,
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
				userId: "", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
		{name: "empty userName",
			args: orderInfoArgs{
				userId: "user", userName: "", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
		{name: "over limit userName",
			args: orderInfoArgs{
				userId: "user", userName: "12345678901", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
		{name: "empty email",
			args: orderInfoArgs{
				userId: "1234", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "", userTelNo: "123456789",
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
		{name: "incorrect email format",
			args: orderInfoArgs{
				userId: "1234", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "abc", userTelNo: "123456789",
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
		{name: "empty telNo",
			args: orderInfoArgs{
				userId: "1234", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "",
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
		{name: "incorrect telNo",
			args: orderInfoArgs{
				userId: "1234", userName: "ユーザーABC", memo: "12", pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "abc1234",
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
				userId: "123", userName: "ユーザーABC", memo: tests.MakeRandomStr(501), pickupDateTime: "2160/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
				userId: "123", userName: "ユーザーABC", memo: "123", pickupDateTime: "2010/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
				userId: "123", userName: "ユーザーABC", memo: "123", pickupDateTime: "21001210 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
				userId: "123", userName: "ユーザーABC", memo: "123", pickupDateTime: "2020/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				stockItems: []commonItemInfoArgs{},
				foodItems:  []commonItemInfoArgs{},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
		{name: "option stock item is duplicated",
			args: orderInfoArgs{
				userId: "123", userName: "ユーザーABC", memo: "123", pickupDateTime: "2120/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
						options: []optionItemInfo{
							{itemId: "1", name: "opt1", price: 10},
							{itemId: "2", name: "opt2", price: 11},
							{itemId: "2", name: "opt2", price: 11},
						},
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
		{name: "option food item is duplicated",
			args: orderInfoArgs{
				userId: "123", userName: "ユーザーABC", memo: "123", pickupDateTime: "2120/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
				stockItems: []commonItemInfoArgs{
					{
						name: "item1", itemId: "12", price: 100, quantity: 10,
						options: []optionItemInfo{},
					},
				},
				foodItems: []commonItemInfoArgs{
					{
						name: "item2", itemId: "13", price: 200, quantity: 1,
						options: []optionItemInfo{
							{itemId: "1", name: "opt1", price: 10},
							{itemId: "2", name: "opt2", price: 11},
							{itemId: "2", name: "opt2", price: 11},
						},
					},
				},
			},
			hasValidationErr: true,
			hasNotFoundErr:   false,
		},
	}

	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		fOrders := []OrderFoodItem{}
		for _, food := range tt.args.foodItems {
			opts, err := food.toOptions()
			assert.NoError(t, err)
			fOrder, err := NewOrderFoodItem(food.itemId, food.name, food.price, food.quantity, opts)
			if err != nil {
				assertOderInfoRoot(t, tt, nil, err)
				// assert.Fail(t, "failed to init OrderFoodItem")
				return
			}
			fOrders = append(fOrders, *fOrder)
		}
		sOrders := []OrderStockItem{}
		for _, stock := range tt.args.stockItems {
			opts, err := stock.toOptions()
			assert.NoError(t, err)
			order, err := NewOrderStockItem(stock.itemId, stock.name, stock.price, stock.quantity, opts)
			if err != nil {
				assertOderInfoRoot(t, tt, nil, err)
				// assert.Fail(t, "failed to init OrderStockItem")
				return
			}
			sOrders = append(sOrders, *order)
		}
		got, err := NewOrderInfo(tt.args.userId, tt.args.userName, tt.args.userEmail, tt.args.userTelNo, tt.args.memo, tt.args.pickupDateTime, sOrders, fOrders)
		assertOderInfoRoot(t, tt, got, err)
	}
}

func TestOrderInfoSetCancel(t *testing.T) {
	inputs := []orderInfoInput{
		{name: "normal check",
			args: orderInfoArgs{
				userId: "abc", userName: "name", memo: "12", pickupDateTime: "2120/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789",
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
				userId: "abc", userName: "name", memo: "12", pickupDateTime: "2120/12/10 10:15", userEmail: "user1@hoge.com", userTelNo: "123456789", canceled: true,
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

		foodOrders := []OrderFoodItem{}
		for _, food := range tt.args.foodItems {
			opts, err := food.toOptions()
			assert.NoError(t, err)
			foodOrder, err := NewOrderFoodItem(food.itemId, food.name, food.price, food.quantity, opts)
			if err != nil {
				assert.Fail(t, "failed to init OrderFoodItem")
				return
			}
			foodOrders = append(foodOrders, *foodOrder)
		}
		sOrders := []OrderStockItem{}
		for _, stock := range tt.args.stockItems {
			opts, err := stock.toOptions()
			assert.NoError(t, err)
			order, err := NewOrderStockItem(stock.itemId, stock.name, stock.price, stock.quantity, opts)
			if err != nil {
				assert.Fail(t, "failed to init OrderStockItem")
				return
			}
			sOrders = append(sOrders, *order)
		}
		got, err := NewOrderInfo(tt.args.userId, tt.args.userName, tt.args.userEmail, tt.args.userTelNo, tt.args.memo, tt.args.pickupDateTime, sOrders, foodOrders)
		if err == nil {
			got.SetCancel()
		}
		assertOderInfoRoot(t, tt, got, err)
	}
}
