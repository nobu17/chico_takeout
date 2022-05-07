package order

import (
	"chico/takeout/common"
	idomains "chico/takeout/domains/item"
	domains "chico/takeout/domains/order"
)

// type OrderInfoModel struct {
// 	Id             string
// 	UserId         string
// 	Memo           string
// 	OrderDateTime  string
// 	PickupDateTime string
// 	StockItems     []CommonItemOrderInfoModel
// 	FoodItems      []CommonItemOrderInfoModel
// 	Canceled       bool
// }

// type CommonItemOrderInfoModel struct {
// 	CommonItemOrderModel
// }

type OrderInfoCreateModel struct {
	UserId         string
	Memo           string
	PickupDateTime string
	StockItems     []CommonItemOrderModel
	FoodItems      []CommonItemOrderModel
}

type CommonItemOrderModel struct {
	ItemId   string
	Quantity int
}

type OrderInfoUseCase struct {
	orderInfoRepository domains.OrderInfoRepository
	stockRepo           idomains.StockItemRepository
	factory             domains.OrderInfoFactory
	stockConsumer       domains.StockItemRemainCheckAndConsumer
	foodRemainChecker   domains.FoodItemRemainChecker
}

func NewOrderInfoUseCase(
	orderInfoRepository domains.OrderInfoRepository,
	stockRepo idomains.StockItemRepository,
	foodRepo idomains.FoodItemRepository) *OrderInfoUseCase {
	return &OrderInfoUseCase{
		orderInfoRepository: orderInfoRepository,
		stockRepo:           stockRepo,
		factory:             *domains.NewOrderInfoFactory(stockRepo, foodRepo),
		stockConsumer:       *domains.NewStockItemRemainCheckAndConsumer(stockRepo),
		foodRemainChecker:   *domains.NewFoodItemRemainChecker(orderInfoRepository, foodRepo),
	}
}

func (o *OrderInfoUseCase) Create(model OrderInfoCreateModel) (string, error) {
	stockOrders := []domains.ItemOrder{}
	for _, item := range model.StockItems {
		stockOrders = append(stockOrders, *domains.NewItemOrder(item.ItemId, item.Quantity))
	}
	foodOrders := []domains.ItemOrder{}
	for _, item := range model.FoodItems {
		foodOrders = append(foodOrders, *domains.NewItemOrder(item.ItemId, item.Quantity))
	}
	order, err := o.factory.Create(model.UserId, model.Memo, model.PickupDateTime, stockOrders, foodOrders)
	if err != nil {
		return "", err
	}
	// check and update stock remain
	err = o.stockConsumer.ConsumeRemainStock(stockOrders)
	if err != nil {
		return "", err
	}
	// check food remain
	err = o.foodRemainChecker.CheckRemain(order.GetPickupDate(), order.GetFoodItems())
	if err != nil {
		return "", err
	}
	// create order
	return o.orderInfoRepository.Create(*order)
}

func (o *OrderInfoUseCase) Cancel(id string) error {
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
	return o.orderInfoRepository.UpdateOrderStatus(*order)
}