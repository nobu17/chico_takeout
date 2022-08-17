package item_test

import (
	"errors"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/item"
	"chico/takeout/tests"

	"github.com/stretchr/testify/assert"
)

type stockItemTest struct {
	name             string
	args             stockItemTestArgs
	want             stockItemTestArgs
	hasValidationErr bool
}

type stockItemTestArgs struct {
	name        string
	priority    int
	maxOrder    int
	price       int
	description string
	kindId      string
	enabled     bool
	remain      int
	imageUrl    string
}

var maxDescStr = tests.MakeRandomStr(150)
var maxName = tests.MakeRandomStr(15)

var stockItemInputs = []stockItemTest{
	{name: "normal case1",
		args:             stockItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{name: "test1", priority: 1, maxOrder: 2, price: 1, description: maxDescStr, kindId: "123", enabled: true, imageUrl: "http://google.com"},
		hasValidationErr: false,
	},
	{name: "normal case2",
		args:             stockItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 20000, description: "ttt", kindId: "abc", enabled: false, imageUrl: ""},
		want:             stockItemTestArgs{name: maxName, priority: 3, maxOrder: 30, price: 20000, description: "ttt", kindId: "abc", enabled: false, imageUrl: ""},
		hasValidationErr: false,
	},
	{name: "error:empty name",
		args:             stockItemTestArgs{name: "", priority: 3, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular name(length 16)",
		args:             stockItemTestArgs{name: tests.MakeRandomStr(16), priority: 3, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular priority(0)",
		args:             stockItemTestArgs{name: "test2", priority: 0, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular priority(-1)",
		args:             stockItemTestArgs{name: "test2", priority: -1, maxOrder: 4, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(0)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 0, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(-1)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: -1, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular maxOrder(31)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 31, price: 140, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular price(0)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 0, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular price(-1)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: -1, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular price(20001)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "ttt", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular description(empty)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "", kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular description(over150 length)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: tests.MakeRandomStr(151), kindId: "abc", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular kindId(empty)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 20001, description: "123", kindId: "", enabled: false, imageUrl: "http://google.com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error:irregular imageUrl(not url format)",
		args:             stockItemTestArgs{name: "test2", priority: 1, maxOrder: 20, price: 100, description: "ttt", kindId: "abc", enabled: false, imageUrl: "com"},
		want:             stockItemTestArgs{},
		hasValidationErr: true,
	},
}

func TestNewStockItem(t *testing.T) {
	for _, tt := range stockItemInputs {
		got, err := item.NewStockItem(tt.args.name, tt.args.description, tt.args.priority, tt.args.maxOrder, tt.args.price, tt.args.kindId, tt.args.enabled, tt.args.imageUrl)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewStockItem() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("New() should have error")
			return
		}
		expect := tt.want
		assert.Equal(t, expect.name, got.GetName())
		assert.Equal(t, expect.priority, got.GetPriority())
		assert.Equal(t, expect.maxOrder, got.GetMaxOrder())
		assert.Equal(t, expect.price, got.GetPrice())
		assert.Equal(t, expect.description, got.GetDescription())
		assert.Equal(t, expect.kindId, got.GetKindId())
		assert.Equal(t, expect.enabled, got.GetEnabled())
		assert.Equal(t, 0, got.GetRemain()) // remain is always 0
		assert.Equal(t, expect.imageUrl, got.GetImageUrl())
	}
}

func TestSetStockItem(t *testing.T) {
	for _, tt := range stockItemInputs {

		got, err := item.NewStockItem("test", "desc", 4, 4, 12, "123", false, "https://yahoo.com")
		if err != nil {
			t.Errorf("init test is failed")
			continue
		}
		err = got.Set(tt.args.name, tt.args.description, tt.args.priority, tt.args.maxOrder, tt.args.price, tt.args.kindId, tt.args.enabled, tt.args.imageUrl)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewStockItem() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
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
		assert.Equal(t, 0, got.GetRemain()) // remain is always 0
		assert.Equal(t, expect.imageUrl, got.GetImageUrl())
	}
}

func TestSetRemain(t *testing.T) {
	var setRemainInput = []stockItemTest{
		{name: "normal",
			args:             stockItemTestArgs{remain: 10},
			want:             stockItemTestArgs{remain: 10},
			hasValidationErr: false,
		},
		{name: "normal (0)",
			args:             stockItemTestArgs{remain: 0},
			want:             stockItemTestArgs{remain: 0},
			hasValidationErr: false,
		},
		{name: "validation error(-1)",
			args:             stockItemTestArgs{remain: -1},
			hasValidationErr: true,
		},
	}

	for _, tt := range setRemainInput {
		// arrange
		got, err := item.NewStockItem("test", "desc", 4, 4, 12, "123", false, "https://yahoo.com")
		// ensure no error
		assert.NoError(t, err, "init test is failed")

		// act
		err = got.SetRemain(tt.args.remain)
		if tt.hasValidationErr {
			assert.Error(t, err)
			assert.IsType(t, err, common.NewValidationError("", ""))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.want.remain, got.GetRemain())
	}
}

func TestConsumeRemain(t *testing.T) {
	var setRemainInput = []stockItemTest{
		{name: "normal",
			args:             stockItemTestArgs{remain: 1},
			want:             stockItemTestArgs{remain: 3},
			hasValidationErr: false,
		},
		{name: "normal (0)",
			args:             stockItemTestArgs{remain: 4},
			want:             stockItemTestArgs{remain: 0},
			hasValidationErr: false,
		},
		{name: "validation error(remain is over stock)",
			args:             stockItemTestArgs{remain: 5},
			hasValidationErr: true,
		},
		{name: "validation error(input is 0)",
			args:             stockItemTestArgs{remain: 0},
			hasValidationErr: true,
		},
		{name: "validation error(input is -1)",
			args:             stockItemTestArgs{remain: -1},
			hasValidationErr: true,
		},
	}

	for _, tt := range setRemainInput {
		// arrange init stock
		got, err := item.NewStockItem("test", "desc", 4, 4, 12, "123", false, "https://yahoo.com")
		assert.NoError(t, err, "init test is failed")
		// arrange initial stock (4)
		got.SetRemain(4)
		assert.NoError(t, err, "init remain test is failed")

		// act
		err = got.ConsumeRemain(tt.args.remain)
		if tt.hasValidationErr {
			assert.Error(t, err)
			assert.IsType(t, err, common.NewValidationError("", ""))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.want.remain, got.GetRemain())
	}
}

func TestIncrementRemain(t *testing.T) {
	var setRemainInput = []stockItemTest{
		{name: "normal",
			args:             stockItemTestArgs{remain: 1},
			want:             stockItemTestArgs{remain: 5},
			hasValidationErr: false,
		},
		{name: "normal",
			args:             stockItemTestArgs{remain: 4},
			want:             stockItemTestArgs{remain: 8},
			hasValidationErr: false,
		},
		{name: "normal (remain is just max(999))",
			args:             stockItemTestArgs{remain: 995},
			want:             stockItemTestArgs{remain: 999},
			hasValidationErr: false,
		},
		{name: "validation error(remain is over(999))",
			args:             stockItemTestArgs{remain: 996},
			hasValidationErr: true,
		},
		{name: "validation error(input is 0)",
			args:             stockItemTestArgs{remain: 0},
			hasValidationErr: true,
		},
		{name: "validation error(input is -1)",
			args:             stockItemTestArgs{remain: -1},
			hasValidationErr: true,
		},
	}

	for _, tt := range setRemainInput {
		// arrange init stock
		got, err := item.NewStockItem("test", "desc", 4, 4, 12, "123", false, "https://yahoo.com")
		assert.NoError(t, err, "init test is failed")
		// arrange initial stock (4)
		got.SetRemain(4)
		assert.NoError(t, err, "init remain test is failed")

		// act
		err = got.IncreaseRemain(tt.args.remain)
		if tt.hasValidationErr {
			assert.Error(t, err)
			assert.IsType(t, err, common.NewValidationError("", ""))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.want.remain, got.GetRemain())
	}
}
