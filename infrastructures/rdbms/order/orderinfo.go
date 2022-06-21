package order

import (
	"time"

	"chico/takeout/common"
	domains "chico/takeout/domains/order"
	"chico/takeout/infrastructures/rdbms"
	"chico/takeout/infrastructures/rdbms/items"

	"gorm.io/gorm"
)

type OrderInfoRepository struct {
	rdbms.BaseRepository
}

func NewOrderInfoRepository(db *gorm.DB) (*OrderInfoRepository, error) {
	return &OrderInfoRepository{
		BaseRepository: rdbms.BaseRepository{ Db: db },
	}, nil
}

type OrderInfoModel struct {
	rdbms.BaseModel
	UserID          string
	Memo            string
	OrderDateTime   time.Time
	PickupDateTime  time.Time
	Canceled        bool
	StockItemModels []items.StockItemModel `gorm:"many2many:orderInfo_stockItems;"`
	FoodItemModels  []items.FoodItemModel  `gorm:"many2many:orderInfo_foodItems;"`
}

type OrderedStockItemModel struct {
	OrderInfoModelID string `gorm:"primaryKey"`
	StockItemModelID string `gorm:"primaryKey"`
	Name             string
	Price            int
	Quantity         int
}

type OrderedFoodItemModel struct {
	OrderInfoModelID string `gorm:"primaryKey"`
	FoodItemModelID  string `gorm:"primaryKey"`
	Name             string
	Price            int
	Quantity         int
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
	model.Memo = order.GetMemo()
	model.OrderDateTime = *orderDateTime
	model.PickupDateTime = *pickupDateTime

	stockModels := []items.StockItemModel{}
	for _, stock := range model.StockItemModels {
		stockModel := items.StockItemModel{}
		stockModel.ID = stock.ID
		stockModels = append(stockModels, stockModel)
	}
	model.StockItemModels = stockModels

	foodModels := []items.FoodItemModel{}
	for _, food := range model.FoodItemModels {
		foodModel := items.FoodItemModel{}
		foodModel.ID = food.ID
		foodModels = append(foodModels, foodModel)
	}
	model.FoodItemModels = foodModels

	return model, nil
}

func (s *OrderInfoModel) toDomain(stocks []OrderedStockItemModel, foods []OrderedFoodItemModel) (*domains.OrderInfo, error) {
	stockDoms := []domains.OrderStockItem{}
	for _, stock := range stocks {
		stockDom, err := domains.NewOrderStockItem(stock.StockItemModelID, stock.Name, stock.Price, stock.Quantity)
		if err != nil {
			return nil, err
		}
		stockDoms = append(stockDoms, *stockDom)
	}
	foodDoms := []domains.OrderFoodItem{}
	for _, food := range foods {
		foodDom, err := domains.NewOrderFoodItem(food.FoodItemModelID, food.Name, food.Price, food.Quantity)
		if err != nil {
			return nil, err
		}
		foodDoms = append(foodDoms, *foodDom)
	}

	pickUp := common.ConvertTimeToDateTimeStr(s.PickupDateTime)
	ordered := common.ConvertTimeToDateTimeStr(s.OrderDateTime)
	dom, err := domains.NewOrderInfoForOrm(s.ID, s.UserID, s.Memo, pickUp, ordered, stockDoms, foodDoms, s.Canceled)
	if err != nil {
		return nil, err
	}
	return dom, nil
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
	models := []OrderInfoModel{}

	err := o.Db.Preload("StockItemModels").Preload("FoodItemModels").Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		stocks := []OrderedStockItemModel{}
		err = o.Db.Model(&model).Association("StockItemModels").Find(&stocks)
		if err != nil {
			return nil, err
		}
		foods := []OrderedFoodItemModel{}
		err = o.Db.Model(&model).Association("FoodItemModels").Find(&foods)
		if err != nil {
			return nil, err
		}
		order, err := model.toDomain(stocks, foods)
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
	err = o.Db.Where("pickup_date_time >= ? and pickup_date_time<= ?", pickupDateStart, pickupDateEnd).Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := []domains.OrderInfo{}
	for _, model := range models {
		stocks := []OrderedStockItemModel{}
		err = o.Db.Model(&model).Association("StockItemModels").Find(&stocks)
		if err != nil {
			return nil, err
		}
		foods := []OrderedFoodItemModel{}
		err = o.Db.Model(&model).Association("FoodItemModels").Find(&foods)
		if err != nil {
			return nil, err
		}
		order, err := model.toDomain(stocks, foods)
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
		err = o.Db.Create(&model).Error
		if err != nil {
			gError = err
			return err
		}
		stocks := []OrderedStockItemModel{}
		for _, stock := range order.GetStockItems() {
			stockModel := OrderedStockItemModel{}
			stockModel.OrderInfoModelID = order.GetId()
			stockModel.StockItemModelID = stock.GetItemId()
			stockModel.Name = stock.GetName()
			stockModel.Price = stock.GetPrice()
			stockModel.Quantity = stock.GetQuantity()
			stocks = append(stocks, stockModel)
		}
		if len(stocks) > 0 {
			err = o.Db.Create(&stocks).Error
			if err != nil {
				gError = err
				return err
			}
		}

		foods := []OrderedFoodItemModel{}
		for _, food := range order.GetFoodItems() {
			foodModel := OrderedFoodItemModel{}
			foodModel.OrderInfoModelID = order.GetId()
			foodModel.FoodItemModelID = food.GetItemId()
			foodModel.Name = food.GetName()
			foodModel.Price = food.GetPrice()
			foodModel.Quantity = food.GetQuantity()
			foods = append(foods, foodModel)
		}
		if len(foods) > 0 {
			err = o.Db.Create(&foods).Error
			if err != nil {
				gError = err
				return err
			}
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
