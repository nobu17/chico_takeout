package main

import (
	"net/http"
	"time"

	itemHandler "chico/takeout/handlers/item"
	orderHandler "chico/takeout/handlers/order"
	storeHandler "chico/takeout/handlers/store"
	"chico/takeout/infrastructures/memory"
	itemUseCase "chico/takeout/usecase/item"
	orderUseCase "chico/takeout/usecase/order"
	storeUseCase "chico/takeout/usecase/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8086")
}

func setupRouter() *gin.Engine {
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

	itemKindRepo := memory.NewItemKindMemoryRepository()
	// itemkind
	kind := r.Group("/item/kind")
	{
		useCase := itemUseCase.NewItemKindUseCase(itemKindRepo)
		handler := itemHandler.NewItemKindHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	kindRepo := memory.NewItemKindMemoryRepository()
	stockRepo := memory.NewStockItemMemoryRepository()
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

	businessHoursRepo := memory.NewBusinessHoursMemoryRepository()
	foodRepo := memory.NewFoodItemMemoryRepository()
	// todo idのGET紐付け
	food := r.Group("/item/food")
	{
		useCase := itemUseCase.NewFoodItemUseCase(foodRepo, itemKindRepo, businessHoursRepo)
		handler := itemHandler.NewFoodItemHandler(useCase)
		food.GET("/:id", handler.Get)
		food.GET("/", handler.GetAll)
		food.POST("/", handler.Post)
		food.PUT("/:id", handler.Put)
		food.DELETE("/:id", handler.Delete)
	}

	// hour
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	hour := r.Group("/store/hour")
	{
		useCase := storeUseCase.NewBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		handler := storeHandler.NewbusinessHoursHandler(useCase)
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
	holidayRepo := memory.NewSpecialHolidayMemoryRepository()
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
	order := r.Group("/order")
	{
		orderRepo := memory.NewOrderInfoMemoryRepository()
		useCase := orderUseCase.NewOrderInfoUseCase(orderRepo, stockRepo, foodRepo, businessHoursRepo, spBusinessHourRepo, holidayRepo)
		handler := orderHandler.NewOrderInfoHandler(useCase)
		order.GET("/:id", handler.Get)
		order.POST("/", handler.PostCreate)
		order.PUT("/:id", handler.PutCancel)
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
