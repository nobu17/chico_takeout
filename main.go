package main

import (
	"fmt"
	"net/http"
	"time"

	itemHandler "chico/takeout/handlers/item"
	orderHandler "chico/takeout/handlers/order"
	storeHandler "chico/takeout/handlers/store"
	itemRDBMS "chico/takeout/infrastructures/rdbms/items"
	orderRDBMS "chico/takeout/infrastructures/rdbms/order"
	storeRDBMS "chico/takeout/infrastructures/rdbms/store"
	orderQueryRDBMS "chico/takeout/infrastructures/rdbms/order/query"
	itemUseCase "chico/takeout/usecase/item"
	orderUseCase "chico/takeout/usecase/order"
	orderQueryUseCase "chico/takeout/usecase/order/query"
	storeUseCase "chico/takeout/usecase/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := setUpDb()
	sqlDb, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	defer sqlDb.Close()

	r := setupRouter(db)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8086")
}

func setupRouter(db *gorm.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Setting Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
			// "http://localhost:3000/",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-Requested-With",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	r.LoadHTMLGlob("frontend/build/*.html")
	r.Static("/static", "./frontend/build/static")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	kindRepo := itemRDBMS.NewItemKindRepository(db)
	//itemKindRepo := memory.NewItemKindMemoryRepository()
	// itemkind
	kind := r.Group("/item/kind")
	{
		useCase := itemUseCase.NewItemKindUseCase(kindRepo)
		handler := itemHandler.NewItemKindHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	// kindRepo := memory.NewItemKindMemoryRepository()
	//stockRepo := memory.NewStockItemMemoryRepository()
	stockRepo := itemRDBMS.NewStockItemRepository(db)
	// stock
	stock := r.Group("/item/stock")
	{
		useCase := itemUseCase.NewStockItemUseCase(stockRepo, kindRepo)
		handler := itemHandler.NewStockItemHandler(useCase)
		stock.GET("/:id", handler.Get)
		stock.GET("/", handler.GetAll)
		stock.POST("/", handler.Post)
		stock.PUT("/:id", handler.Put)
		stock.PUT("/:id/remain", handler.PutRemain)
		stock.DELETE("/:id", handler.Delete)
	}

	//businessHoursRepo := memory.NewBusinessHoursMemoryRepository()
	businessHoursRepo := storeRDBMS.NewBusinessHoursRepository(db)
	foodRepo := itemRDBMS.NewFoodItemRepository(db)
	// foodRepo := memory.NewFoodItemMemoryRepository()
	// todo idのGET紐付け
	food := r.Group("/item/food")
	{
		useCase := itemUseCase.NewFoodItemUseCase(foodRepo, kindRepo, businessHoursRepo)
		handler := itemHandler.NewFoodItemHandler(useCase)
		food.GET("/:id", handler.Get)
		food.GET("/", handler.GetAll)
		food.POST("/", handler.Post)
		food.PUT("/:id", handler.Put)
		food.DELETE("/:id", handler.Delete)
	}

	// hour
	// spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	spBusinessHourRepo := storeRDBMS.NewSpecialBusinessHoursRepository(db)
	hour := r.Group("/store/hour")
	{
		useCase := storeUseCase.NewBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		// init
		err := useCase.InitIfNotExists()
		if err != nil {
			panic(err)
		}
		handler := storeHandler.BusinessHoursHandler(useCase)
		hour.GET("/", handler.Get)
		hour.PUT("/:id", handler.Put)
	}
	specialHour := r.Group("/store/special_hour")
	{
		useCase := storeUseCase.NewSpecialBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		handler := storeHandler.NewSpecialBusinessHourHandler(useCase)
		specialHour.GET("/:id", handler.Get)
		specialHour.GET("/", handler.GetAll)
		specialHour.POST("/", handler.Post)
		specialHour.PUT("/:id", handler.Put)
		specialHour.DELETE("/:id", handler.Delete)
	}

	//holiday
	// holidayRepo := memory.NewSpecialHolidayMemoryRepository()
	holidayRepo := storeRDBMS.NewSpecialHolidayRepository(db)
	holiday := r.Group("/store/holiday")
	{
		useCase := storeUseCase.NewSpecialHolidayUseCase(holidayRepo)
		handler := storeHandler.NewSpecialHolidayHandler(useCase)
		holiday.GET("/:id", handler.Get)
		holiday.GET("/", handler.GetAll)
		holiday.POST("/", handler.Post)
		holiday.PUT("/:id", handler.Put)
		holiday.DELETE("/:id", handler.Delete)
	}

	// order
	orderRepo, err := orderRDBMS.NewOrderInfoRepository(db)
	if err != nil {
		panic(err)
	}
	order := r.Group("/order")
	{
		// orderRepo := memory.NewOrderInfoMemoryRepository()
		useCase := orderUseCase.NewOrderInfoUseCase(orderRepo, stockRepo, foodRepo, businessHoursRepo, spBusinessHourRepo, holidayRepo)
		handler := orderHandler.NewOrderInfoHandler(useCase)
		order.GET("/:id", handler.Get)
		order.GET("/user/:userId", handler.GetByUser)
		order.GET("/user/active/:userId", handler.GetActiveByUser)
		order.POST("/", handler.PostCreate)
		order.PUT("/:id", handler.PutCancel)
	}
	orderable := r.Group("/orderable")
	{
		// orderRepo := memory.NewOrderInfoMemoryRepository()
		qService := orderQueryRDBMS.NewOrderableInfoRdbmsQueryService(db)
		useCase := orderQueryUseCase.NewOrderQueryUseCase(qService)
		handler := orderHandler.NewOrderableInfoHandler(useCase)
		orderable.GET("/", handler.Get)
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func setUpDb() *gorm.DB {
	user := "test"
	pass := "password"
	server := "localhost"
	port := "15432"
	dbName := "chicoDB"

	dsn := "host=" + server + " user=" + user + " password=" + pass + " dbname=" + dbName + " port=" + port + " sslmode=disable"
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("db connected: ", &db)

	migrateDb(db)

	if err != nil {
		panic(err.Error())
	}

	return db
}

func migrateDb(db *gorm.DB) {
	err := db.AutoMigrate(&itemRDBMS.ItemKindModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&itemRDBMS.StockItemModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&storeRDBMS.BusinessHourModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&storeRDBMS.SpecialBusinessHourModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&storeRDBMS.SpecialHolidayModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&storeRDBMS.WeekDaysModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&itemRDBMS.FoodItemModel{})
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&orderRDBMS.OrderedStockItemModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&orderRDBMS.OrderedFoodItemModel{})
	if err != nil {
		panic(err.Error())
	}
	// create join define
	err = db.SetupJoinTable(&orderRDBMS.OrderInfoModel{}, "StockItemModels", &orderRDBMS.OrderedStockItemModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.SetupJoinTable(&orderRDBMS.OrderInfoModel{}, "FoodItemModels", &orderRDBMS.OrderedFoodItemModel{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&orderRDBMS.OrderInfoModel{})
	if err != nil {
		panic(err.Error())
	}
}
