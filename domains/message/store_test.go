package message_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"chico/takeout/common"
	"chico/takeout/domains/message"

	"github.com/stretchr/testify/assert"
)

type messageInfoArg struct {
	id      string
	content string
	edited  time.Time
	created time.Time
}

type messageInfoTestInput struct {
	name             string
	args             messageInfoArg
	want             messageInfoArg
	hasValidationErr bool
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func TestNewStoreTopMessage(t *testing.T) {
	mockedTime := time.Date(2023, 10, 5, 10, 20, 30, 0, jst)
	common.MockNow(func() time.Time {
		return mockedTime
	})
	inputs := []messageInfoTestInput{
		{
			name:             "normal",
			args:             messageInfoArg{content: "msg1", created: mockedTime, edited: mockedTime},
			want:             messageInfoArg{id: "1", content: "msg1", created: mockedTime, edited: mockedTime},
			hasValidationErr: false,
		},
		{
			name:             "empty error",
			args:             messageInfoArg{content: "", created: mockedTime, edited: mockedTime},
			hasValidationErr: true,
		},
	}
	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		msg, err := message.NewMessage("1", tt.args.content)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewStoreTopMessage() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("New() should have error")
			return
		}
		assert.NoError(t, err, "no error should be")
		assert.Equal(t, tt.want.id, msg.GetId())
		assert.Equal(t, tt.want.content, msg.GetContent())
		assert.Equal(t, tt.want.created, msg.GetCreated())
		assert.Equal(t, tt.want.edited, msg.GetEdited())
	}
}

func TestNewMyPageMessage(t *testing.T) {
	mockedTime := time.Date(2023, 10, 5, 10, 20, 30, 0, jst)
	common.MockNow(func() time.Time {
		return mockedTime
	})
	inputs := []messageInfoTestInput{
		{
			name:             "normal",
			args:             messageInfoArg{content: "msg1", created: mockedTime, edited: mockedTime},
			want:             messageInfoArg{id: "2", content: "msg1", created: mockedTime, edited: mockedTime},
			hasValidationErr: false,
		},
		{
			name:             "empty error",
			args:             messageInfoArg{content: "", created: mockedTime, edited: mockedTime},
			hasValidationErr: true,
		},
	}
	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		msg, err := message.NewMessage("2", tt.args.content)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewMyPageMessage() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("New() should have error")
			return
		}
		assert.NoError(t, err, "no error should be")
		assert.Equal(t, tt.want.id, msg.GetId())
		assert.Equal(t, tt.want.content, msg.GetContent())
		assert.Equal(t, tt.want.created, msg.GetCreated())
		assert.Equal(t, tt.want.edited, msg.GetEdited())
	}
}


func TestNotSupportedMessage(t *testing.T) {
	mockedTime := time.Date(2023, 10, 5, 10, 20, 30, 0, jst)
	common.MockNow(func() time.Time {
		return mockedTime
	})
	inputs := []messageInfoTestInput{
		{
			name:             "not supported id",
			args:             messageInfoArg{id: "3", content: "test", created: mockedTime, edited: mockedTime},
			hasValidationErr: true,
		},
	}
	for _, tt := range inputs {
		fmt.Println("name:", tt.name)

		msg, err := message.NewMessage(tt.args.id, tt.args.content)
		if err != nil {
			var vErr *common.ValidationError
			if errors.As(err, &vErr) {
				if tt.hasValidationErr {
					continue
				}
			}
			t.Errorf("NewMyPageMessage() error = %v, hasValidationErr %v", err, tt.hasValidationErr)
			return
		}
		if tt.hasValidationErr {
			t.Errorf("New() should have error")
			return
		}
		assert.NoError(t, err, "no error should be")
		assert.Equal(t, tt.want.id, msg.GetId())
		assert.Equal(t, tt.want.content, msg.GetContent())
		assert.Equal(t, tt.want.created, msg.GetCreated())
		assert.Equal(t, tt.want.edited, msg.GetEdited())
	}
}
