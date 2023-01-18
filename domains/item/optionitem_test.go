package item_test

import (
	"errors"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/item"
	"chico/takeout/tests"
)

type optionItemTestArgs struct {
	name        string
	priority    int
	price       int
	description string
	enabled     bool
}

type optionItemTestInput struct {
	name             string
	args             optionItemTestArgs
	want             optionItemTestArgs
	hasValidationErr bool
}

var optionTestInputs = []optionItemTestInput{
	{name: "normal case1",
		args:             optionItemTestArgs{name: "test1", priority: 1, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{name: "test1", priority: 1, price: 100, description: "test", enabled: true},
		hasValidationErr: false,
	},
	{name: "normal case2",
		args:             optionItemTestArgs{name: "test2", priority: 2, price: 101, description: "a", enabled: false},
		want:             optionItemTestArgs{name: "test2", priority: 2, price: 101, description: "a", enabled: false},
		hasValidationErr: false,
	},
	{name: "normal case:name edge case(25char)",
		args:             optionItemTestArgs{name: "あいうえおかきくけこたちつてとあいうえおあいうえお", priority: 1, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{name: "あいうえおかきくけこたちつてとあいうえおあいうえお", priority: 1, price: 100, description: "test", enabled: true},
		hasValidationErr: false,
	},
	{name: "error case:name empty",
		args:             optionItemTestArgs{name: "", priority: 1, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:name over limit(26)",
		args:             optionItemTestArgs{name: "あいうえおかきくけこたちつてとあいうえおあいうえおA", priority: 1, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:priority is 0",
		args:             optionItemTestArgs{name: "aaa", priority: 0, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:priority is minus",
		args:             optionItemTestArgs{name: "aaa", priority: -1, price: 100, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "normal case:price edge case(20000)",
		args:             optionItemTestArgs{name: "aa", priority: 1, price: 20000, description: "test", enabled: true},
		want:             optionItemTestArgs{name: "aa", priority: 1, price: 20000, description: "test", enabled: true},
		hasValidationErr: false,
	},
	{name: "error case:price is over limit",
		args:             optionItemTestArgs{name: "aaa", priority: 1, price: 20001, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:price is 0",
		args:             optionItemTestArgs{name: "aaa", priority: 1, price: 0, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:price is minus",
		args:             optionItemTestArgs{name: "aaa", priority: 1, price: -1, description: "test", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:description is over limit(151)",
		args:             optionItemTestArgs{name: "aaa", priority: 1, price: 100, description: tests.MakeRandomStr(151), enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
	{name: "error case:description empty",
		args:             optionItemTestArgs{name: "aaa", priority: 1, price: 100, description: "", enabled: true},
		want:             optionItemTestArgs{},
		hasValidationErr: true,
	},
}

func TestNewOptionItem(t *testing.T) {
	for _, tt := range optionTestInputs {
		got, err := item.NewOptionItem(tt.args.name, tt.args.description, tt.args.priority, tt.args.price, tt.args.enabled)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewOptionItem() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("NewOptionItem() should have error")
			return
		}
		if got.GetName() != tt.want.name {
			t.Errorf("NewOptionItem() name should be:%s, actual:%s", tt.want.name, got.GetName())
			return
		}
		if got.GetDescription() != tt.want.description {
			t.Errorf("NewOptionItem() description should be:%s, actual:%s", tt.want.description, got.GetDescription())
			return
		}
		if got.GetPrice() != tt.want.price {
			t.Errorf("NewOptionItem() price should be:%d, actual:%d", tt.want.price, got.GetPrice())
			return
		}
		if got.GetPriority() != tt.want.priority {
			t.Errorf("NewOptionItem() priority should be:%d, actual:%d", tt.want.priority, got.GetPriority())
			return
		}
		if got.GetEnabled() != tt.want.enabled {
			t.Errorf("NewOptionItem() enabled should be:%v, actual:%v", tt.want.enabled, got.GetEnabled())
			return
		}
	}
}

func TestSetOptionItem(t *testing.T) {
	for _, tt := range optionTestInputs {
		got, err := item.NewOptionItem("test", "de", 12, 222, true);
		if err != nil {
			t.Errorf("NewOptionItem() unexpected error = %v", err)
			return
		}
		err = got.Set(tt.args.name, tt.args.description, tt.args.priority, tt.args.price, tt.args.enabled)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewOptionItem() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("NewOptionItem() should have error")
			return
		}
		if got.GetName() != tt.want.name {
			t.Errorf("NewOptionItem() name should be:%s, actual:%s", tt.want.name, got.GetName())
			return
		}
		if got.GetDescription() != tt.want.description {
			t.Errorf("NewOptionItem() description should be:%s, actual:%s", tt.want.description, got.GetDescription())
			return
		}
		if got.GetPrice() != tt.want.price {
			t.Errorf("NewOptionItem() price should be:%d, actual:%d", tt.want.price, got.GetPrice())
			return
		}
		if got.GetPriority() != tt.want.priority {
			t.Errorf("NewOptionItem() priority should be:%d, actual:%d", tt.want.priority, got.GetPriority())
			return
		}
		if got.GetEnabled() != tt.want.enabled {
			t.Errorf("NewOptionItem() enabled should be:%v, actual:%v", tt.want.enabled, got.GetEnabled())
			return
		}
	}
}