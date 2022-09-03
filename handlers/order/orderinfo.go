package order

import (
	"chico/takeout/handlers"
	usecases "chico/takeout/usecase/order"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type OrderInfoData struct {
	Id             string                `json:"id" binding:"required"`
	UserId         string                `json:"userId" binding:"required"`
	UserName       string                `json:"userName" binding:"required"`
	UserEmail      string                `json:"userEmail" binding:"required"`
	UserTelNo      string                `json:"userTelNo" binding:"required"`
	Memo           string                `json:"memo" binding:"required"`
	PickupDateTime string                `json:"pickupDateTime" binding:"required"`
	OrderDateTime  string                `json:"orderDateTime" binding:"required"`
	StockItems     []CommonItemOrderData `json:"stockItems" binding:"required"`
	FoodItems      []CommonItemOrderData `json:"foodItems" binding:"required"`
	Canceled       bool                  `json:"canceled" binding:"required"`
}

type CommonItemOrderData struct {
	ItemId   string `json:"itemId" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Price    int    `json:"price" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

func newCommonItemOrderData(itemId, name string, price, quantity int) *CommonItemOrderData {
	return &CommonItemOrderData{
		ItemId:   itemId,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
}

func newOrderInfoData(item *usecases.OrderInfoModel) *OrderInfoData {
	stocks := []CommonItemOrderData{}
	for _, stock := range item.StockItems {
		stocks = append(stocks, *newCommonItemOrderData(stock.ItemId, stock.Name, stock.Price, stock.Quantity))
	}
	foods := []CommonItemOrderData{}
	for _, stock := range item.FoodItems {
		foods = append(foods, *newCommonItemOrderData(stock.ItemId, stock.Name, stock.Price, stock.Quantity))
	}
	return &OrderInfoData{
		Id:             item.Id,
		UserId:         item.UserId,
		UserName:       item.UserName,
		UserEmail:      item.UserEmail,
		UserTelNo:      item.UserTelNo,
		Memo:           item.Memo,
		OrderDateTime:  item.OrderDateTime,
		PickupDateTime: item.PickupDateTime,
		Canceled:       item.Canceled,
		StockItems:     stocks,
		FoodItems:      foods,
	}
}

type OrderInfoCreateRequest struct {
	UserId         string                   `json:"userId" binding:"required"`
	UserName       string                   `json:"userName" binding:"required"`
	UserEmail      string                   `json:"userEmail" binding:"required"`
	UserTelNo      string                   `json:"userTelNo" binding:"required"`
	Memo           string                   `json:"memo"`
	PickupDateTime string                   `json:"pickupDateTime" binding:"required"`
	StockItems     []CommonItemOrderRequest `json:"stockItems" binding:"required"`
	FoodItems      []CommonItemOrderRequest `json:"foodItems" binding:"required"`
}

func (o *OrderInfoCreateRequest) toModel() *usecases.OrderInfoCreateModel {
	stocks := []usecases.CommonItemOrderCreateModel{}
	for _, stock := range o.StockItems {
		stocks = append(stocks, *stock.toModel())
	}
	foods := []usecases.CommonItemOrderCreateModel{}
	for _, food := range o.FoodItems {
		foods = append(foods, *food.toModel())
	}
	return &usecases.OrderInfoCreateModel{
		UserId:         o.UserId,
		UserName:       o.UserName,
		UserEmail:      o.UserEmail,
		UserTelNo:      o.UserTelNo,
		Memo:           o.Memo,
		PickupDateTime: o.PickupDateTime,
		StockItems:     stocks,
		FoodItems:      foods,
	}
}

type OrderInfoCreateResponse struct {
	Id string `json:"id" binding:"required"`
}

type CommonItemOrderRequest struct {
	ItemId   string `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type OrderInfoCancelRequest struct {
	Id string
}

func (c *CommonItemOrderRequest) toModel() *usecases.CommonItemOrderCreateModel {
	return &usecases.CommonItemOrderCreateModel{
		ItemId:   c.ItemId,
		Quantity: c.Quantity,
	}
}

type OrderUserInfoUpdateRequest struct {
	OrderId   string
	UserId    string
	UserName  string `json:"userName" binding:"required"`
	UserEmail string `json:"userEmail" binding:"required"`
	UserTelNo string `json:"userTelNo" binding:"required"`
	Memo      string `json:"memo"`
}

func (o *OrderUserInfoUpdateRequest) toModel() *usecases.OrderUserInfoUpdateModel {
	return &usecases.OrderUserInfoUpdateModel{
		OrderId:   o.OrderId,
		UserId:    o.UserId,
		UserName:  o.UserName,
		UserEmail: o.UserEmail,
		UserTelNo: o.UserTelNo,
		Memo:      o.Memo,
	}
}

type orderInfoHandler struct {
	*handlers.BaseHandler
	usecase usecases.OrderInfoUseCase
}

func NewOrderInfoHandler(usecase usecases.OrderInfoUseCase) *orderInfoHandler {
	return &orderInfoHandler{
		usecase: usecase,
	}
}

func (s *orderInfoHandler) InitContext(ctx context.Context) {
	s.usecase.InitContext(ctx)
}

func (s *orderInfoHandler) GetAll(c *gin.Context) {
	models, err := s.usecase.FindAll()
	if err != nil {
		s.HandleError(c, err)
		return
	}

	orders := []OrderInfoData{}
	for _, model := range models {
		order := newOrderInfoData(&model)
		orders = append(orders, *order)
	}
	s.HandleOK(c, orders)
}

func (s *orderInfoHandler) Get(c *gin.Context) {
	id := c.Param("id")
	model, err := s.usecase.Find(id)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, newOrderInfoData(model))
}

func (s *orderInfoHandler) GetByUser(c *gin.Context) {
	id := c.Param("userId")
	models, err := s.usecase.FindByUserId(id)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	orders := []OrderInfoData{}
	for _, model := range models {
		order := newOrderInfoData(&model)
		orders = append(orders, *order)
	}
	s.HandleOK(c, orders)
}

func (s *orderInfoHandler) GetActiveByUser(c *gin.Context) {
	id := c.Param("userId")
	models, err := s.usecase.FindActiveByUserId(id)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	orders := []OrderInfoData{}
	for _, model := range models {
		order := newOrderInfoData(&model)
		orders = append(orders, *order)
	}
	s.HandleOK(c, orders)
}

func (s *orderInfoHandler) GetActiveByDate(c *gin.Context) {
	date := c.Param("date")
	// remove slash (optional url is added)
	date = strings.Replace(date, "/", "", -1)
	models, err := s.usecase.FindActiveByPickupDate(date)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	orders := []OrderInfoData{}
	for _, model := range models {
		order := newOrderInfoData(&model)
		orders = append(orders, *order)
	}
	s.HandleOK(c, orders)
}

func (s *orderInfoHandler) PostCreate(c *gin.Context) {
	var req OrderInfoCreateRequest
	if !s.ShouldBind(c, &req) {
		return
	}
	id, err := s.usecase.Create(req.toModel())
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, OrderInfoCreateResponse{Id: id})
}

func (s *orderInfoHandler) PutCancel(c *gin.Context) {
	id := c.Param("id")
	req := OrderInfoCancelRequest{Id: id}
	err := s.usecase.Cancel(req.Id)
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, nil)
}

func (s *orderInfoHandler) PutUpdateUserInfo(c *gin.Context) {
	userId := c.Param("userId")
	orderId := c.Param("orderId")
	var req OrderUserInfoUpdateRequest
	if !s.ShouldBind(c, &req) {
		return
	}
	req.OrderId = orderId
	req.UserId = userId
	err := s.usecase.UpdateUserInfo(req.toModel())
	if err != nil {
		s.HandleError(c, err)
		return
	}
	s.HandleOK(c, nil)
}
