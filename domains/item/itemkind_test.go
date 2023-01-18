package item_test

import (
	"errors"
	"testing"

	"chico/takeout/common"
	"chico/takeout/domains/item"
)

type args struct {
	name          string
	priority      int
	optionItemIds []string
}
type input struct {
	name             string
	args             args
	want             args
	hasValidationErr bool
}

var inputs = []input{
	{name: "normal case",
		args:             args{name: "test1", priority: 1, optionItemIds: []string{}},
		want:             args{name: "test1", priority: 1, optionItemIds: []string{}},
		hasValidationErr: false,
	},
	{name: "normal case (has ids)",
		args:             args{name: "test1", priority: 1, optionItemIds: []string{"1", "2"}},
		want:             args{name: "test1", priority: 1, optionItemIds: []string{"1", "2"}},
		hasValidationErr: false,
	},
	{name: "normal case:name edge case(15char)",
		args:             args{name: "あいうえおかきくけこたちつてと", priority: 1, optionItemIds: []string{}},
		want:             args{name: "あいうえおかきくけこたちつてと", priority: 1, optionItemIds: []string{}},
		hasValidationErr: false,
	},
	{name: "error case:name empty",
		args:             args{name: "", priority: 1, optionItemIds: []string{}},
		hasValidationErr: true,
	},
	{name: "error case:name is over limit(16) 2byte",
		args:             args{name: "あいうえおかきくけこたちつてとな", priority: 1, optionItemIds: []string{}},
		hasValidationErr: true,
	},
	{name: "error case:name is over limit(16) 1byte",
		args:             args{name: "1234567890123455", priority: 1, optionItemIds: []string{}},
		hasValidationErr: true,
	},
	{name: "error case:priority is 0",
		args:             args{name: "", priority: 1, optionItemIds: []string{}},
		hasValidationErr: true,
	},
	{name: "error case:priority is minus",
		args:             args{name: "", priority: -1, optionItemIds: []string{}},
		hasValidationErr: true,
	},
	{name: "error case:duplicated ids",
		args:             args{name: "test1", priority: 1, optionItemIds: []string{"1", "2", "1"}},
		hasValidationErr: true,
	},
}

func TestNewItemKind(t *testing.T) {
	for _, tt := range inputs {
		got, err := item.NewItemKind(tt.args.name, tt.args.priority, tt.args.optionItemIds)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewItemKind() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("NewItemKind() should have error")
			return
		}
		if got.GetName() != tt.want.name {
			t.Errorf("NewItemKind() ItemKind name should be:%s, actual:%s", tt.want.name, got.GetName())
			return
		}
		if got.GetPriority() != tt.want.priority {
			t.Errorf("NewItemKind() ItemKind priority should be:%d, actual:%d", tt.want.priority, got.GetPriority())
			return
		}
	}
}

func TestSetItemKind(t *testing.T) {
	for _, tt := range inputs {
		got, _ := item.NewItemKind("test", 4, []string{})
		err := got.Set(tt.args.name, tt.args.priority, tt.args.optionItemIds)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewItemKind() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("NewItemKind() should have error")
			return
		}
		if got.GetName() != tt.want.name {
			t.Errorf("NewItemKind() ItemKind name should be:%s, actual:%s", tt.want.name, got.GetName())
			return
		}
		if got.GetPriority() != tt.want.priority {
			t.Errorf("NewItemKind() ItemKind priority should be:%d, actual:%d", tt.want.priority, got.GetPriority())
			return
		}
	}
}
