package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/order"
	"chico/takeout/infrastructures/rdbms"
	"chico/takeout/infrastructures/rdbms/items"

	"database/sql/driver"

	"gorm.io/gorm"
)

type OrderInfoRepository struct {
	rdbms.BaseRepository
}

func NewOrderInfoRepository(db *gorm.DB) (*OrderInfoRepository, error) {
	return &OrderInfoRepository{
		BaseRepository: rdbms.BaseRepository{Db: db},
	}, nil
}

type OrderInfoModel struct {
	rdbms.BaseModel
	UserID                 string
	UserName               string
	UserEmail              string
	UserTelNo              string
	Memo                   string
	OrderDateTime          time.Time
	PickupDateTime         time.Time
	Canceled               bool
	StockItemModels        []items.StockItemModel `gorm:"many2many:orderInfo_stockItems;"`
	FoodItemModels         []items.FoodItemModel  `gorm:"many2many:orderInfo_foodItems;"`
	OrderedStockItemModels []OrderedStockItemModel
	OrderedFoodItemModels  []OrderedFoodItemModel
}

type OrderedStockItemModel struct {
	OrderInfoModelID string `gorm:"primaryKey"`
	StockItemModelID string `gorm:"primaryKey"`
	Name             string
	Price            int
	Quantity         int
	Options          []OrderedStockOptionItemModel `gorm:"serializer:json"`
}

type OrderedStockOptionItemModel struct {
	OptionItemModelID       string
	Name                    string
	Price                   int
	Quantity                int
}

func (r *OrderedStockOptionItemModel) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}

	return json.Unmarshal([]byte(val), r)
}

func (r OrderedStockOptionItemModel) Value() (driver.Value, error) {
	val, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	return val, nil
}

type OrderedFoodItemModel struct {
	OrderInfoModelID string `gorm:"primaryKey"`
	FoodItemModelID  string `gorm:"primaryKey"`
	Name             string
	Price            int
	Quantity         int
	Options          []OrderedFoodOptionItemModel `gorm:"serializer:json"`
}

type OrderedFoodOptionItemModel struct {
	OptionItemModelID      string
	Name                   string
	Price                  int
	Quantity               int
}

func (r *OrderedFoodOptionItemModel) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}

	return json.Unmarshal([]byte(val), r)
}

func (r OrderedFoodOptionItemModel) Value() (driver.Value, error) {
	val, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func newOrderInfoModel(order *domains.OrderInfo) (*OrderInfoModel, error) {

	orderDateTime, err := common.ConvertStrToDateTime(order.GetOrderDateTime())
	if err != nil {
		return nil, err
	}
	pickupDateTime, err := common.ConvertStrToDateTime(order.GetPickupDateTime())
	if err != nil {
		return nil, err
	}

	model := &OrderInfoModel{}
	model.ID = order.GetId()
	model.UserID = order.GetUserId()
	model.UserName = order.GetUserName()
	model.UserEmail = order.GetUserEmail()
	model.UserTelNo = order.GetUserTelNo()
	model.Memo = order.GetMemo()
	model.OrderDateTime = *orderDateTime
	model.PickupDateTime = *pickupDateTime

	// below data is not needed to insert

	// stockModels := []items.StockItemModel{}
	// for _, stock := range order.GetStockItems() {
	// 	stockModel := items.StockItemModel{}
	// 	stockModel.ID = stock.GetItemId()
	// 	stockModels = append(stockModels, stockModel)
	// }
	// model.StockItemModels = stockModels

	// foodModels := []items.FoodItemModel{}
	// for _, food := range order.GetFoodItems() {
	// 	foodModel := items.FoodItemModel{}
	// 	foodModel.ID = food.GetItemId()
	// 	foodModels = append(foodModels, foodModel)
	// }
	// model.FoodItemModels = foodModels

	return model, nil
}

func (s *OrderInfoModel) toDomain(stocks []OrderedStockItemModel, foods []OrderedFoodItemModel) (*domains.OrderInfo, error) {
	stockDoms := []domains.OrderStockItem{}
	for _, stock := range stocks {
		opts, err := s.toStockOptDomain(stock.Options)
		if err != nil {
			return nil, err
		}
		stockDom, err := domains.NewOrderStockItem(stock.StockItemModelID, stock.Name, stock.Price, stock.Quantity, opts)
		if err != nil {
			return nil, err
		}
		stockDoms = append(stockDoms, *stockDom)
	}
	foodDoms := []domains.OrderFoodItem{}
	for _, food := range foods {
		opts, err := s.toFoodOptDomain(food.Options)
		if err != nil {
			return nil, err
		}
		foodDom, err := domains.NewOrderFoodItem(food.FoodItemModelID, food.Name, food.Price, food.Quantity, opts)
		if err != nil {
			return nil, err
		}
		foodDoms = append(foodDoms, *foodDom)
	}

	pickUp := common.ConvertTimeToDateTimeStr(s.PickupDateTime)
	ordered := common.ConvertTimeToDateTimeStr(s.OrderDateTime)
	dom, err := domains.NewOrderInfoForOrm(s.ID, s.UserID, s.UserName, s.UserEmail, s.UserTelNo, s.Memo, pickUp, ordered, stockDoms, foodDoms, s.Canceled)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (o *OrderInfoModel) toFoodOptDomain(opts []OrderedFoodOptionItemModel) ([]domains.OptionItemInfo, error) {
	options := []domains.OptionItemInfo{}
	for _, opt := range opts {
		op, err := domains.NewOptionItemInfo(opt.OptionItemModelID, opt.Name, opt.Price)
		if err != nil {
			return nil, err
		}
		options = append(options, *op)
	}
	return options, nil
}
func (o *OrderInfoModel) toStockOptDomain(opts []OrderedStockOptionItemModel) ([]domains.OptionItemInfo, error) {
	options := []domains.OptionItemInfo{}
	for _, opt := range opts {
		op, err := domains.NewOptionItemInfo(opt.OptionItemModelID, opt.Name, opt.Price)
		if err != nil {
			return nil, err
		}
		options = append(options, *op)
	}
	return options, nil
}

func (o *OrderInfoRepository) Find(id string) (*domains.OrderInfo, error) {
	model := OrderInfoModel{}

	err := o.Db.First(&model, "ID=?", id).Error
	if err != nil {
		return nil, err
	}
	stocks := []OrderedStockItemModel{}
	err = o.Db.Where("order_info_model_id = ?", id).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	foods := []OrderedFoodItemModel{}
	err = o.Db.Where("order_info_model_id = ?", id).Find(&foods).Error
	if err != nil {
		return nil, err
	}

	dom, err := model.toDomain(stocks, foods)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (o *OrderInfoRepository) FindAll() ([]domains.OrderInfo, error) {
	maxLimit := 1000
	models := []OrderInfoModel{}

	// until 1000 order by ordered time(latest ordered item)
	err := o.Db.Preload("OrderedStockItemModels").
		Preload("OrderedFoodItemModels").
		Preload("StockItemModels").
		Preload("FoodItemModels").Limit(maxLimit).Order("order_date_time desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		order, err := model.toDomain(model.OrderedStockItemModels, model.OrderedFoodItemModels)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (o *OrderInfoRepository) FindByPickupDate(date string) ([]domains.OrderInfo, error) {
	pickupDateStart, err := common.ConvertStrToDate(date)
	if err != nil {
		return nil, err
	}
	pickupDateEnd := pickupDateStart.AddDate(0, 0, 1)

	models := []OrderInfoModel{}
	err = o.Db.Preload("OrderedStockItemModels").
		Preload("OrderedFoodItemModels").Where("pickup_date_time >= ? and pickup_date_time<= ?", pickupDateStart, pickupDateEnd).Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		order, err := model.toDomain(model.OrderedStockItemModels, model.OrderedFoodItemModels)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (o *OrderInfoRepository) FindByUserId(userId string) ([]domains.OrderInfo, error) {
	models := []OrderInfoModel{}
	err := o.Db.Preload("OrderedStockItemModels").Preload("OrderedFoodItemModels").Where("user_id = ?", userId).Order("pickup_date_time desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		order, err := model.toDomain(model.OrderedStockItemModels, model.OrderedFoodItemModels)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (o *OrderInfoRepository) FindActiveByUserId(userId string) ([]domains.OrderInfo, error) {
	models := []OrderInfoModel{}
	// until 30 minutes passed, treats as active
	targetTime := common.GetNowDate().Add(time.Minute * -30)
	err := o.Db.Preload("OrderedStockItemModels").Preload("OrderedFoodItemModels").Where("user_id = ? and canceled = false and pickup_date_time > ?", userId, targetTime).Order("pickup_date_time desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		order, err := model.toDomain(model.OrderedStockItemModels, model.OrderedFoodItemModels)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (o *OrderInfoRepository) Create(order *domains.OrderInfo) (string, error) {
	model, err := newOrderInfoModel(order)
	if err != nil {
		return "", err
	}
	var gError error = nil
	o.Db.Transaction(func(tx *gorm.DB) error {
		stocks := []OrderedStockItemModel{}
		for _, stock := range order.GetStockItems() {
			stockModel := OrderedStockItemModel{}
			stockModel.OrderInfoModelID = order.GetId()
			stockModel.StockItemModelID = stock.GetItemId()
			stockModel.Name = stock.GetName()
			stockModel.Price = stock.GetPrice()
			stockModel.Quantity = stock.GetQuantity()
			options := []OrderedStockOptionItemModel{}
			for _, opt := range stock.GetOptionItems() {
				option := OrderedStockOptionItemModel{}
				option.OptionItemModelID = opt.GetId()
				option.Name = opt.GetName()
				option.Price = opt.GetPrice()
				options = append(options, option)
			}
			stockModel.Options = options
			stocks = append(stocks, stockModel)
		}
		model.OrderedStockItemModels = stocks

		foods := []OrderedFoodItemModel{}
		for _, food := range order.GetFoodItems() {
			foodModel := OrderedFoodItemModel{}
			foodModel.OrderInfoModelID = order.GetId()
			foodModel.FoodItemModelID = food.GetItemId()
			foodModel.Name = food.GetName()
			foodModel.Price = food.GetPrice()
			foodModel.Quantity = food.GetQuantity()
			options := []OrderedFoodOptionItemModel{}
			for _, opt := range food.GetOptionItems() {
				option := OrderedFoodOptionItemModel{}
				option.OptionItemModelID = opt.GetId()
				option.Name = opt.GetName()
				option.Price = opt.GetPrice()
				options = append(options, option)
			}
			foodModel.Options = options
			foods = append(foods, foodModel)
		}
		model.OrderedFoodItemModels = foods

		err = o.Db.Create(&model).Error
		if err != nil {
			gError = err
			return err
		}

		return nil
	})

	if gError != nil {
		return "", gError
	}
	return order.GetId(), nil
}

func (o *OrderInfoRepository) UpdateOrderStatus(order *domains.OrderInfo) error {
	model := OrderInfoModel{}
	err := o.Db.Model(&model).Where("ID = ?", order.GetId()).Update("Canceled", order.GetCanceled()).Error
	return err
}

func (o *OrderInfoRepository) UpdateUserInfo(order *domains.OrderInfo) error {
	model := OrderInfoModel{}
	err := o.Db.Model(&model).Where("ID = ?", order.GetId()).Updates(OrderInfoModel{UserName: order.GetUserName(), UserEmail: order.GetUserEmail(), UserTelNo: order.GetUserTelNo(), Memo: order.GetMemo()}).Error
	return err
}
