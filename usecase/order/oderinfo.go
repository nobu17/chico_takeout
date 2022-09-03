package order

import (
	"context"
	"fmt"
	"os"

	"chico/takeout/common"
	idomains "chico/takeout/domains/item"
	domains "chico/takeout/domains/order"
	sdomains "chico/takeout/domains/store"
	"chico/takeout/usecase"
)

type OrderInfoModel struct {
	Id             string
	UserId         string
	UserName       string
	UserEmail      string
	UserTelNo      string
	Memo           string
	OrderDateTime  string
	PickupDateTime string
	StockItems     []CommonItemOrderModel
	FoodItems      []CommonItemOrderModel
	Canceled       bool
}

type CommonItemOrderModel struct {
	ItemId   string
	Name     string
	Price    int
	Quantity int
}

func newCommonItemOrderModel(itemId, name string, price, quantity int) *CommonItemOrderModel {
	return &CommonItemOrderModel{
		ItemId:   itemId,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
}

func newOrderInfoModel(item *domains.OrderInfo) *OrderInfoModel {
	stocks := []CommonItemOrderModel{}
	for _, stock := range item.GetStockItems() {
		stocks = append(stocks, *newCommonItemOrderModel(stock.GetItemId(), stock.GetName(), stock.GetPrice(), stock.GetQuantity()))
	}
	foods := []CommonItemOrderModel{}
	for _, food := range item.GetFoodItems() {
		foods = append(foods, *newCommonItemOrderModel(food.GetItemId(), food.GetName(), food.GetPrice(), food.GetQuantity()))
	}
	return &OrderInfoModel{
		Id:             item.GetId(),
		UserId:         item.GetUserId(),
		UserName:       item.GetUserName(),
		UserEmail:      item.GetUserEmail(),
		UserTelNo:      item.GetUserTelNo(),
		Memo:           item.GetMemo(),
		OrderDateTime:  item.GetOrderDateTime(),
		PickupDateTime: item.GetPickupDateTime(),
		Canceled:       item.GetCanceled(),
		StockItems:     stocks,
		FoodItems:      foods,
	}
}

// type CommonItemOrderInfoModel struct {
// 	CommonItemOrderModel
// }

type OrderInfoCreateModel struct {
	UserId         string
	UserName       string
	UserEmail      string
	UserTelNo      string
	Memo           string
	PickupDateTime string
	StockItems     []CommonItemOrderCreateModel
	FoodItems      []CommonItemOrderCreateModel
}

type OrderUserInfoUpdateModel struct {
	OrderId   string
	UserId    string
	UserName  string
	UserEmail string
	UserTelNo string
	Memo      string
}

func NewOrderUserInfoUpdateModel(orderId, userId, userName, userEmail, userTelNo, memo string) *OrderUserInfoUpdateModel {
	return &OrderUserInfoUpdateModel{
		OrderId:   orderId,
		UserId:    userId,
		UserName:  userName,
		UserEmail: userEmail,
		UserTelNo: userTelNo,
		Memo:      memo,
	}
}

type CommonItemOrderCreateModel struct {
	ItemId   string
	Quantity int
}

func newCommonItemOrderCreateModel(itemId string, quantity int) *CommonItemOrderCreateModel {
	return &CommonItemOrderCreateModel{
		ItemId:   itemId,
		Quantity: quantity,
	}
}

type OrderInfoUseCase interface {
	InitContext(ctx context.Context)
	Find(id string) (*OrderInfoModel, error)
	FindAll() ([]OrderInfoModel, error)
	FindByUserId(userId string) ([]OrderInfoModel, error)
	FindActiveByUserId(userId string) ([]OrderInfoModel, error)
	FindActiveByPickupDate(dateStr string) ([]OrderInfoModel, error)
	Create(model *OrderInfoCreateModel) (string, error)
	UpdateUserInfo(model *OrderUserInfoUpdateModel) error
	Cancel(id string) error
}

type orderInfoUseCase struct {
	*usecase.BaseUseCase
	orderInfoRepository   domains.OrderInfoRepository
	stockRepo             idomains.StockItemRepository
	busRepo               sdomains.BusinessHoursRepository
	spBusRepo             sdomains.SpecialBusinessHourRepository
	spHolidayRepo         sdomains.SpecialHolidayRepository
	factory               domains.OrderInfoFactory
	stockConsumer         domains.StockItemRemainCheckAndConsumer
	foodRemainChecker     domains.FoodItemRemainChecker
	orderDuplicateChecker domains.OrderDuplicateChecker
	mailerService         SendOrderMailService
}

func NewOrderInfoUseCase(
	orderInfoRepository domains.OrderInfoRepository,
	stockRepo idomains.StockItemRepository,
	foodRepo idomains.FoodItemRepository,
	busRepo sdomains.BusinessHoursRepository,
	spBusRepo sdomains.SpecialBusinessHourRepository,
	spHolidayRepo sdomains.SpecialHolidayRepository,
	mailerService SendOrderMailService,
) OrderInfoUseCase {
	return &orderInfoUseCase{
		BaseUseCase:           usecase.NewBaseUseCase(),
		orderInfoRepository:   orderInfoRepository,
		stockRepo:             stockRepo,
		busRepo:               busRepo,
		spBusRepo:             spBusRepo,
		spHolidayRepo:         spHolidayRepo,
		factory:               *domains.NewOrderInfoFactory(stockRepo, foodRepo),
		stockConsumer:         *domains.NewStockItemRemainCheckAndConsumer(stockRepo),
		foodRemainChecker:     *domains.NewFoodItemRemainChecker(orderInfoRepository, foodRepo),
		orderDuplicateChecker: *domains.NewOrderDuplicateChecker(orderInfoRepository),
		mailerService:         mailerService,
	}
}

func (o *orderInfoUseCase) Find(id string) (*OrderInfoModel, error) {
	item, err := o.orderInfoRepository.Find(id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, common.NewNotFoundError(id)
	}

	return newOrderInfoModel(item), nil
}

func (o *orderInfoUseCase) FindAll() ([]OrderInfoModel, error) {
	items, err := o.orderInfoRepository.FindAll()
	if err != nil {
		return nil, err
	}

	orders := []OrderInfoModel{}
	for _, item := range items {
		order := newOrderInfoModel(&item)
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o *orderInfoUseCase) FindByUserId(userId string) ([]OrderInfoModel, error) {
	userOrders, err := o.orderInfoRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	orders := []OrderInfoModel{}
	for _, userOrder := range userOrders {
		order := newOrderInfoModel(&userOrder)
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o *orderInfoUseCase) FindActiveByUserId(userId string) ([]OrderInfoModel, error) {
	userOrders, err := o.orderInfoRepository.FindActiveByUserId(userId)
	if err != nil {
		return nil, err
	}

	orders := []OrderInfoModel{}
	for _, userOrder := range userOrders {
		order := newOrderInfoModel(&userOrder)
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o *orderInfoUseCase) FindActiveByPickupDate(dateStr string) ([]OrderInfoModel, error) {
	// empty treats now
	targetDate := common.GetNowDate()
	if dateStr != "" {
		converted, err := common.ConvertHyphenStrToDate(dateStr)
		if err != nil {
			return nil, common.NewValidationError("date", fmt.Sprintf("not allowed date format:%s", dateStr))
		}
		targetDate = converted
	}

	dateOrders, err := o.orderInfoRepository.FindByPickupDate(common.ConvertTimeToDateStr(*targetDate))
	if err != nil {
		return nil, err
	}

	orders := []OrderInfoModel{}
	for _, userOrder := range dateOrders {
		order := newOrderInfoModel(&userOrder)
		orders = append(orders, *order)
	}

	return orders, nil
}

func (o *orderInfoUseCase) Create(model *OrderInfoCreateModel) (string, error) {
	// todo: currently food item schedule id and pickup date time relation is not checking

	// if not admin, can not reserve 2 times.
	if !o.IsAdmin() {
		fmt.Println("not admin. checking order duplicated...")
		duplicated, err := o.orderDuplicateChecker.ActiveOrderExists(model.UserId)
		if err != nil {
			return "", err
		}
		if duplicated {
			return "", common.NewValidationError("UserId", "active order is already exists")
		}
	}

	stockOrders := []domains.ItemOrder{}
	for _, item := range model.StockItems {
		stockOrders = append(stockOrders, *domains.NewItemOrder(item.ItemId, item.Quantity))
	}
	foodOrders := []domains.ItemOrder{}
	for _, item := range model.FoodItems {
		foodOrders = append(foodOrders, *domains.NewItemOrder(item.ItemId, item.Quantity))
	}
	// factory check each item id existence also (will return error)
	// factory check pickup date time is past or not
	order, err := o.factory.Create(model.UserId, model.UserName, model.UserEmail, model.UserTelNo, model.Memo, model.PickupDateTime, stockOrders, foodOrders)
	if err != nil {
		return "", err
	}

	var gError error = nil
	var id = ""
	o.orderInfoRepository.Transact(func() error {
		schedules, err := o.busRepo.Fetch()
		if err != nil {
			gError = err
			return err
		}
		spSchedules, err := o.spBusRepo.FindAll()
		if err != nil {
			gError = err
			return err
		}
		spHoliday, err := o.spHolidayRepo.FindAll()
		if err != nil {
			gError = err
			return err
		}
		// check pickup time is in store business time
		holidaySpec := sdomains.NewHolidaySpecification(*schedules, spSchedules, spHoliday)
		isInBusiness, err := holidaySpec.IsStoreInBusiness(model.PickupDateTime)
		if err != nil {
			gError = err
			return err
		}

		if !isInBusiness {
			gError = common.NewValidationError("PickupDateTime", "pickup time is not in store business")
			return gError
		}

		// check and update stock remain
		err = o.stockConsumer.ConsumeRemainStock(stockOrders)
		if err != nil {
			gError = err
			return err
		}
		// check food remain
		err = o.foodRemainChecker.CheckRemain(order.GetPickupDate(), order.GetFoodItems())
		if err != nil {
			gError = err
			return err
		}
		// create order
		id, err = o.orderInfoRepository.Create(order)
		if err != nil {
			gError = err
			return err
		}
		return nil
	})

	if gError != nil {
		return "", gError
	}

	mError := o.sendCompleteMail(order)
	// mail error not treats as error only displaying as info
	if mError != nil {
		fmt.Printf("mail send error.%s", mError)
	}

	return id, nil
}

func (o *orderInfoUseCase) Cancel(id string) error {
	order, err := o.orderInfoRepository.Find(id)
	if err != nil {
		return err
	}
	if order == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}
	order.SetCancel()
	// increment stock
	err = o.stockConsumer.IncrementCanceledRemain(order.GetStockItems())
	if err != nil {
		return err
	}
	upErr := o.orderInfoRepository.UpdateOrderStatus(order)
	if upErr != nil {
		return upErr
	}

	mError := o.sendCancelMail(order)
	// mail error not treats as error only displaying as info
	if mError != nil {
		fmt.Printf("mail send error.%s", mError)
	}
	return nil
}

func (o *orderInfoUseCase) UpdateUserInfo(model *OrderUserInfoUpdateModel) error {
	order, err := o.orderInfoRepository.Find(model.OrderId)
	if err != nil {
		return err
	}
	if order == nil {
		return common.NewNotFoundError(fmt.Sprintf("order is not exist. id:%s", model.OrderId))
	}
	// check user id is same
	userId := o.GetUserId()
	if model.UserId != userId {
		return common.NewValidationError("UserID", "UserId is invalid. not match authorized user.")
	}

	err = order.UpdateUserInfo(model.UserName, model.UserEmail, model.UserTelNo, model.Memo)
	if err != nil {
		return err
	}

	err = o.orderInfoRepository.UpdateUserInfo(order)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderInfoUseCase) sendCompleteMail(order *domains.OrderInfo) error {
	mailData, err := NewOrderCompleteMailData(order, os.Getenv("MAIL_FROM"))
	if err != nil {
		return err
	}
	return o.mailerService.SendComplete(*mailData)
}

func (o *orderInfoUseCase) sendCancelMail(order *domains.OrderInfo) error {
	mailData, err := NewOrderCancelMailData(order, os.Getenv("MAIL_FROM"))
	if err != nil {
		return err
	}
	return o.mailerService.SendCancel(*mailData)
}
