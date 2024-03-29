package item_test

import (
	"errors"
	"fmt"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/item"
	"chico/takeout/tests"

	"github.com/stretchr/testify/assert"
)

type foodItemTest struct {
	name             string
	args             foodItemTestArgs
	want             foodItemTestArgs
	hasValidationErr bool
}

type foodItemTestArgs struct {
	name           string
	priority       int
	maxOrder       int
	price          int
	description    string
	kindId         string
	enabled        bool
	scheduleIds    []string
	maxOrderPerDay int
	imageUrl       string
	allowDates     []string
}

var foodItemInputs = []foodItemTest{
	{name: "normal case1",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: false,
	},
	{name: "normal case2",
		args:             foodItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 20000, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1"}, maxOrderPerDay: 100, imageUrl: "", allowDates: []string{}},
		want:             foodItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 20000, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1"}, maxOrderPerDay: 100, imageUrl: "", allowDates: []string{}},
		hasValidationErr: false,
	},
	{name: "normal case3(price zero",
		args:             foodItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 0, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1"}, maxOrderPerDay: 100, imageUrl: "", allowDates: []string{}},
		want:             foodItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 0, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1"}, maxOrderPerDay: 100, imageUrl: "", allowDates: []string{}},
		hasValidationErr: false,
	},
	{name: "normal case3(allow dates single)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{"2022/12/11"}},
		want:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{"2022/12/11"}},
		hasValidationErr: false,
	},
	{name: "normal case3(allow dates multiple)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{"2022/12/11", "2022/12/13"}},
		want:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 1, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{"2022/12/11", "2022/12/13"}},
		hasValidationErr: false,
	},
	{name: "error:empty name",
		args:             foodItemTestArgs{name: "", priority: 3, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "https://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular name(length 26)",
		args:             foodItemTestArgs{name: tests.MakeRandomStr(26), priority: 3, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular priority(0)",
		args:             foodItemTestArgs{name: "test2", priority: 0, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular priority(-1)",
		args:             foodItemTestArgs{name: "test2", priority: -1, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(0)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 0, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(-1)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: -1, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(31)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 31, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular price(-1)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: -1, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular price(20001)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular description(empty)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular description(over150 length)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: tests.MakeRandomStr(151), kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular kindId(empty)",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "123", kindId: "", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 30, imageUrl: "http://google.com", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error: irregular maxOrderPerDay(0)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 0, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error: irregular maxOrderPerDay(-1)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: -1, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error: irregular maxOrderPerDay(101)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 101, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error: irregular scheduleIds(empty)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{}, maxOrderPerDay: 10, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error: irregular scheduleIds(duplicated)",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2", "1"}, maxOrderPerDay: 10, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error: irregular maxOrder > maxOrderPerDay",
		args:             foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 1, imageUrl: "http://google.com", allowDates: []string{}},
		hasValidationErr: true,
	},
	{name: "error:not url format ImageUrl",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 30, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 10, imageUrl: "incorrect", allowDates: []string{}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:not url format incorrect date format",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 3, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 10, imageUrl: "http://google.com", allowDates: []string{"1234"}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:not url format incorrect date format",
		args:             foodItemTestArgs{name: "test2", priority: 1, maxOrder: 3, price: 140, description: "ttt", kindId: "abc", enabled: false, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 10, imageUrl: "http://google.com", allowDates: []string{"2012/12"}},
		want:             foodItemTestArgs{},
		hasValidationErr: true,
	},
}

func TestNewFoodItem(t *testing.T) {
	for _, tt := range foodItemInputs {
		fmt.Println("name:", tt.name)
		got, err := item.NewFoodItem(tt.args.name, tt.args.description, tt.args.priority, tt.args.maxOrder, tt.args.maxOrderPerDay, tt.args.price, tt.args.kindId, tt.args.scheduleIds, tt.args.enabled, tt.args.imageUrl, tt.args.allowDates)
		if err != nil {
			fmt.Println(err)
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewFoodItem() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			continue
		}
		if tt.hasValidationErr {
			t.Errorf("New() should have error")
			continue
		}
		expect := tt.want
		assert.Equal(t, expect.name, got.GetName())
		assert.Equal(t, expect.priority, got.GetPriority())
		assert.Equal(t, expect.maxOrder, got.GetMaxOrder())
		assert.Equal(t, expect.price, got.GetPrice())
		assert.Equal(t, expect.description, got.GetDescription())
		assert.Equal(t, expect.kindId, got.GetKindId())
		assert.Equal(t, expect.enabled, got.GetEnabled())
		assert.ElementsMatch(t, expect.scheduleIds, got.GetScheduleIds())
		assert.Equal(t, expect.maxOrderPerDay, got.GetMaxOrderPerDay())
		assert.Equal(t, expect.imageUrl, got.GetImageUrl())
		assert.ElementsMatch(t, expect.allowDates, got.GetAllowDates())
	}
}

func TestFoodItem_Update(t *testing.T) {
	for _, tt := range foodItemInputs {
		fmt.Println("name:", tt.name)

		// arrange
		init := foodItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, scheduleIds: []string{"1", "2"}, maxOrderPerDay: 4, imageUrl: "http://ho.com", allowDates: []string{}}
		got, err := item.NewFoodItem(init.name, init.description, init.priority, init.maxOrder, init.maxOrderPerDay, init.price, init.kindId, init.scheduleIds, init.enabled, init.imageUrl, init.allowDates)
		if err != nil {
			assert.Fail(t, "failed to initialize", err)
			continue
		}
		// act
		err = got.Set(tt.args.name, tt.args.description, tt.args.priority, tt.args.maxOrder, tt.args.maxOrderPerDay, tt.args.price, tt.args.kindId, tt.args.scheduleIds, tt.args.enabled, tt.args.imageUrl, tt.args.allowDates)
		if err != nil {
			fmt.Println(err)
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			assert.Fail(t, "unexpected error is raised.", err)
			continue
		}
		if tt.hasValidationErr {
			assert.Fail(t, "should has error.")
			continue
		}
		// assert
		expect := tt.want
		assert.Equal(t, expect.name, got.GetName())
		assert.Equal(t, expect.priority, got.GetPriority())
		assert.Equal(t, expect.maxOrder, got.GetMaxOrder())
		assert.Equal(t, expect.price, got.GetPrice())
		assert.Equal(t, expect.description, got.GetDescription())
		assert.Equal(t, expect.kindId, got.GetKindId())
		assert.Equal(t, expect.enabled, got.GetEnabled())
		assert.ElementsMatch(t, expect.scheduleIds, got.GetScheduleIds())
		assert.Equal(t, expect.maxOrderPerDay, got.GetMaxOrderPerDay())
		assert.Equal(t, expect.imageUrl, got.GetImageUrl())
		assert.ElementsMatch(t, expect.allowDates, got.GetAllowDates())
	}
}
