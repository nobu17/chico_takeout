package order_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/order"
	stDomains "chico/takeout/domains/store"
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

	// ensure no-order mail is sent
	assert.Equal(t, "from@dummy.co.jp", mail.Sent[0].SendFrom)
	assert.Equal(t, []string{"admin@dummy.co.jp"}, mail.Sent[0].SendTo)
	assert.Equal(t, "", mail.Sent[0].Bcc)
	assert.Equal(t, "本日のオーダーはありません(2020/07/19)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "本日のオーダーはありません。"))
}

func TestNotifyDailyOrder_NoOrderTime(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// order has until 12:00
	date := time.Date(2050, 12, 10, 12, 1, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	// ensure no-order mail is sent
	assert.Equal(t, "from@dummy.co.jp", mail.Sent[0].SendFrom)
	assert.Equal(t, []string{"admin@dummy.co.jp"}, mail.Sent[0].SendTo)
	assert.Equal(t, "", mail.Sent[0].Bcc)
	assert.Equal(t, "本日のオーダーはありません(2050/12/10)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "本日のオーダーはありません。"))
}

func TestNotifyDailyOrder_CanceledOrder_NotSent(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// 20250/12/11 has only canceled order
	date := time.Date(2050, 12, 11, 07, 1, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	// ensure no-order mail is sent
	assert.Equal(t, "from@dummy.co.jp", mail.Sent[0].SendFrom)
	assert.Equal(t, []string{"admin@dummy.co.jp"}, mail.Sent[0].SendTo)
	assert.Equal(t, "", mail.Sent[0].Bcc)
	assert.Equal(t, "本日のオーダーはありません(2050/12/11)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "本日のオーダーはありません。"))
}

func TestNotifyDailyOrder_SentMail(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	date := time.Date(2050, 12, 10, 1, 0, 0, 0, jst)
	useCase.NotifyDailyOrder(date)

	assert.Equal(t, "from@dummy.co.jp", mail.Sent[0].SendFrom)
	assert.Equal(t, []string{"admin@dummy.co.jp"}, mail.Sent[0].SendTo)
	assert.Equal(t, "", mail.Sent[0].Bcc)
	assert.Equal(t, "本日のオーダー情報(2050/12/10)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "注文数:1"))
}

func TestNotifyOrderByHour_NotSent_TimeIsBefore(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// morning:7:00 - 2:00 => 5:00 ~ 5:25
	current := time.Date(2050, 12, 10, 4, 10, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// edge time
	current = time.Date(2050, 12, 10, 4, 59, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55
	current = time.Date(2050, 12, 10, 9, 29, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25
	current = time.Date(2050, 12, 10, 15, 59, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))
}

func TestNotifyOrderByHour_NotSent_TimeIsPassed(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// morning:7:00 - 2:00 => 5:00 ~ 5:25 => 5:26
	current := time.Date(2050, 12, 10, 5, 26, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	current = time.Date(2050, 12, 10, 6, 0, 0, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55
	current = time.Date(2050, 12, 10, 9, 56, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25
	current = time.Date(2050, 12, 10, 16, 26, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))
}

func TestNotifyOrderByHour_NotSent_HourIsNotIncluded(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25 but Sun is not included
	current := time.Date(2050, 12, 11, 16, 5, 0, 0, jst) // sun
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))
}

func TestNotifyOrderByHour_Sent_WithOrder(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// morning:7:00 - 2:00 => 5:00 ~ 5:25 => 5:00
	current := time.Date(2050, 12, 10, 5, 00, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/10 07:00 ~ 09:30)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "注文数:1"))

	// 5:01
	current = time.Date(2050, 12, 10, 5, 1, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/10 07:00 ~ 09:30)", mail.Sent[1].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[1].Message, "注文数:1"))

	// 5:24.59.999
	current = time.Date(2050, 12, 10, 5, 24, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 3, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/10 07:00 ~ 09:30)", mail.Sent[2].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[2].Message, "注文数:1"))

	// 5:25.59.999
	current = time.Date(2050, 12, 10, 5, 25, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 4, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/10 07:00 ~ 09:30)", mail.Sent[3].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[3].Message, "注文数:1"))
}

func TestNotifyOrderByHour_Sent_WithMultipleOrder(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCaseWithMultipleOrders()

	// morning:7:00 - 2:00 => 5:00 ~ 5:25
	current := time.Date(2050, 12, 10, 5, 00, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/10 07:00 ~ 09:30)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "注文数:2"))
}

func TestNotifyOrderByHour_Sent_NoOrder(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25 => 16:00
	current := time.Date(2050, 12, 10, 16, 00, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 18:00 ~ 21:00)", mail.Sent[0].Title)

	// 16:01
	current = time.Date(2050, 12, 10, 16, 1, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 18:00 ~ 21:00)", mail.Sent[1].Title)

	// 16:24
	current = time.Date(2050, 12, 10, 16, 24, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 3, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 18:00 ~ 21:00)", mail.Sent[2].Title)

	// 16:25
	current = time.Date(2050, 12, 10, 16, 25, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 4, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 18:00 ~ 21:00)", mail.Sent[3].Title)

	// 16:25.59.999
	current = time.Date(2050, 12, 10, 16, 25, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 5, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 18:00 ~ 21:00)", mail.Sent[4].Title)
}

func TestNotifyOrderByHour_Sent_NoOrderDueToCancel(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55
	current := time.Date(2050, 12, 10, 9, 30, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 11:30 ~ 15:00)", mail.Sent[0].Title)

	current = time.Date(2050, 12, 10, 9, 31, 0, 0, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 11:30 ~ 15:00)", mail.Sent[1].Title)

	// 9:54
	current = time.Date(2050, 12, 10, 9, 54, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 3, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 11:30 ~ 15:00)", mail.Sent[2].Title)

	// 9:55
	current = time.Date(2050, 12, 10, 9, 55, 00, 000, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 4, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 11:30 ~ 15:00)", mail.Sent[3].Title)

	// 9:55.59.999
	current = time.Date(2050, 12, 10, 9, 55, 59, 999, jst) // sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 5, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/10 11:30 ~ 15:00)", mail.Sent[4].Title)
}

// special holiday
func TestNotifyOrderByHour_NotSent_SpecialHoliday(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()
	// 2050/11/20 ~ 2050/12/08 holidays

	// morning:7:00 - 2:00 => 5:00 ~ 5:25
	// 5:05
	current := time.Date(2050, 12, 3, 5, 5, 0, 0, jst) // 12/3 sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55
	// 9:35
	current = time.Date(2050, 12, 3, 9, 35, 0, 0, jst) // 12/3 sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25
	// 16:05
	current = time.Date(2050, 12, 3, 16, 5, 0, 0, jst) // 12/3 sat
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// 5:05
	current = time.Date(2050, 12, 4, 5, 5, 0, 0, jst) // 12/4 sun
	useCase.NotifyOrderByHour(current)
	assert.Equal(t, 0, len(mail.Sent))
	current = time.Date(2050, 12, 5, 5, 5, 0, 0, jst) // 12/5 mon
	useCase.NotifyOrderByHour(current)
	assert.Equal(t, 0, len(mail.Sent))
	current = time.Date(2050, 12, 6, 5, 5, 0, 0, jst) // 12/6 tue
	useCase.NotifyOrderByHour(current)
	assert.Equal(t, 0, len(mail.Sent))
	current = time.Date(2050, 12, 7, 5, 5, 0, 0, jst) // 12/7 wed
	useCase.NotifyOrderByHour(current)
	assert.Equal(t, 0, len(mail.Sent))
	current = time.Date(2050, 12, 8, 5, 5, 0, 0, jst) // 12/8 thr
	useCase.NotifyOrderByHour(current)
	assert.Equal(t, 0, len(mail.Sent))
}

// special hour (single, multiple)

func TestNotifyOrderByHour_Sent_NoOrder_SingleSpecialHours(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()
	// sp morning "2050/12/17" 08:00 - 12:00

	// morning:8:00 - 2:00 => 6:00 ~ 6:25
	// 5:59
	current := time.Date(2050, 12, 17, 5, 59, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// 6:00
	current = time.Date(2050, 12, 17, 6, 0, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/17 08:00 ~ 12:00)", mail.Sent[0].Title)

	// 6:25.59
	current = time.Date(2050, 12, 17, 6, 25, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/17 08:00 ~ 12:00)", mail.Sent[1].Title)

	// 6:26 (not sent)
	current = time.Date(2050, 12, 17, 6, 26, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55 => should not sent
	// 9:35
	current = time.Date(2050, 12, 17, 9, 35, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25 => should not sent
	// 16:10
	current = time.Date(2050, 12, 17, 16, 10, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
}

func TestNotifyOrderByHour_Sent_NoOrder_MultiSpecialHours(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()
	// sp morning "2050/12/18" 08:00 - 12:00
	// sp dinner "2050/12/18" 16:00 - 19:00

	// morning:8:00 - 2:00 => 6:00 ~ 6:25
	// 5:59
	current := time.Date(2050, 12, 18, 5, 59, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// 6:00
	current = time.Date(2050, 12, 18, 6, 0, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/18 08:00 ~ 12:00)", mail.Sent[0].Title)

	// 6:25.59
	current = time.Date(2050, 12, 18, 6, 25, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/18 08:00 ~ 12:00)", mail.Sent[1].Title)

	// 6:26 (not sent)
	current = time.Date(2050, 12, 18, 6, 26, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55 => should not sent
	// 9:35
	current = time.Date(2050, 12, 18, 9, 35, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// dinner:16:00 - 19:00 => 14:00 ~ 14:25
	// 16:10
	current = time.Date(2050, 12, 18, 14, 10, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 3, len(mail.Sent))
	assert.Equal(t, "オーダーはありません(2050/12/18 16:00 ~ 19:00)", mail.Sent[2].Title)
}

func TestNotifyOrderByHour_Sent_WithOrder_SingleSpecialHours(t *testing.T) {
	setUpEnv(t)
	useCase, mail := setUpUseCase()
	// sp morning "2050/12/19" 08:00 - 12:00

	// morning:8:00 - 2:00 => 6:00 ~ 6:25
	// 5:59
	current := time.Date(2050, 12, 19, 5, 59, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 0, len(mail.Sent))

	// 6:00
	current = time.Date(2050, 12, 19, 6, 0, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 1, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/19 08:00 ~ 12:00)", mail.Sent[0].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[0].Message, "注文数:1"))

	// 6:25.59
	current = time.Date(2050, 12, 19, 6, 25, 59, 999, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
	assert.Equal(t, "本日のオーダー情報(2050/12/19 08:00 ~ 12:00)", mail.Sent[1].Title)
	assert.Equal(t, true, strings.Contains(mail.Sent[1].Message, "注文数:1"))

	// 6:26 (not sent)
	current = time.Date(2050, 12, 19, 6, 26, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// lunch:11:30 - 2:00 => 9:30 ~ 9:55 => should not sent
	// 9:35
	current = time.Date(2050, 12, 19, 9, 35, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))

	// dinner:"18:00" ~ "21:00" => 16:00 ~ 16:25 => should not sent
	// 16:10
	current = time.Date(2050, 12, 19, 16, 10, 0, 0, jst)
	useCase.NotifyOrderByHour(current)

	assert.Equal(t, 2, len(mail.Sent))
}

func setUpUseCase() (order.OrderTaskUseCase, *memory.MemorySendOrderMail) {
	repo := memory.NewOrderInfoMemoryRepository()
	orders := repo.GetMemory()
	// create additional order
	foodOrders1 := []domains.OrderFoodItem{}
	foodOrder1, err := domains.NewOrderFoodItem("abc", "food123", 210, 3, []domains.OptionItemInfo{})
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders1 = append(foodOrders1, *foodOrder1)

	// morning
	stockOrders1 := []domains.OrderStockItem{}
	order1, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 9:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orders[order1.GetId()] = order1

	// create additional canceled order
	order2, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 9:00", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	order2.SetCancel()
	orders[order2.GetId()] = order2

	// lunch and cancel
	order3, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/11 12:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	order3.SetCancel()
	orders[order3.GetId()] = order3

	// morning, _ := NewBusinessHour("morning", "07:00", "09:30", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	// lunch, _ := NewBusinessHour("lunch", "11:30", "15:00", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	// dinner, _ := NewBusinessHour("dinner", "18:00", "21:00", []Weekday{Wednesday, Saturday})

	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	holidayRepo := memory.NewSpecialHolidayMemoryRepository()

	holiday1, err := stDomains.NewSpecialHoliday("おやすみXX", "2050/11/20", "2050/12/08")
	if err != nil {
		fmt.Println(err)
		panic("failed to create holidays")
	}
	holidayRepo.Create(holiday1)

	schedules := businessHourRepo.GetMemory().GetSchedules()
	spHour1, err := stDomains.NewSpecialBusinessHour("特別モーニング", "2050/12/17", "08:00", "12:00", schedules[0].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	spBusinessHourRepo.Create(spHour1)

	spHour2, err := stDomains.NewSpecialBusinessHour("特別モーニング2", "2050/12/18", "08:00", "12:00", schedules[0].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	spBusinessHourRepo.Create(spHour2)

	spHour3, err := stDomains.NewSpecialBusinessHour("特別ディナー2", "2050/12/18", "16:00", "19:00", schedules[2].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	spBusinessHourRepo.Create(spHour3)

	// sp morning with order
	spHour4, err := stDomains.NewSpecialBusinessHour("特別モーニング3", "2050/12/19", "08:00", "12:00", schedules[0].GetId())
	if err != nil {
		fmt.Println(err)
		panic("failed to create special holiday")
	}
	spBusinessHourRepo.Create(spHour4)

	orderSp, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/19 9:10", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orders[orderSp.GetId()] = orderSp

	mail := memory.NewMemorySendOrderMail()
	useCase := order.NewOrderTaskUseCase(repo, mail, businessHourRepo, holidayRepo, spBusinessHourRepo)

	return useCase, mail
}

func setUpUseCaseWithMultipleOrders() (order.OrderTaskUseCase, *memory.MemorySendOrderMail) {
	repo := memory.NewOrderInfoMemoryRepository()
	orders := repo.GetMemory()
	// create additional order
	foodOrders1 := []domains.OrderFoodItem{}
	foodOrder1, err := domains.NewOrderFoodItem("abc", "food123", 210, 3, []domains.OptionItemInfo{})
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders1 = append(foodOrders1, *foodOrder1)

	// morning
	stockOrders1 := []domains.OrderStockItem{}
	order1, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 7:00", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orders[order1.GetId()] = order1

	// morning
	foodOrders2 := []domains.OrderFoodItem{}
	foodOrder21, err := domains.NewOrderFoodItem("abc", "food123", 410, 3, []domains.OptionItemInfo{})
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders2 = append(foodOrders2, *foodOrder21)
	foodOrder22, err := domains.NewOrderFoodItem("abcd", "food1234", 510, 1, []domains.OptionItemInfo{})
	if err != nil {
		fmt.Println(err)
		panic("failed to create food order")
	}
	foodOrders2 = append(foodOrders2, *foodOrder22)

	stockOrders2 := []domains.OrderStockItem{}
	stockOrder2, err := domains.NewOrderStockItem("abc", "food123", 210, 3, []domains.OptionItemInfo{})
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	stockOrders2 = append(stockOrders2, *stockOrder2)
	order2, err := domains.NewOrderInfo("user3", "ユーザー3", "user3@hoge.com", "123456789", "memo1", "2050/12/10 9:00", stockOrders2, foodOrders2)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	orders[order2.GetId()] = order2

	// morning and cancel
	order3, err := domains.NewOrderInfo("user2", "ユーザー2", "user2@hoge.com", "123456789", "memo1", "2050/12/10 8:30", stockOrders1, foodOrders1)
	if err != nil {
		fmt.Println(err)
		panic("failed to create stock order")
	}
	order3.SetCancel()
	orders[order3.GetId()] = order3

	// morning, _ := NewBusinessHour("morning", "07:00", "09:30", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	// lunch, _ := NewBusinessHour("lunch", "11:30", "15:00", []Weekday{Tuesday, Wednesday, Friday, Saturday, Sunday})
	// dinner, _ := NewBusinessHour("dinner", "18:00", "21:00", []Weekday{Wednesday, Saturday})

	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	holidayRepo := memory.NewSpecialHolidayMemoryRepository()

	mail := memory.NewMemorySendOrderMail()
	useCase := order.NewOrderTaskUseCase(repo, mail, businessHourRepo, holidayRepo, spBusinessHourRepo)

	return useCase, mail
}

func setUpEnv(t *testing.T) {
	t.Setenv("MAIL_FROM", "from@dummy.co.jp")
	t.Setenv("MAIL_ADMIN", "admin@dummy.co.jp")

	t.Setenv("APP_PORT", "80")
	t.Setenv("GOOGLE_CREDENTIALS_JSON", "test")

	if err := common.InitConfig(true); err != nil {
		assert.Fail(t, "failed to init config")
	}
}
