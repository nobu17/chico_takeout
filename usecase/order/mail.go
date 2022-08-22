package order

import (
	"fmt"
	"os"
	"strings"
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/order"
	"chico/takeout/domains/shared/validator"
)

type SendOrderMailService interface {
	SendComplete(data OrderCompleteMailData) error
	SendCancel(data OrderCancelMailData) error
	SendDailySummary(data ReservationSummaryMailData) error
}

type OrderCompleteMailData struct {
	commonMailData
}

func NewOrderCompleteMailData(order *domains.OrderInfo, sendFrom string) (*OrderCompleteMailData, error) {
	title := "予約完了のお知らせ.(CHICO SPICE)"

	b := &strings.Builder{}
	b.WriteString("テイクアウトの予約が完了いたしました。")
	b.WriteString("\n\n")
	b.WriteString("**ご注文内容に関してはマイページからご確認下さい。**")
	b.WriteString("\n")
	b.WriteString("**決済は当日、店舗にて実施させていただきます。**")
	b.WriteString("\n\n")

	b.WriteString("--予約情報--")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("受取日時:%s", order.GetPickupDateTime()))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("氏名:%s", order.GetUserName()))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("E-mail:%s", order.GetUserEmail()))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("TEL:%s", order.GetUserTelNo()))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("要望やメッセージ:%s", order.GetMemo()))
	b.WriteString("\n")

	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("本メールに心当たりが無い方は、お手数ですが(%s)宛にご連絡をお願いいたします。(本メールは送信専用アドレスから送信しているため、直接の返信は不可能です。)", sendFrom))
	b.WriteString("\n")

	message := b.String()
	sendTo := []string{order.GetUserEmail()}
	bcc := os.Getenv("MAIL_ADMIN")

	comm, err := newCommonMailData(title, message, sendFrom, bcc, sendTo)
	if err != nil {
		return nil, err
	}
	return &OrderCompleteMailData{
		commonMailData: *comm,
	}, nil
}

type OrderCancelMailData struct {
	commonMailData
}

func NewOrderCancelMailData(order *domains.OrderInfo, sendFrom string) (*OrderCancelMailData, error) {
	title := "キャンセル完了のお知らせ.(CHICO SPICE)"

	b := &strings.Builder{}
	b.WriteString("テイクアウトの予約をキャンセルいたしました。")
	b.WriteString("\n")
	b.WriteString("またのご利用をお待ちしております。")
	b.WriteString("\n\n")

	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("本メールに心当たりが無い方は、お手数ですが(%s)宛にご連絡をお願いいたします。(本メールは送信専用アドレスから送信しているため、直接の返信は不可能です。)", sendFrom))
	b.WriteString("\n")

	message := b.String()
	sendTo := []string{order.GetUserEmail()}
	bcc := os.Getenv("MAIL_ADMIN")

	comm, err := newCommonMailData(title, message, sendFrom, bcc, sendTo)
	if err != nil {
		return nil, err
	}
	return &OrderCancelMailData{
		commonMailData: *comm,
	}, nil
}

type ReservationSummaryMailData struct {
	commonMailData
}

func NewReservationSummaryMailData(orders []domains.OrderInfo, sendFrom string, startDateTime time.Time) (*ReservationSummaryMailData, error) {
	title := fmt.Sprintf("本日のオーダー情報(%s)", common.ConvertTimeToDateStr(startDateTime))

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("本日のオーダー情報は下記になります。(%s 以降)", common.ConvertTimeToTimeStr(startDateTime)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("注文数:%d", len(orders)))
	b.WriteString("\n")

	for _, order := range orders {
		b.WriteString("-------------")
		b.WriteString("\n")
		b.WriteString("**注文者 情報**")
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("受取日時:%s", order.GetPickupDateTime()))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("氏名:%s", order.GetUserName()))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("E-mail:%s", order.GetUserEmail()))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("TEL:%s", order.GetUserTelNo()))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("要望やメッセージ:%s", order.GetMemo()))
		b.WriteString("\n")
		b.WriteString("\n")
		b.WriteString("**注文 情報**")
		b.WriteString("\n")
		for _, food := range order.GetFoodItems() {
			b.WriteString(fmt.Sprintf("%s, %d円, %d個", food.GetName(), food.GetPrice(), food.GetQuantity()))
			b.WriteString("\n")
		}
		for _, stock := range order.GetStockItems() {
			b.WriteString(fmt.Sprintf("%s, %d円, %d個", stock.GetName(), stock.GetPrice(), stock.GetQuantity()))
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("合計:: %d円", order.GetTotalCost()))
		b.WriteString("\n\n")
	}

	message := b.String()
	sendTo := []string{os.Getenv("MAIL_ADMIN")}
	bcc := ""

	comm, err := newCommonMailData(title, message, sendFrom, bcc, sendTo)
	if err != nil {
		return nil, err
	}
	return &ReservationSummaryMailData{
		commonMailData: *comm,
	}, nil
}

type commonMailData struct {
	Title    string
	SendTo   []string
	SendFrom string
	Bcc      string
	Message  string
}

func newCommonMailData(title, message, sendFrom, bcc string, sendTo []string) (*commonMailData, error) {
	if strings.TrimSpace(title) == "" {
		return nil, common.NewValidationError("title", "empty is not allowed.")
	}
	if strings.TrimSpace(message) == "" {
		return nil, common.NewValidationError("message", "empty is not allowed.")
	}
	check := validator.NewEmailValidator("SendFrom")
	if err := check.Validate(sendFrom); err != nil {
		return nil, err
	}

	// bcc is optional
	if bcc != "" {
		check = validator.NewEmailValidator("Bcc")
		if err := check.Validate(bcc); err != nil {
			return nil, err
		}
	}

	if len(sendTo) == 0 {
		return nil, common.NewValidationError("sendTo", "empty slice is not allowed.")
	}
	for _, mail := range sendTo {
		check = validator.NewEmailValidator("SendTo")
		if err := check.Validate(mail); err != nil {
			return nil, err
		}
	}

	return &commonMailData{
		Title:    title,
		SendTo:   sendTo,
		Bcc:      bcc,
		SendFrom: sendFrom,
		Message:  message,
	}, nil
}