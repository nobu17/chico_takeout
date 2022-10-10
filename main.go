package main

import (
	"fmt"
	"net/http"
	"time"

	"chico/takeout/common"
	itemHandler "chico/takeout/handlers/item"
	orderHandler "chico/takeout/handlers/order"
	storeHandler "chico/takeout/handlers/store"
	itemRDBMS "chico/takeout/infrastructures/rdbms/items"
	orderRDBMS "chico/takeout/infrastructures/rdbms/order"
	orderQueryRDBMS "chico/takeout/infrastructures/rdbms/order/query"
	storeRDBMS "chico/takeout/infrastructures/rdbms/store"
	"chico/takeout/infrastructures/smtp"
	"chico/takeout/middleware"
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
	cfg, err := loadConfig()
	if err != nil {
		panic("failed to load config")
	}

	db := setUpDb(cfg.Db)
	sqlDb, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	defer sqlDb.Close()

	auth := initAuthService()
	r := setupRouter(db, auth)

	go scheduleTask(db)

	r.Run(":" + cfg.AppPort)
}

func loadConfig() (*common.Config, error) {
	if err := common.InitConfig(false); err != nil {
		return nil, err
	}
	cfg := common.GetConfig()
	return &cfg, nil
}

func initAuthService() middleware.AuthService {
	service, err := middleware.NewFirebaseApp()
	if err != nil {
		fmt.Println(err)
		panic("failed to init auth service.")
	}
	return service
}

func setupRouter(db *gorm.DB, auth middleware.AuthService) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// set auth info at first
	r.Use(middleware.SetAuthInfo())

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
	r.Static("/images", "./frontend/build/images")
	r.Static("/static", "./frontend/build/static")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	kindRepo := itemRDBMS.NewItemKindRepository(db)
	kind := r.Group("/item/kind")
	{
		kind.Use(middleware.CheckAuthInfo(auth))
		useCase := itemUseCase.NewItemKindUseCase(kindRepo)
		handler := itemHandler.NewItemKindHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", middleware.CheckAdmin(), handler.Post)
		kind.PUT("/:id", middleware.CheckAdmin(), handler.Put)
		kind.DELETE("/:id", middleware.CheckAdmin(), handler.Delete)
	}

	stockRepo := itemRDBMS.NewStockItemRepository(db)
	stock := r.Group("/item/stock")
	{
		stock.Use(middleware.CheckAuthInfo(auth))
		useCase := itemUseCase.NewStockItemUseCase(stockRepo, kindRepo)
		handler := itemHandler.NewStockItemHandler(useCase)
		stock.GET("/:id", handler.Get)
		stock.GET("/", handler.GetAll)
		stock.POST("/", middleware.CheckAdmin(), handler.Post)
		stock.PUT("/:id", middleware.CheckAdmin(), handler.Put)
		stock.PUT("/:id/remain", middleware.CheckAdmin(), handler.PutRemain)
		stock.DELETE("/:id", middleware.CheckAdmin(), handler.Delete)
	}

	businessHoursRepo := storeRDBMS.NewBusinessHoursRepository(db)
	foodRepo := itemRDBMS.NewFoodItemRepository(db)
	// todo idのGET紐付け
	food := r.Group("/item/food")
	{
		food.Use(middleware.CheckAuthInfo(auth))
		useCase := itemUseCase.NewFoodItemUseCase(foodRepo, kindRepo, businessHoursRepo)
		handler := itemHandler.NewFoodItemHandler(useCase)
		food.GET("/:id", handler.Get)
		food.GET("/", handler.GetAll)
		food.POST("/", middleware.CheckAdmin(), handler.Post)
		food.PUT("/:id", middleware.CheckAdmin(), handler.Put)
		food.DELETE("/:id", middleware.CheckAdmin(), handler.Delete)
	}

	spBusinessHourRepo := storeRDBMS.NewSpecialBusinessHoursRepository(db)
	hour := r.Group("/store/hour")
	{
		hour.Use(middleware.CheckAuthInfo(auth))
		useCase := storeUseCase.NewBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		// init
		err := useCase.InitIfNotExists()
		if err != nil {
			panic(err)
		}
		handler := storeHandler.BusinessHoursHandler(useCase)
		hour.GET("/", handler.Get)
		hour.PUT("/:id", middleware.CheckAdmin(), handler.Put)
	}

	specialHour := r.Group("/store/special_hour")
	{
		specialHour.Use(middleware.CheckAuthInfo(auth))
		useCase := storeUseCase.NewSpecialBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		handler := storeHandler.NewSpecialBusinessHourHandler(useCase)
		specialHour.GET("/:id", handler.Get)
		specialHour.GET("/", handler.GetAll)
		specialHour.POST("/", middleware.CheckAdmin(), handler.Post)
		specialHour.PUT("/:id", middleware.CheckAdmin(), handler.Put)
		specialHour.DELETE("/:id", middleware.CheckAdmin(), handler.Delete)
	}

	holidayRepo := storeRDBMS.NewSpecialHolidayRepository(db)
	holiday := r.Group("/store/holiday")
	{
		holiday.Use(middleware.CheckAuthInfo(auth))
		useCase := storeUseCase.NewSpecialHolidayUseCase(holidayRepo)
		handler := storeHandler.NewSpecialHolidayHandler(useCase)
		holiday.GET("/:id", handler.Get)
		holiday.GET("/", handler.GetAll)
		holiday.POST("/", middleware.CheckAdmin(), handler.Post)
		holiday.PUT("/:id", middleware.CheckAdmin(), handler.Put)
		holiday.DELETE("/:id", middleware.CheckAdmin(), handler.Delete)
	}

	orderRepo, err := orderRDBMS.NewOrderInfoRepository(db)
	if err != nil {
		panic(err)
	}
	order := r.Group("/order")
	{
		mailer := smtp.NewSmtpSendOrderMail()
		useCase := orderUseCase.NewOrderInfoUseCase(orderRepo, stockRepo, foodRepo, businessHoursRepo, spBusinessHourRepo, holidayRepo, mailer)
		handler := orderHandler.NewOrderInfoHandler(useCase)

		order.Use(middleware.CheckAuthInfo(auth))
		order.Use(middleware.SetContext(handler.InitContext))
		order.GET("/:id", handler.Get)
		order.GET("/user/:userId", handler.GetByUser)
		order.GET("/user/active/:userId", handler.GetActiveByUser)
		order.POST("/", handler.PostCreate)
		order.PUT("/:id", handler.PutCancel)
		order.PUT("user/:userId/:orderId", handler.PutUpdateUserInfo)
		order.GET("/admin_all/", middleware.CheckAdmin(), handler.GetAll)
		order.GET("/active/:date", middleware.CheckAdmin(), handler.GetActiveByDate)
	}
	orderable := r.Group("/orderable")
	{
		orderable.Use(middleware.CheckAuthInfo(auth))
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

func setUpDb(cfg common.DbConfig) *gorm.DB {
	dsn := "host=" + cfg.Server + " user=" + cfg.User + " password=" + cfg.Pass + " dbname=" + cfg.DbName + " port=" + cfg.Port + " sslmode=disable"
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("db connected: ", &db)

	sqlDb, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDb.SetConnMaxLifetime(time.Hour)

	migrateDb(db)

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

func scheduleTask(db *gorm.DB) {
	mailer := smtp.NewSmtpSendOrderMail()
	orderRepo, err := orderRDBMS.NewOrderInfoRepository(db)
	if err != nil {
		panic(err)
	}
	useCase := orderUseCase.NewOrderTaskUseCase(orderRepo, mailer)
	task, err := common.NewDailySchedularTask("08:00", func() {
		mErr := useCase.NotifyDailyOrder(*common.GetNowDate())
		if mErr != nil {
			fmt.Println("mail send fail", mErr)
			return
		}
	})
	if err != nil {
		panic("failed to init schedular")
	}
	timer := time.NewTicker(time.Minute * 10)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			task.CheckAndExecTask()
		}
	}
}
