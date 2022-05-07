package order

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/order"

	"github.com/gin-gonic/gin"
)

type OrderInfoCreateRequest struct {
	UserId         string                   `json:"userId" binding:"required"`
	Memo           string                   `json:"memo" binding:"required"`
	PickupDateTime string                   `json:"pickupDateTime" binding:"required"`
	StockItems     []CommonItemOrderRequest `json:"stockItems" binding:"required"`
	FoodItems      []CommonItemOrderRequest `json:"foodItems" binding:"required"`
}

func (o *OrderInfoCreateRequest) toModel() *usecases.OrderInfoCreateModel {
	stocks := []usecases.CommonItemOrderModel{}
	for _, stock := range o.StockItems {
		stocks = append(stocks, *stock.toModel())
	}
	foods := []usecases.CommonItemOrderModel{}
	for _, food := range o.FoodItems {
		foods = append(stocks, *food.toModel())
	}
	return &usecases.OrderInfoCreateModel{
		UserId:         o.UserId,
		Memo:           o.Memo,
		PickupDateTime: o.PickupDateTime,
		StockItems:     stocks,
		FoodItems:      foods,
	}
}

type OrderInfoCreateResponce struct {
	Id string `json:"id" binding:"required"`
}

type CommonItemOrderRequest struct {
	ItemId   string `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type OrderInfoCancelRequest struct {
	Id string `json:"id" binding:"required"`
}

func (c *CommonItemOrderRequest) toModel() *usecases.CommonItemOrderModel {
	return &usecases.CommonItemOrderModel{
		ItemId:   c.ItemId,
		Quantity: c.Quantity,
	}
}

type orderInfoHandler struct {
	*handlers.BaseHandler
	usecase usecases.OrderInfoUseCase
}

func NewSpecialHolidayHandler(usecase usecases.OrderInfoUseCase) *orderInfoHandler {
	return &orderInfoHandler{
		usecase: usecase,
	}
}

func (s *orderInfoHandler) PostCreate(c *gin.Context) {
	var req OrderInfoCreateRequest
	// validation is executed model
	c.ShouldBind(&req)
	id, err := s.usecase.Create(*req.toModel())
	if err != nil {
		s.HandleError(c, err)
	}
	s.HandleOK(c, OrderInfoCreateResponce{Id: id})
}

func (s *orderInfoHandler) PostCancel(c *gin.Context) {
	var req OrderInfoCancelRequest
	// validation is executed model
	c.ShouldBind(&req)
	err := s.usecase.Cancel(req.Id)
	if err != nil {
		s.HandleError(c, err)
	}
	s.HandleOK(c, nil)
}
