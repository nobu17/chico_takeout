package order_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	domains "chico/takeout/domains/order"
	"chico/takeout/infrastructures/memory"
	"chico/takeout/usecase/order"

	"github.com/stretchr/testify/assert"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func TestNotifyDailyOrder_NoOrderDate(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	date := time.Date(2020, 7, 19, 10, 0, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	// ensure mail is not sent
	assert.Nil(t, mail.DummyData)
}

func TestNotifyDailyOrder_NoOrderTime(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// order has until 12:00
	date := time.Date(2050, 12, 10, 12, 1, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	// ensure mail is not sent
	assert.Nil(t, mail.DummyData)
}

func TestNotifyDailyOrder_CanceledOrder_NotSent(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// 20250/12/11 has only canceled order
	date := time.Date(2050, 12, 11, 07, 1, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	// ensure mail is not sent
	assert.Nil(t, mail.DummyData)
}

func TestNotifyDailyOrder_SentMail(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	date := time.Date(2050, 12, 10, 1, 0, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	assert.Equal(t, "from@dummy.co.jp", mail.DummyData.SendFrom)
	assert.Equal(t, []string{"admin@dummy.co.jp"}, mail.DummyData.SendTo)
	assert.Equal(t, "", mail.DummyData.Bcc)
	assert.Equal(t, "本日のオーダー情報(2050/12/10)", mail.DummyData.Title)
	assert.Equal(t, true, strings.Contains(mail.DummyData.Message, "注文数:2"))
}

func setUpUseCase() (order.OrderTaskUseCase, *memory.MemorySendOrderMail) {
	repo := memory.NewOrderInfoMemoryRepository()
	orders := repo.GetMemory()
	// create additional order
	foodOrders1 := []domains.OrderFoodItem{}
	foodOrder1, err := domains.NewOrderFoodItem("abc", "food123", 210, 3)
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders1 = append(foodOrders1, *foodOrder1)

	stockOrders1 := []domains.OrderStockItem{}
	order1, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 10:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orders[order1.GetId()] = order1

	// create additional canceled order
	order2, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 11:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	order2.SetCancel()
	orders[order2.GetId()] = order2

	order3, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/11 11:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	order3.SetCancel()
	orders[order3.GetId()] = order3

	mail := memory.NewMemorySendOrderMail()
	useCase := order.NewOrderTaskUseCase(repo, mail)

	return useCase, mail
}

func setUpEnv(t *testing.T) {
	t.Setenv("MAIL_FROM", "from@dummy.co.jp")
	t.Setenv("MAIL_ADMIN", "admin@dummy.co.jp")
}
